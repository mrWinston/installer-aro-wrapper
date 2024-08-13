// Package gcp generates Machine objects for gcp.
package gcp

import (
	"context"
	"encoding/json"
	"fmt"

	machineapi "github.com/openshift/api/machine/v1beta1"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != gcp.Name {
		return nil, fmt.Errorf("non-GCP configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != gcp.Name {
		return nil, fmt.Errorf("non-GCP machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.GCP
	mpool := pool.Platform.GCP
	azs := mpool.Zones

	credentialsMode := config.CredentialsMode

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		azIndex := int(idx) % len(azs)
		provider, err := provider(clusterID, platform, mpool, osImage, azIndex, role, userDataSecret, credentialsMode)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}
		machine := machineapi.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.openshift.io/v1beta1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-machine-api",
				Name:      fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx),
				Labels: map[string]string{
					"machine.openshift.io/cluster-api-cluster":      clusterID,
					"machine.openshift.io/cluster-api-machine-role": role,
					"machine.openshift.io/cluster-api-machine-type": role,
				},
			},
			Spec: machineapi.MachineSpec{
				ProviderSpec: machineapi.ProviderSpec{
					Value: &runtime.RawExtension{Object: provider},
				},
				// we don't need to set Versions, because we control those via operators.
			},
		}

		machines = append(machines, machine)
	}

	return machines, nil
}

func provider(clusterID string, platform *gcp.Platform, mpool *gcp.MachinePool, osImage string, azIdx int, role, userDataSecret string, credentialsMode types.CredentialsMode) (*machineapi.GCPMachineProviderSpec, error) {
	az := mpool.Zones[azIdx]
	if len(platform.Licenses) > 0 {
		osImage = fmt.Sprintf("%s-rhcos-image", clusterID)
	} else if mpool.OSImage != nil {
		osImage = fmt.Sprintf("projects/%s/global/images/%s", mpool.OSImage.Project, mpool.OSImage.Name)
	}
	network, subnetwork, err := getNetworks(platform, clusterID, role)
	if err != nil {
		return nil, err
	}

	var encryptionKey *machineapi.GCPEncryptionKeyReference

	if mpool.OSDisk.EncryptionKey != nil {
		encryptionKey = &machineapi.GCPEncryptionKeyReference{
			KMSKey: &machineapi.GCPKMSKeyReference{
				Name:      mpool.OSDisk.EncryptionKey.KMSKey.Name,
				KeyRing:   mpool.OSDisk.EncryptionKey.KMSKey.KeyRing,
				ProjectID: mpool.OSDisk.EncryptionKey.KMSKey.ProjectID,
				Location:  mpool.OSDisk.EncryptionKey.KMSKey.Location,
			},
			KMSKeyServiceAccount: mpool.OSDisk.EncryptionKey.KMSKeyServiceAccount,
		}
	}

	instanceServiceAccount := fmt.Sprintf("%s-%s@%s.iam.gserviceaccount.com", clusterID, role[0:1], platform.ProjectID)
	// Passthrough service accounts are only needed for GCP XPN.
	if len(platform.NetworkProjectID) > 0 && credentialsMode == types.PassthroughCredentialsMode {
		sess, err := gcpconfig.GetSession(context.TODO())
		if err != nil {
			return nil, err
		}

		var found bool
		serviceAccount := make(map[string]interface{})
		err = json.Unmarshal([]byte(sess.Credentials.JSON), &serviceAccount)
		if err != nil {
			return nil, err
		}
		instanceServiceAccount, found = serviceAccount["client_email"].(string)
		if !found {
			return nil, errors.New("could not find google service account")
		}
	}

	return &machineapi.GCPMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machine.openshift.io/v1beta1",
			Kind:       "GCPMachineProviderSpec",
		},
		UserDataSecret:    &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret: &corev1.LocalObjectReference{Name: "gcp-cloud-credentials"},
		Disks: []*machineapi.GCPDisk{{
			AutoDelete:    true,
			Boot:          true,
			SizeGB:        mpool.OSDisk.DiskSizeGB,
			Type:          mpool.OSDisk.DiskType,
			Image:         osImage,
			EncryptionKey: encryptionKey,
		}},
		NetworkInterfaces: []*machineapi.GCPNetworkInterface{{
			Network:    network,
			ProjectID:  platform.NetworkProjectID,
			Subnetwork: subnetwork,
		}},
		ServiceAccounts: []machineapi.GCPServiceAccount{{
			Email:  instanceServiceAccount,
			Scopes: []string{"https://www.googleapis.com/auth/cloud-platform"},
		}},
		Tags:        append(mpool.Tags, []string{fmt.Sprintf("%s-%s", clusterID, role)}...),
		MachineType: mpool.InstanceType,
		Region:      platform.Region,
		Zone:        az,
		ProjectID:   platform.ProjectID,
	}, nil
}

// ConfigMasters assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string, publish types.PublishingStrategy) {
	var targetPools []string
	if publish == types.ExternalPublishingStrategy {
		targetPools = append(targetPools, fmt.Sprintf("%s-api", clusterID))
	}

	for _, machine := range machines {
		providerSpec := machine.Spec.ProviderSpec.Value.Object.(*machineapi.GCPMachineProviderSpec)
		providerSpec.TargetPools = targetPools
	}
}
func getNetworks(platform *gcp.Platform, clusterID, role string) (string, string, error) {
	if platform.Network == "" {
		return fmt.Sprintf("%s-network", clusterID), fmt.Sprintf("%s-%s-subnet", clusterID, role), nil
	}

	switch role {
	case "worker":
		return platform.Network, platform.ComputeSubnet, nil
	case "master":
		return platform.Network, platform.ControlPlaneSubnet, nil
	default:
		return "", "", fmt.Errorf("unrecognized machine role %s", role)
	}
}
