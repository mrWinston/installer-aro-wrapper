package manifests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/util/yaml"
	k8syaml "sigs.k8s.io/yaml"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/manifests/staticnetworkconfig"
	"github.com/openshift/installer/pkg/types"
	agenttype "github.com/openshift/installer/pkg/types/agent"
)

var (
	nmStateConfigFilename = filepath.Join(clusterManifestDir, "nmstateconfig.yaml")
)

// NMStateConfig generates the nmstateconfig.yaml file.
type NMStateConfig struct {
	File                *asset.File
	StaticNetworkConfig []*models.HostStaticNetworkConfig
	Config              []*aiv1beta1.NMStateConfig
}

type nmStateConfig struct {
	Interfaces []struct {
		IPV4 struct {
			Address []struct {
				IP string `yaml:"ip,omitempty"`
			} `yaml:"address,omitempty"`
		} `yaml:"ipv4,omitempty"`
		IPV6 struct {
			Address []struct {
				IP string `yaml:"ip,omitempty"`
			} `yaml:"address,omitempty"`
		} `yaml:"ipv6,omitempty"`
	} `yaml:"interfaces,omitempty"`
}

var _ asset.WritableAsset = (*NMStateConfig)(nil)

// Name returns a human friendly name for the asset.
func (*NMStateConfig) Name() string {
	return "NMState Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*NMStateConfig) Dependencies() []asset.Asset {
	return []asset.Asset{
		&agentconfig.AgentHosts{},
		&agent.OptionalInstallConfig{},
	}
}

// Generate generates the NMStateConfig manifest.
func (n *NMStateConfig) Generate(dependencies asset.Parents) error {

	agentHosts := &agentconfig.AgentHosts{}
	installConfig := &agent.OptionalInstallConfig{}
	dependencies.Get(agentHosts, installConfig)

	staticNetworkConfig := []*models.HostStaticNetworkConfig{}
	nmStateConfigs := []*aiv1beta1.NMStateConfig{}
	var data string
	var isNetworkConfigAvailable bool

	if len(agentHosts.Hosts) == 0 {
		return nil
	}
	if err := validateHostCount(installConfig.Config, agentHosts); err != nil {
		return err
	}

	for i, host := range agentHosts.Hosts {
		if host.NetworkConfig.Raw != nil {
			isNetworkConfigAvailable = true

			nmStateConfig := aiv1beta1.NMStateConfig{
				TypeMeta: metav1.TypeMeta{
					Kind:       "NMStateConfig",
					APIVersion: "agent-install.openshift.io/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf(getNMStateConfigName(installConfig)+"-%d", i),
					Namespace: getObjectMetaNamespace(installConfig),
					Labels:    getNMStateConfigLabels(installConfig),
				},
				Spec: aiv1beta1.NMStateConfigSpec{
					NetConfig: aiv1beta1.NetConfig{
						Raw: []byte(host.NetworkConfig.Raw),
					},
				},
			}
			for _, hostInterface := range host.Interfaces {
				intrfc := aiv1beta1.Interface{
					Name:       hostInterface.Name,
					MacAddress: hostInterface.MacAddress,
				}
				nmStateConfig.Spec.Interfaces = append(nmStateConfig.Spec.Interfaces, &intrfc)

			}
			nmStateConfigs = append(nmStateConfigs, &nmStateConfig)

			staticNetworkConfig = append(staticNetworkConfig, &models.HostStaticNetworkConfig{
				MacInterfaceMap: buildMacInterfaceMap(nmStateConfig),
				NetworkYaml:     string(nmStateConfig.Spec.NetConfig.Raw),
			})

			// Marshal the nmStateConfig one at a time
			// and add a yaml separator with new line
			// so as not to marshal the nmStateConfigs
			// as a yaml list in the generated nmstateconfig.yaml
			nmStateConfigData, err := k8syaml.Marshal(nmStateConfig)

			if err != nil {
				return errors.Wrap(err, "failed to marshal agent installer NMStateConfig")
			}
			data = fmt.Sprint(data, fmt.Sprint(string(nmStateConfigData), "---\n"))
		}
	}

	if isNetworkConfigAvailable {
		n.Config = nmStateConfigs
		n.StaticNetworkConfig = staticNetworkConfig

		n.File = &asset.File{
			Filename: nmStateConfigFilename,
			Data:     []byte(data),
		}
	}
	return n.finish()
}

// Files returns the files generated by the asset.
func (n *NMStateConfig) Files() []*asset.File {
	if n.File != nil {
		return []*asset.File{n.File}
	}
	return []*asset.File{}
}

