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

package controller

import (
	"fmt"
	"log"
	"os"
	//"os/exec"
	//"path"
	//"strings"

	//"bytes"
	//k8sresources "github.com/IBM/ubiquity-k8s/resources"
	//"github.com/IBM/ubiquity-k8s/volume"
	"github.com/IBM/ubiquity/remote"
	"github.com/IBM/ubiquity/resources"
	"github.com/IBM/ubiquity/utils"
	"github.com/container-storage-interface/spec/lib/go/csi"
	//"path/filepath"
)

//Controller this is a structure that controls volume management
type Controller struct {
	Client resources.StorageClient
	Name   string
	logger *log.Logger
	exec   utils.Executor
}

//NewController allows to instantiate a controller
func NewController(logger *log.Logger, name string, storageApiURL string, config resources.UbiquityPluginConfig) (*Controller, error) {

	remoteClient, err := remote.NewRemoteClient(logger, storageApiURL, config)
	if err != nil {
		return nil, err
	}
	return &Controller{logger: logger, Name: name, Client: remoteClient, exec: utils.NewExecutor()}, nil
}

//NewControllerWithClient is made for unit testing purposes where we can pass a fake client
func NewControllerWithClient(logger *log.Logger, client resources.StorageClient, exec utils.Executor) *Controller {
	utils.NewExecutor()
	return &Controller{logger: logger, Client: client, exec: exec}
}

func (c *Controller) CreateVolume(request csi.CreateVolumeRequest) (csi.CreateVolumeResponse, error) {
	return csi.CreateVolumeResponse{}, nil
	//in := &resources.CreateVolumeRequest{}
	//opts := make(map[string]interface{})
	//// set the volume size
	//var capacity *csi.CapacityRange
	//if capacity = req.GetCapacityRange(); capacity != nil {
	//	opts["quota"] = capacity.LimitBytes
	//}
	//
	//// set additional options
	//params := req.GetParameters()
	//for k, v := range params {
	//	opts[k] = v
	//}
	//
	//in.Name = req.GetName()
	//in.Backend = req.Parameters["backend"]

	//volumeID := csi.VolumeID{Values: map[string]string{"Name": volume.Name, "ID": fmt.Sprintf("%d", volume.ID)}}
	//volumeInfo := csi.VolumeInfo{CapacityBytes: capacity.LimitBytes,
	//	AccessMode: nil,
	//	Id:         &volumeID,
	//	Metadata:   &csi.VolumeMetadata{Values: map[string]string{"backend": volume.Backend}}}
	//return &csi.CreateVolumeResponse{
	//	Reply: &csi.CreateVolumeResponse_Result_{
	//		Result: &csi.CreateVolumeResponse_Result{
	//			VolumeInfo: &volumeInfo,
	//		},
	//	},
	//}, csi.Error_CreateVolumeError{}

}

func (c *Controller) DeleteVolume(request csi.DeleteVolumeRequest) (csi.DeleteVolumeResponse, error) {

	return csi.DeleteVolumeResponse{}, nil
}

func (c *Controller) Attach(request csi.ControllerPublishVolumeRequest) (csi.ControllerPublishVolumeResponse, error) {
	//id, ok := req.GetVolumeId().GetValues()["id"]
	//if !ok {
	//	// INVALID_VOLUME_ID
	//	return nil, csi.Error_ControllerPublishVolumeError{ErrorCode: 3, ErrorDescription: "missing id val"}
	//}
	//
	//nid := req.GetNodeId()
	//if nid == nil {
	//	// INVALID_NODE_ID
	//	return nil, csi.Error_ControllerPublishVolumeError{ErrorCode: 7, ErrorDescription: "missing node id"}
	//}
	return csi.ControllerPublishVolumeResponse{}, nil
}

func (c *Controller) Detach(request csi.ControllerUnpublishVolumeRequest) (csi.ControllerUnpublishVolumeResponse, error) {
	//id, ok := req.GetVolumeId().GetValues()["id"]
	//if !ok {
	//	// INVALID_VOLUME_ID
	//	return nil, csi.Error_ControllerUnpublishVolumeError{ErrorCode: 3, ErrorDescription: "missing id val"}
	//}
	//
	//nid := req.GetNodeId()
	//if nid == nil {
	//	// INVALID_NODE_ID
	//	return nil, csi.Error_ControllerUnpublishVolumeError{ErrorCode: 7, ErrorDescription: "missing node id"}
	//}
	//
	//nidv := nid.GetValues()
	//if len(nidv) == 0 {
	//	// INVALID_NODE_ID
	//	return nil, csi.Error_ControllerUnpublishVolumeError{ErrorCode: 7, "missing node id"}
	//}
	//
	//nidid, ok := nidv["id"]
	//if !ok {
	//	// NODE_ID_REQUIRED
	//	return  nil,csi.Error_ControllerUnpublishVolumeError{ErrorCode: 9, ErrorDescription: "node id required"}
	//}
	//
	//_ = id
	//_ = nidid
	return csi.ControllerUnpublishVolumeResponse{}, nil
}

