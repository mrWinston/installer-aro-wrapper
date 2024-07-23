package cluster

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster/alibabacloud"
	"github.com/openshift/installer/pkg/asset/cluster/aws"
	"github.com/openshift/installer/pkg/asset/cluster/azure"
	"github.com/openshift/installer/pkg/asset/cluster/baremetal"
	"github.com/openshift/installer/pkg/asset/cluster/gcp"
	"github.com/openshift/installer/pkg/asset/cluster/ibmcloud"
	"github.com/openshift/installer/pkg/asset/cluster/libvirt"
	"github.com/openshift/installer/pkg/asset/cluster/nutanix"
	"github.com/openshift/installer/pkg/asset/cluster/openstack"
	"github.com/openshift/installer/pkg/asset/cluster/ovirt"
	"github.com/openshift/installer/pkg/asset/cluster/powervs"
	"github.com/openshift/installer/pkg/asset/cluster/vsphere"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	alibabacloudtypes "github.com/openshift/installer/pkg/types/alibabacloud"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	externaltypes "github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/featuregates"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	libvirttypes "github.com/openshift/installer/pkg/types/libvirt"
	nonetypes "github.com/openshift/installer/pkg/types/none"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

const (
	metadataFileName = "metadata.json"
)

// Metadata contains information needed to destroy clusters.
type Metadata struct {
	File *asset.File
}

var _ asset.WritableAsset = (*Metadata)(nil)

// Name returns the human-friendly name of the asset.
func (m *Metadata) Name() string {
	return "Metadata"
}

// Dependencies returns the direct dependencies for the metadata
// asset.
func (m *Metadata) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		&installconfig.InstallConfig{},
		&bootstrap.Bootstrap{},
	}
}

// Generate generates the metadata asset.
func (m *Metadata) Generate(parents asset.Parents) (err error) {
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	parents.Get(clusterID, installConfig)

	featureSet := installConfig.Config.FeatureSet
	var customFS *configv1.CustomFeatureGates
	if featureSet == configv1.CustomNoUpgrade {
		customFS = featuregates.GenerateCustomFeatures(installConfig.Config.FeatureGates)
	}

	metadata := &types.ClusterMetadata{
		ClusterName:      installConfig.Config.ObjectMeta.Name,
		ClusterID:        clusterID.UUID,
		InfraID:          clusterID.InfraID,
		FeatureSet:       featureSet,
		CustomFeatureSet: customFS,
	}

	switch installConfig.Config.Platform.Name() {
	case awstypes.Name:
		metadata.ClusterPlatformMetadata.AWS = aws.Metadata(clusterID.UUID, clusterID.InfraID, installConfig.Config)
	case libvirttypes.Name:
		metadata.ClusterPlatformMetadata.Libvirt = libvirt.Metadata(installConfig.Config)
	case openstacktypes.Name:
		metadata.ClusterPlatformMetadata.OpenStack = openstack.Metadata(clusterID.InfraID, installConfig.Config)
	case azuretypes.Name:
		metadata.ClusterPlatformMetadata.Azure = azure.Metadata(installConfig.Config)
	case gcptypes.Name:
		metadata.ClusterPlatformMetadata.GCP = gcp.Metadata(installConfig.Config)
	case ibmcloudtypes.Name:
		metadata.ClusterPlatformMetadata.IBMCloud = ibmcloud.Metadata(clusterID.InfraID, installConfig.Config)
	case baremetaltypes.Name:
		metadata.ClusterPlatformMetadata.BareMetal = baremetal.Metadata(installConfig.Config)
	case ovirttypes.Name:
		metadata.ClusterPlatformMetadata.Ovirt = ovirt.Metadata(installConfig.Config)
	case vspheretypes.Name:
		metadata.ClusterPlatformMetadata.VSphere = vsphere.Metadata(installConfig.Config)
	case alibabacloudtypes.Name:
		metadata.ClusterPlatformMetadata.AlibabaCloud = alibabacloud.Metadata(installConfig.Config)
	case powervstypes.Name:
		metadata.ClusterPlatformMetadata.PowerVS = powervs.Metadata(installConfig.Config, installConfig.PowerVS)
	case externaltypes.Name, nonetypes.Name:
	case nutanixtypes.Name:
		metadata.ClusterPlatformMetadata.Nutanix = nutanix.Metadata(installConfig.Config)
	default:
		return errors.Errorf("no known platform")
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal ClusterMetadata")
	}

	m.File = &asset.File{
		Filename: metadataFileName,
		Data:     data,
	}

	return nil
}

// Files returns the metadata file generated by the asset.
func (m *Metadata) Files() []*asset.File {
	if m.File != nil {
		return []*asset.File{m.File}
	}
	return []*asset.File{}
}

// Load is a no-op, because we never want to load broken metadata from
// the disk.
func (m *Metadata) Load(f asset.FileFetcher) (found bool, err error) {
	return false, nil
}

// LoadMetadata loads the cluster metadata from an asset directory.
func LoadMetadata(dir string) (*types.ClusterMetadata, error) {
	path := filepath.Join(dir, metadataFileName)
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata *types.ClusterMetadata
	if err = json.Unmarshal(raw, &metadata); err != nil {
		return nil, errors.Wrapf(err, "failed to Unmarshal data from %q to types.ClusterMetadata", path)
	}

	return metadata, err
}
