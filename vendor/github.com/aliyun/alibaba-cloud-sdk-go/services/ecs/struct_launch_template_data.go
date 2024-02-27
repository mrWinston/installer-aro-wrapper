package ecs

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// LaunchTemplateData is a nested struct in ecs response
type LaunchTemplateData struct {
	DeploymentSetId                string                                            `json:"DeploymentSetId" xml:"DeploymentSetId"`
	VpcId                          string                                            `json:"VpcId" xml:"VpcId"`
	SystemDiskPerformanceLevel     string                                            `json:"SystemDisk.PerformanceLevel" xml:"SystemDisk.PerformanceLevel"`
	KeyPairName                    string                                            `json:"KeyPairName" xml:"KeyPairName"`
	SecurityGroupId                string                                            `json:"SecurityGroupId" xml:"SecurityGroupId"`
	NetworkType                    string                                            `json:"NetworkType" xml:"NetworkType"`
	SpotStrategy                   string                                            `json:"SpotStrategy" xml:"SpotStrategy"`
	EnableVmOsConfig               bool                                              `json:"EnableVmOsConfig" xml:"EnableVmOsConfig"`
	Description                    string                                            `json:"Description" xml:"Description"`
	SpotDuration                   int                                               `json:"SpotDuration" xml:"SpotDuration"`
	InstanceName                   string                                            `json:"InstanceName" xml:"InstanceName"`
	SecurityEnhancementStrategy    string                                            `json:"SecurityEnhancementStrategy" xml:"SecurityEnhancementStrategy"`
	UserData                       string                                            `json:"UserData" xml:"UserData"`
	SystemDiskDiskName             string                                            `json:"SystemDisk.DiskName" xml:"SystemDisk.DiskName"`
	SystemDiskSize                 int                                               `json:"SystemDisk.Size" xml:"SystemDisk.Size"`
	SpotPriceLimit                 float64                                           `json:"SpotPriceLimit" xml:"SpotPriceLimit"`
	PasswordInherit                bool                                              `json:"PasswordInherit" xml:"PasswordInherit"`
	PrivateIpAddress               string                                            `json:"PrivateIpAddress" xml:"PrivateIpAddress"`
	ImageId                        string                                            `json:"ImageId" xml:"ImageId"`
	SystemDiskDeleteWithInstance   bool                                              `json:"SystemDisk.DeleteWithInstance" xml:"SystemDisk.DeleteWithInstance"`
	SystemDiskCategory             string                                            `json:"SystemDisk.Category" xml:"SystemDisk.Category"`
	AutoReleaseTime                string                                            `json:"AutoReleaseTime" xml:"AutoReleaseTime"`
	SystemDiskDescription          string                                            `json:"SystemDisk.Description" xml:"SystemDisk.Description"`
	ImageOwnerAlias                string                                            `json:"ImageOwnerAlias" xml:"ImageOwnerAlias"`
	HostName                       string                                            `json:"HostName" xml:"HostName"`
	SystemDiskIops                 int                                               `json:"SystemDisk.Iops" xml:"SystemDisk.Iops"`
	SystemDiskAutoSnapshotPolicyId string                                            `json:"SystemDisk.AutoSnapshotPolicyId" xml:"SystemDisk.AutoSnapshotPolicyId"`
	InternetMaxBandwidthOut        int                                               `json:"InternetMaxBandwidthOut" xml:"InternetMaxBandwidthOut"`
	InternetMaxBandwidthIn         int                                               `json:"InternetMaxBandwidthIn" xml:"InternetMaxBandwidthIn"`
	InstanceType                   string                                            `json:"InstanceType" xml:"InstanceType"`
	Period                         int                                               `json:"Period" xml:"Period"`
	InstanceChargeType             string                                            `json:"InstanceChargeType" xml:"InstanceChargeType"`
	IoOptimized                    string                                            `json:"IoOptimized" xml:"IoOptimized"`
	RamRoleName                    string                                            `json:"RamRoleName" xml:"RamRoleName"`
	VSwitchId                      string                                            `json:"VSwitchId" xml:"VSwitchId"`
	ResourceGroupId                string                                            `json:"ResourceGroupId" xml:"ResourceGroupId"`
	InternetChargeType             string                                            `json:"InternetChargeType" xml:"InternetChargeType"`
	ZoneId                         string                                            `json:"ZoneId" xml:"ZoneId"`
	Ipv6AddressCount               int                                               `json:"Ipv6AddressCount" xml:"Ipv6AddressCount"`
	SystemDiskProvisionedIops      int64                                             `json:"SystemDisk.ProvisionedIops" xml:"SystemDisk.ProvisionedIops"`
	SystemDiskBurstingEnabled      bool                                              `json:"SystemDisk.BurstingEnabled" xml:"SystemDisk.BurstingEnabled"`
	SystemDiskEncrypted            string                                            `json:"SystemDisk.Encrypted" xml:"SystemDisk.Encrypted"`
	DeletionProtection             bool                                              `json:"DeletionProtection" xml:"DeletionProtection"`
	SecurityGroupIds               SecurityGroupIdsInDescribeLaunchTemplateVersions  `json:"SecurityGroupIds" xml:"SecurityGroupIds"`
	DataDisks                      DataDisks                                         `json:"DataDisks" xml:"DataDisks"`
	NetworkInterfaces              NetworkInterfacesInDescribeLaunchTemplateVersions `json:"NetworkInterfaces" xml:"NetworkInterfaces"`
	Tags                           TagsInDescribeLaunchTemplateVersions              `json:"Tags" xml:"Tags"`
}