// Load returns the NMStateConfig asset from the disk.
func (n *NMStateConfig) Load(f asset.FileFetcher) (bool, error) {

	file, err := f.FetchByName(nmStateConfigFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrapf(err, "failed to load file %s", nmStateConfigFilename)
	}

	// Split up the file into multiple YAMLs if it contains NMStateConfig for more than one node
	yamlList, err := GetMultipleYamls[aiv1beta1.NMStateConfig](file.Data)
	if err != nil {
		return false, errors.Wrapf(err, "could not decode YAML for %s", nmStateConfigFilename)
	}

	var staticNetworkConfig []*models.HostStaticNetworkConfig
	var nmStateConfigList []*aiv1beta1.NMStateConfig

	for i := range yamlList {
		nmStateConfig := yamlList[i]
		staticNetworkConfig = append(staticNetworkConfig, &models.HostStaticNetworkConfig{
			MacInterfaceMap: buildMacInterfaceMap(nmStateConfig),
			NetworkYaml:     string(nmStateConfig.Spec.NetConfig.Raw),
		})
		nmStateConfigList = append(nmStateConfigList, &nmStateConfig)
	}

	n.File, n.StaticNetworkConfig, n.Config = file, staticNetworkConfig, nmStateConfigList
	if err = n.finish(); err != nil {
		return false, err
	}
	return true, nil
}

func (n *NMStateConfig) finish() error {

	if err := n.validateWithNMStateCtl(); err != nil {
		return err
	}

	if errList := n.validateNMStateConfig().ToAggregate(); errList != nil {
		return errors.Wrapf(errList, "invalid NMStateConfig configuration")
	}
	return nil
}

func (n *NMStateConfig) validateWithNMStateCtl() error {
	level := logrus.GetLevel()
	logrus.SetLevel(logrus.WarnLevel)
	staticNetworkConfigGenerator := staticnetworkconfig.New(logrus.WithField("pkg", "manifests"), staticnetworkconfig.Config{MaxConcurrentGenerations: 2})
	defer logrus.SetLevel(level)

	// Validate the network config using nmstatectl
	if err := staticNetworkConfigGenerator.ValidateStaticConfigParams(context.Background(), n.StaticNetworkConfig); err != nil {
		return errors.Wrapf(err, "staticNetwork configuration is not valid")
	}
	return nil
}

func (n *NMStateConfig) validateNMStateConfig() field.ErrorList {
	allErrs := field.ErrorList{}

	if err := n.validateNMStateLabels(); err != nil {
		allErrs = append(allErrs, err...)
	}

	return allErrs
}

func (n *NMStateConfig) validateNMStateLabels() field.ErrorList {

	var allErrs field.ErrorList

	fieldPath := field.NewPath("ObjectMeta", "Labels")

	for _, nmStateConfig := range n.Config {
		if len(nmStateConfig.ObjectMeta.Labels) == 0 {
			allErrs = append(allErrs, field.Required(fieldPath, fmt.Sprintf("%s does not have any label set", nmStateConfig.Name)))
		}
	}

	return allErrs
}

func getFirstIP(nmstateRaw []byte) (string, error) {
	var nmStateConfig nmStateConfig
	err := yaml.Unmarshal(nmstateRaw, &nmStateConfig)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling NMStateConfig: %w", err)
	}

	for _, intf := range nmStateConfig.Interfaces {
		for _, addr4 := range intf.IPV4.Address {
			if addr4.IP != "" {
				return addr4.IP, nil
			}
		}
		for _, addr6 := range intf.IPV6.Address {
			if addr6.IP != "" {
				return addr6.IP, nil
			}
		}
	}

	return "", nil
}

// GetNodeZeroIP retrieves the first IP to be set as the node0 IP.
// The method prioritizes the search by trying to scan first the NMState configs defined
// in the agent-config hosts - so that it would be possible to skip the worker nodes - and then
// the NMStateConfig.
func GetNodeZeroIP(hosts []agenttype.Host, nmStateConfigs []*aiv1beta1.NMStateConfig) (string, error) {
	rawConfigs := []aiv1beta1.RawNetConfig{}

	// Select first the configs from the hosts, if defined
	// Skip worker hosts (or without an explicit role assigned)
	for _, host := range hosts {
		if host.Role != "master" {
			continue
		}
		rawConfigs = append(rawConfigs, host.NetworkConfig.Raw)
	}

	// Add other hosts without explicit role with a lower
	// priority as potential candidates
	for _, host := range hosts {
		if host.Role != "" {
			continue
		}
		rawConfigs = append(rawConfigs, host.NetworkConfig.Raw)
	}

	// Fallback on nmstate configs (in case hosts weren't found or didn't have static configuration)
	for _, nmStateConfig := range nmStateConfigs {
		rawConfigs = append(rawConfigs, nmStateConfig.Spec.NetConfig.Raw)
	}

	// Try to look for an eligible IP
	for _, raw := range rawConfigs {
		nodeZeroIP, err := getFirstIP(raw)
		if err != nil {
			return "", fmt.Errorf("error unmarshalling NMStateConfig: %w", err)
		}
		if nodeZeroIP == "" {
			continue
		}
		if net.ParseIP(nodeZeroIP) == nil {
			return "", fmt.Errorf("could not parse static IP: %s", nodeZeroIP)
		}
		return nodeZeroIP, nil
	}

	return "", fmt.Errorf("invalid NMState configurations provided, no interface IPs set")
}