func (c *Controller) ListVolumes(request csi.ListVolumesRequest) (csi.ListVolumesResponse, error) {
	//entries := make([]*csi.ListVolumesResponse_Result_Entry, len(listResponse.Volumes))
	//for x, volume := range listResponse.Volumes {
	//	entries[x] = &csi.ListVolumesResponse_Result_Entry{
	//		VolumeInfo: volume,
	//	}
	//}
	return csi.ListVolumesResponse{}, nil
}

func (c *Controller) ValidateCapabilities(request csi.ValidateVolumeCapabilitiesRequest) (csi.ValidateVolumeCapabilitiesResponse, error) {
	//entries := make([]*csi.ListVolumesResponse_Result_Entry, len(listResponse.Volumes))
	//for x, volume := range listResponse.Volumes {
	//	entries[x] = &csi.ListVolumesResponse_Result_Entry{
	//		VolumeInfo: volume,
	//	}
	//}
	return csi.ValidateVolumeCapabilitiesResponse{}, nil
}

func (c *Controller) GetCapacity(request csi.GetCapacityRequest) (csi.GetCapacityResponse, error) {

	return csi.GetCapacityResponse{}, nil
}

func (c *Controller) ControllerGetCapabilities(request csi.ControllerGetCapabilitiesRequest) (csi.ControllerGetCapabilitiesResponse, error) {

	return csi.ControllerGetCapabilitiesResponse{
		Reply: &csi.ControllerGetCapabilitiesResponse_Result_{
			Result: &csi.ControllerGetCapabilitiesResponse_Result{
				Capabilities: []*csi.ControllerServiceCapability{
					{
						Type: &csi.ControllerServiceCapability_Rpc{
							Rpc: &csi.ControllerServiceCapability_RPC{
								// CREATE_DELETE_VOLUME
								Type: 1,
							},
						},
					},
					{
						Type: &csi.ControllerServiceCapability_Rpc{
							Rpc: &csi.ControllerServiceCapability_RPC{
								// PUBLISH_UNPUBLISH_VOLUME
								Type: 2,
							},
						},
					},
					{
						Type: &csi.ControllerServiceCapability_Rpc{
							Rpc: &csi.ControllerServiceCapability_RPC{
								// LIST_VOLUMES
								Type: 3,
							},
						},
					},
					{
						Type: &csi.ControllerServiceCapability_Rpc{
							Rpc: &csi.ControllerServiceCapability_RPC{
								// GET_CAPACITY
								Type: 4,
							},
						},
					},
				},
			},
		},
	}, nil
}

func (c *Controller) GetSupportedVersions(request csi.GetSupportedVersionsRequest) (csi.GetSupportedVersionsResponse, error) {
	return csi.GetSupportedVersionsResponse{
		Reply: &csi.GetSupportedVersionsResponse_Result_{
			Result: &csi.GetSupportedVersionsResponse_Result{
				SupportedVersions: []*csi.Version{
					{
						Major: 0,
						Minor: 1,
						Patch: 0,
					},
				},
			},
		},
	}, nil
}

func (c *Controller) GetPluginInfos(request csi.GetPluginInfoRequest) (csi.GetPluginInfoResponse, error) {

	return csi.GetPluginInfoResponse{
		Reply: &csi.GetPluginInfoResponse_Result_{
			Result: &csi.GetPluginInfoResponse_Result{
				Name:          c.Name,
				VendorVersion: "0.1.0",
				Manifest:      nil,
			},
		},
	}, nil
}

func (c *Controller) Mount(request csi.NodePublishVolumeRequest) (csi.NodePublishVolumeResponse, error) {
	return csi.NodePublishVolumeResponse{}, nil
}

