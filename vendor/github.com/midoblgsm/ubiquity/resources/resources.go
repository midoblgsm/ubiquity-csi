/**
 * Copyright 2017 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package resources

import (
	"github.com/jinzhu/gorm"
	"github.com/midoblgsm/ubiquity/csi"
)

const (
	SpectrumScale    string = "spectrum-scale"
	SpectrumScaleNFS string = "spectrum-scale-nfs"
	SoftlayerNFS     string = "softlayer-nfs"
	SCBE             string = "scbe"
)

type UbiquityServerConfig struct {
	Port                int
	LogPath             string
	ConfigPath          string
	SpectrumScaleConfig SpectrumScaleConfig
	ScbeConfig          ScbeConfig
	BrokerConfig        BrokerConfig
	DefaultBackend      string
	LogLevel            string
}

// TODO we should consider to move dedicated backend structs to the backend resource file instead of this one.
type SpectrumScaleConfig struct {
	DefaultFilesystemName string
	NfsServerAddr         string
	SshConfig             SshConfig
	RestConfig            RestConfig
	ForceDelete           bool
}

type CredentialInfo struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Group    string `json:"group"`
}

type ConnectionInfo struct {
	CredentialInfo CredentialInfo
	Port           int
	ManagementIP   string
	SkipVerifySSL  bool
}

type ScbeConfig struct {
	ConfigPath           string // TODO consider to remove later
	ConnectionInfo       ConnectionInfo
	DefaultService       string // SCBE storage service to be used by default if not mentioned by plugin
	DefaultVolumeSize    string // The default volume size in case not specified by user
	UbiquityInstanceName string // Prefix for the volume name in the storage side (max length 15 char)

	DefaultFilesystemType string // The default filesystem type to create on new provisioned volume during attachment to the host
}

const UbiquityInstanceNameMaxSize = 15
const DefaultForScbeConfigParamDefaultVolumeSize = "1"    // if customer don't mention size, then the default is 1gb
const DefaultForScbeConfigParamDefaultFilesystem = "ext4" // if customer don't mention fstype, then the default is ext4
const PathToMountUbiquityBlockDevices = "/ubiquity/%s"    // %s is the WWN of the volume # TODO this should be moved to docker plugin side
const OptionNameForVolumeFsType = "fstype"                // the option name of the fstype and also the key in the volumeConfig

type SshConfig struct {
	User string
	Host string
	Port string
}

type RestConfig struct {
	Endpoint string
	User     string
	Password string
	Hostname string
}

type SpectrumNfsRemoteConfig struct {
	ClientConfig string
}

type BrokerConfig struct {
	ConfigPath string
	Port       int //for CF Service broker
}

type UbiquityPluginConfig struct {
	DockerPlugin            UbiquityDockerPluginConfig
	LogPath                 string
	UbiquityServer          UbiquityServerConnectionInfo
	SpectrumNfsRemoteConfig SpectrumNfsRemoteConfig
	ScbeRemoteConfig        ScbeRemoteConfig
	Backends                []string
	LogLevel                string
}

type UbiquityDockerPluginConfig struct {
	//Address          string
	Port             int
	PluginsDirectory string
}

type UbiquityServerConnectionInfo struct {
	Address string
	Port    int
}

type ScbeRemoteConfig struct {
	SkipRescanISCSI bool
}

//go:generate counterfeiter -o ../fakes/fake_storage_client.go . StorageClient

type StorageClient interface {
	Activate(activateRequest ActivateRequest) ActivateResponse
	CreateVolume(createVolumeRequest CreateVolumeRequest) CreateVolumeResponse
	RemoveVolume(removeVolumeRequest RemoveVolumeRequest) RemoveVolumeResponse
	ListVolumes(listVolumeRequest ListVolumesRequest) ListVolumesResponse
	GetVolume(getVolumeRequest GetVolumeRequest) GetVolumeResponse
	GetVolumeConfig(getVolumeConfigRequest GetVolumeConfigRequest) GetVolumeConfigResponse
	Attach(attachRequest AttachRequest) AttachResponse
	Detach(detachRequest DetachRequest) DetachResponse
}

//go:generate counterfeiter -o ../fakes/fake_mounter.go . Mounter

type Mounter interface {
	Mount(mountRequest MountRequest) MountResponse
	Unmount(unmountRequest UnmountRequest) UnmountResponse
	ActionAfterDetach(request AfterDetachRequest) AfterDetachResponse
}

type ActivateRequest struct {
	Backends []string
	Opts     map[string]string
}

type CreateVolumeRequest struct {
	Name          string
	Backend       string
	ID            csi.VolumeID
	CapacityBytes uint64
	Metadata      csi.VolumeMetadata
	Opts          map[string]interface{}
}

type RemoveVolumeRequest struct {
	Name string
}

type ListVolumesRequest struct {
	//TODO add filter
	Backends []string
}

type AttachRequest struct {
	Name string
	Host string
}

type DetachRequest struct {
	Name string
	Host string
}
type GetVolumeRequest struct {
	Name string
}
type GetVolumeConfigRequest struct {
	Name string
}
type ActivateResponse struct {
	Error error
}

type CreateVolumeResponse struct {
	Volume Volume
	Error  error
}

type RemoveVolumeResponse struct {
	Error error
}

type ListVolumesResponse struct {
	Volumes []Volume
	Error   error
}

type GenericResponse struct {
	Err string
}

type GenericRequest struct {
	Name string
}

type MountRequest struct {
	Mountpoint   string
	VolumeConfig map[string]interface{}
}
type UnmountRequest struct {
	VolumeConfig map[string]interface{}
}
type AfterDetachRequest struct {
	VolumeConfig map[string]interface{}
}
type AttachResponse struct {
	Mountpoint string
	Error      error
}

type MountResponse struct {
	Mountpoint string
	Error      error
}

type UnmountResponse struct {
	Error error
}

type DetachResponse struct {
	Error error
}

type AfterDetachResponse struct {
	Error error
}

type GetVolumeResponse struct {
	Volume Volume
	Error  error
}
type GetVolumeConfigResponse struct {
	VolumeConfig map[string]interface{}
	Error        error
}

type DockerGetResponse struct {
	Volume map[string]interface{}
	Err    string
}

type Volume struct {
	gorm.Model
	Name          string
	ID            csi.VolumeID
	CapacityBytes uint64
	Metadata      csi.VolumeMetadata
	Backend       string
	Mountpoint    string
}

type GetConfigResponse struct {
	VolumeConfig map[string]interface{}
	Err          string
}

type ListResponse struct {
	Volumes []Volume
	Err     string
}