// GetNMIgnitionFiles returns the list of NetworkManager configuration files
func GetNMIgnitionFiles(staticNetworkConfig []*models.HostStaticNetworkConfig) ([]staticnetworkconfig.StaticNetworkConfigData, error) {

	level := logrus.GetLevel()
	logrus.SetLevel(logrus.WarnLevel)
	staticNetworkConfigGenerator := staticnetworkconfig.New(logrus.WithField("pkg", "manifests"), staticnetworkconfig.Config{MaxConcurrentGenerations: 2})
	defer logrus.SetLevel(level)

	networkConfigStr, err := staticNetworkConfigGenerator.FormatStaticNetworkConfigForDB(staticNetworkConfig)
	if err != nil {
		err = fmt.Errorf("error marshalling StaticNetwork configuration: %w", err)
		return nil, err
	}

	filesList, err := staticNetworkConfigGenerator.GenerateStaticNetworkConfigData(context.Background(), networkConfigStr)
	if err != nil {
		err = fmt.Errorf("failed to create StaticNetwork config data: %w", err)
		return nil, err
	}

	return filesList, err
}

// GetMultipleYamls reads a YAML file containing multiple YAML definitions of the same format
// Each specific format must be of type DecodeFormat
func GetMultipleYamls[T any](contents []byte) ([]T, error) {

	r := bytes.NewReader(contents)
	dec := yaml.NewYAMLToJSONDecoder(r)

	var outputList []T
	for {
		decodedData := new(T)
		err := dec.Decode(&decodedData)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, errors.Wrapf(err, "Error reading multiple YAMLs")
		}

		if decodedData != nil {
			outputList = append(outputList, *decodedData)
		}
	}

	return outputList, nil
}

func buildMacInterfaceMap(nmStateConfig aiv1beta1.NMStateConfig) models.MacInterfaceMap {

	// TODO - this eventually will move to another asset so the interface definition can be shared with Butane
	macInterfaceMap := make(models.MacInterfaceMap, 0, len(nmStateConfig.Spec.Interfaces))
	for _, cfg := range nmStateConfig.Spec.Interfaces {
		logrus.Debug("adding MAC interface map to host static network config - Name: ", cfg.Name, " MacAddress:", cfg.MacAddress)
		macInterfaceMap = append(macInterfaceMap, &models.MacInterfaceMapItems0{
			MacAddress:     cfg.MacAddress,
			LogicalNicName: cfg.Name,
		})
	}
	return macInterfaceMap
}

func validateHostCount(installConfig *types.InstallConfig, agentHosts *agentconfig.AgentHosts) error {
	numRequiredMasters, numRequiredWorkers := agent.GetReplicaCount(installConfig)

	numMasters := int64(0)
	numWorkers := int64(0)
	// Check for hosts explicitly defined
	for _, host := range agentHosts.Hosts {
		switch host.Role {
		case "master":
			numMasters++
		case "worker":
			numWorkers++
		}
	}

	// If role is not defined it will first be assigned as a master
	for _, host := range agentHosts.Hosts {
		if host.Role == "" {
			if numMasters < numRequiredMasters {
				numMasters++
			} else {
				numWorkers++
			}
		}
	}

	if numMasters != 0 && numMasters < numRequiredMasters {
		logrus.Warnf("not enough master hosts defined (%v) to support all the configured ControlPlane replicas (%v)", numMasters, numRequiredMasters)
	}
	if numMasters > numRequiredMasters {
		return fmt.Errorf("the number of master hosts defined (%v) exceeds the configured ControlPlane replicas (%v)", numMasters, numRequiredMasters)
	}

	if numWorkers != 0 && numWorkers < numRequiredWorkers {
		logrus.Warnf("not enough worker hosts defined (%v) to support all the configured Compute replicas (%v)", numWorkers, numRequiredWorkers)
	}
	if numWorkers > numRequiredWorkers {
		return fmt.Errorf("the number of worker hosts defined (%v) exceeds the configured Compute replicas (%v)", numWorkers, numRequiredWorkers)
	}

	return nil
}