func (c *Controller) Unount(request csi.NodeUnpublishVolumeRequest) (csi.NodeUnpublishVolumeResponse, error) {
	return csi.NodeUnpublishVolumeResponse{}, nil
	//id, ok := req.GetVolumeId().GetValues()["id"]
	//if !ok {
	//// VOLUME_DOES_NOT_EXIST
	//return csi.Error_NodePublishVolumeError{ErrorCode: 2, ErrorDescription: "missing id val"}, nil
	//}
	//
	//_ = id
	//
	//return &csi.NodeUnpublishVolumeResponse{
	//Reply: &csi.NodeUnpublishVolumeResponse_Result_{
	//Result: &csi.NodeUnpublishVolumeResponse_Result{},
	//},
	//}, nil
}
func (c *Controller) GetNodeID(request csi.GetNodeIDRequest) (csi.GetNodeIDResponse, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return csi.GetNodeIDResponse{}, err
	}
	return csi.GetNodeIDResponse{
		Reply: &csi.GetNodeIDResponse_Result_{
			Result: &csi.GetNodeIDResponse_Result{
				NodeId: &csi.NodeID{
					Values: map[string]string{
						"instanceID": hostname,
					},
				},
			},
		},
	}, nil
}

func (c *Controller) ProbeNode(request csi.ProbeNodeRequest) (csi.ProbeNodeResponse, error) {
	return csi.ProbeNodeResponse{
		Reply: &csi.ProbeNodeResponse_Result_{
			Result: &csi.ProbeNodeResponse_Result{},
		},
	}, nil
}

func (c *Controller) GetNodeCapabilities(request csi.NodeGetCapabilitiesRequest) (csi.NodeGetCapabilitiesResponse, error) {
	return csi.NodeGetCapabilitiesResponse{}, nil
}

//id, ok := req.GetVolumeId().GetValues()["id"]
////Init method is to initialize the k8sresourcesvolume
//func (c *Controller) Init(config resources.UbiquityPluginConfig) k8sresources.FlexVolumeResponse {
//	c.logger.Println("controller-activate-start")
//	defer c.logger.Println("controller-activate-end")
//
//	activateRequest := resources.ActivateRequest{Backends: config.Backends}
//	err := c.Client.Activate(activateRequest)
//	if err != nil {
//		return k8sresources.FlexVolumeResponse{
//			Status:  "Failure",
//			Message: fmt.Sprintf("Plugin init failed %#v ", err),
//		}
//
//	}
//

//	return k8sresources.FlexVolumeResponse{
//		Status:  "Success",
//		Message: "Plugin init successfully",
//	}
//}
//
////Attach method attaches a volume to a host
//func (c *Controller) Attach(attachRequest k8sresources.FlexVolumeAttachRequest) k8sresources.FlexVolumeResponse {
//	c.logger.Println("controller-attach-start")
//	defer c.logger.Println("controller-attach-end")
//
//	if attachRequest.Version == k8sresources.KubernetesVersion_1_5 {
//		c.logger.Printf("k8s 1.5 attach just returning Success")
//		return k8sresources.FlexVolumeResponse{
//			Status: "Success",
//		}
//	}
//	c.logger.Printf("For k8s version 1.6 and higher, Ubiquity just returns NOT supported for Attach API. This might change in the future")
//	return k8sresources.FlexVolumeResponse{
//		Status: "Not supported",
//	}
//}
//
////GetVolumeName checks if volume is attached
//func (c *Controller) GetVolumeName(getVolumeNameRequest k8sresources.FlexVolumeGetVolumeNameRequest) k8sresources.FlexVolumeResponse {
//	c.logger.Println("controller-isAttached-start")
//	defer c.logger.Println("controller-isAttached-end")
//
//	return k8sresources.FlexVolumeResponse{
//		Status: "Not supported",
//	}
//}
//
////WaitForAttach Waits for a volume to get attached to the node
//func (c *Controller) WaitForAttach(waitForAttachRequest k8sresources.FlexVolumeWaitForAttachRequest) k8sresources.FlexVolumeResponse {
//	c.logger.Println("controller-waitForAttach-start")
//	return k8sresources.FlexVolumeResponse{
//		Status: "Not supported",
//	}
//}
//
////IsAttached checks if volume is attached
//func (c *Controller) IsAttached(isAttachedRequest k8sresources.FlexVolumeIsAttachedRequest) k8sresources.FlexVolumeResponse {
//	c.logger.Println("controller-isAttached-start")
//	return k8sresources.FlexVolumeResponse{
//		Status: "Not supported",
//	}
//}
//
////Detach detaches the volume/ fileset from the pod
//func (c *Controller) Detach(detachRequest k8sresources.FlexVolumeDetachRequest) k8sresources.FlexVolumeResponse {
//	c.logger.Println("controller-detach-start")
//	defer c.logger.Println("controller-detach-end")
//	if detachRequest.Version == k8sresources.KubernetesVersion_1_5 {
//		return k8sresources.FlexVolumeResponse{
//			Status: "Success",
//		}
//	}
//	return k8sresources.FlexVolumeResponse{
//		Status: "Not supported",
//	}
//}
//
////MountDevice mounts a device in a given location
//func (c *Controller) MountDevice(mountDeviceRequest k8sresources.FlexVolumeMountDeviceRequest) k8sresources.FlexVolumeResponse {
//	c.logger.Println("controller-MountDevice-start")
//	defer c.logger.Println("controller-MountDevice-end")
//	return k8sresources.FlexVolumeResponse{
//		Status: "Not supported",
//	}
//}
//
////UnmountDevice checks if volume is unmounted
//func (c *Controller) UnmountDevice(unmountDeviceRequest k8sresources.FlexVolumeUnmountDeviceRequest) k8sresources.FlexVolumeResponse {
//	c.logger.Println("controller-UnmountDevice-start")
//	defer c.logger.Println("controller-UnmountDevice-end")
//	return k8sresources.FlexVolumeResponse{
//		Status: "Not supported",
//	}
//}
//
////Mount method allows to mount the volume/fileset to a given location for a pod
//func (c *Controller) Mount(mountRequest k8sresources.FlexVolumeMountRequest) k8sresources.FlexVolumeResponse {
//	c.logger.Println("controller-mount-start")
//	defer c.logger.Println("controller-mount-end")
//	c.logger.Println(fmt.Sprintf("mountRequest [%#v]", mountRequest))
//	var lnPath string
//	attachRequest := resources.AttachRequest{Name: mountRequest.MountDevice, Host: getHost()}
//	mountedPath, err := c.Client.Attach(attachRequest)
//
//	if err != nil {
//		msg := fmt.Sprintf("Failed to mount volume [%s], Error: %#v", mountRequest.MountDevice, err)
//		c.logger.Println(msg)
//		return k8sresources.FlexVolumeResponse{
//			Status:  "Failure",
//			Message: msg,
//		}
//	}
//	if mountRequest.Version == k8sresources.KubernetesVersion_1_5 {
//		//For k8s 1.5, by the time we do the attach/mount, the mountDir (MountPath) is not created trying to do mount and ln will fail because the dir is not found, so we need to create the directory before continuing
//		dir := filepath.Dir(mountRequest.MountPath)
//		c.logger.Printf("mountrequest.MountPath %s", mountRequest.MountPath)
//		lnPath = mountRequest.MountPath
//		k8sRequiredMountPoint := path.Join(mountRequest.MountPath, mountRequest.MountDevice)
//		if _, err = os.Stat(k8sRequiredMountPoint); err != nil {
//			if os.IsNotExist(err) {
//
//				c.logger.Printf("creating volume directory %s", dir)
//				err = os.MkdirAll(dir, 0777)
//				if err != nil && !os.IsExist(err) {
//					msg := fmt.Sprintf("Failed creating volume directory %#v", err)
//					c.logger.Println(msg)
//
//					return k8sresources.FlexVolumeResponse{
//						Status:  "Failure",
//						Message: msg,
//					}
//
//				}
//			}
//		}
//		// For k8s 1.6 and later kubelet creates a folder as the MountPath, including the volume name, whenwe try to create the symlink this will fail because the same name exists. This is why we need to remove it before continuing.
//	} else {
//		lnPath, _ = path.Split(mountRequest.MountPath)
//		c.logger.Printf("removing folder %s", mountRequest.MountPath)
//
//		err = os.Remove(mountRequest.MountPath)
//		if err != nil && !os.IsExist(err) {
//			msg := fmt.Sprintf("Failed removing existing volume directory %#v", err)
//			c.logger.Println(msg)
//
//			return k8sresources.FlexVolumeResponse{
//				Status:  "Failure",
//				Message: msg,
//			}
//
//		}
//
//	}
//	symLinkCommand := "/bin/ln"
//	args := []string{"-s", mountedPath, lnPath}
//	c.logger.Printf(fmt.Sprintf("creating slink from %s -> %s", mountedPath, lnPath))
//
//	var stderr bytes.Buffer
//	cmd := exec.Command(symLinkCommand, args...)
//	cmd.Stderr = &stderr
//
//	err = cmd.Run()
//	if err != nil {
//		msg := fmt.Sprintf("Controller: mount failed to symlink %#v", stderr.String())
//		c.logger.Println(msg)
//		return k8sresources.FlexVolumeResponse{
//			Status:  "Failure",
//			Message: msg,
//		}
//
//	}
//	msg := fmt.Sprintf("Volume mounted successfully to %s", mountedPath)
//	c.logger.Println(msg)
//
//	return k8sresources.FlexVolumeResponse{
//		Status:  "Success",
//		Message: msg,
//	}
//}
//
////Unmount methods unmounts the volume from the pod
//func (c *Controller) Unmount(unmountRequest k8sresources.FlexVolumeUnmountRequest) k8sresources.FlexVolumeResponse {
//	c.logger.Println("Controller: unmount start")
//	defer c.logger.Println("Controller: unmount end")
//	c.logger.Printf("unmountRequest %#v", unmountRequest)
//	var detachRequest resources.DetachRequest
//	var pvName string
//
//	// Validate that the mountpoint is a symlink as ubiquity expect it to be
//	realMountPoint, err := c.exec.EvalSymlinks(unmountRequest.MountPath)
//	if err != nil {
//		msg := fmt.Sprintf("Cannot execute umount because the mountPath [%s] is not a symlink as expected. Error: %#v", unmountRequest.MountPath, err)
//		c.logger.Println(msg)
//		return k8sresources.FlexVolumeResponse{Status: "Failure", Message: msg, Device: ""}
//	}
//	ubiquityMountPrefix := fmt.Sprintf(resources.PathToMountUbiquityBlockDevices, "")
//	if strings.HasPrefix(realMountPoint, ubiquityMountPrefix) {
//		// SCBE backend flow
//		pvName = path.Base(unmountRequest.MountPath)
//
//		detachRequest = resources.DetachRequest{Name: pvName, Host: getHost()}
//		err = c.Client.Detach(detachRequest)
//		if err != nil {
//			msg := fmt.Sprintf(
//				"Failed to unmount volume [%s] on mountpoint [%s]. Error: %#v",
//				pvName,
//				unmountRequest.MountPath,
//				err)
//			c.logger.Println(msg)
//			return k8sresources.FlexVolumeResponse{Status: "Failure", Message: msg, Device: ""}
//		}
//
//		c.logger.Println(fmt.Sprintf("Removing the slink [%s] to the real mountpoint [%s]", unmountRequest.MountPath, realMountPoint))
//		err := c.exec.Remove(unmountRequest.MountPath)
//		if err != nil {
//			msg := fmt.Sprintf("fail to remove slink %s. Error %#v", unmountRequest.MountPath, err)
//			c.logger.Println(msg)
//			return k8sresources.FlexVolumeResponse{Status: "Failure", Message: msg, Device: ""}
//		}
//
//	} else {
//
//		listVolumeRequest := resources.ListVolumesRequest{}
//		volumes, err := c.Client.ListVolumes(listVolumeRequest)
//		if err != nil {
//			msg := fmt.Sprintf("Error getting the volume list from ubiquity server %#v", err)
//			c.logger.Println(msg)
//			return k8sresources.FlexVolumeResponse{
//				Status:  "Failure",
//				Message: msg,
//			}
//		}
//
//		volume, err := getVolumeForMountpoint(unmountRequest.MountPath, volumes)
//		if err != nil {
//			msg := fmt.Sprintf(
//				"Error finding the volume with mountpoint [%s] from the list of ubiquity volumes %#v. Error is : %#v",
//				unmountRequest.MountPath,
//				volumes,
//				err)
//			c.logger.Println(msg)
//			return k8sresources.FlexVolumeResponse{
//				Status:  "Failure",
//				Message: msg,
//			}
//		}
//
//		detachRequest = resources.DetachRequest{Name: volume.Name}
//		err = c.Client.Detach(detachRequest)
//		if err != nil && err.Error() != "fileset not linked" {
//			msg := fmt.Sprintf(
//				"Failed to unmount volume [%s] on mountpoint [%s]. Error: %#v",
//				volume.Name,
//				unmountRequest.MountPath,
//				err)
//			c.logger.Println(msg)
//
//			return k8sresources.FlexVolumeResponse{
//				Status:  "Failure",
//				Message: msg,
//			}
//		}
//
//		pvName = volume.Name
//	}
//
//	msg := fmt.Sprintf(
//		"Succeeded to umount volume [%s] on mountpoint [%s]",
//		pvName,
//		unmountRequest.MountPath,
//	)
//	c.logger.Println(msg)
//
//	return k8sresources.FlexVolumeResponse{
//		Status:  "Success",
//		Message: "Volume unmounted successfully",
//	}
//}

func getVolumeForMountpoint(mountpoint string, volumes []resources.Volume) (resources.Volume, error) {
	for _, volume := range volumes {
		if volume.Mountpoint == mountpoint {
			return volume, nil
		}
	}
	return resources.Volume{}, fmt.Errorf("Volume not found")
}

//TODO check os.Host
func getHost() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}
