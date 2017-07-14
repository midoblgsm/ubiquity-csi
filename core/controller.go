package core

import (
	"log"

	"fmt"

	"github.com/IBM/ubiquity/remote"
	"github.com/IBM/ubiquity/resources"
)

type Controller struct {
	client resources.StorageClient
	logger *log.Logger
}

func NewController(logger *log.Logger, storageApiURL string, config resources.UbiquityPluginConfig) (*Controller, error) {
	remoteClient, err := remote.NewRemoteClient(logger, storageApiURL, config)
	if err != nil {
		logger.Fatal("Cannot initialize remote client")
		return nil, err
	}
	return &Controller{logger: logger, client: remoteClient}, nil
}

func NewControllerWithClient(logger *log.Logger, client resources.StorageClient) *Controller {
	return &Controller{logger: logger, client: client}
}

func (c *Controller) Activate() resources.ActivateResponse {
	c.logger.Println("Controller: activate start")
	defer c.logger.Println("Controller: activate end")

	err := c.client.Activate()

	if err != nil {
		return resources.ActivateResponse{}
	}

	return resources.ActivateResponse{Implements: []string{"VolumeDriver"}}
}

func (c *Controller) Create(createRequest *resources.CreateVolumeRequest) (resources.Volume, error) {
	c.logger.Println("Controller: create start")
	defer c.logger.Println("Controller: create end")
	c.logger.Printf("Create details %s, %#v\n", createRequest.Name, createRequest.Opts)

	err := c.client.CreateVolume(createRequest.Name, createRequest.Opts)

	if err != nil {
		return resources.Volume{}, err
	}
	getVolumeRequest := resources.GetVolumeRequest{Name: createRequest.Name}
	volume, err := c.client.GetVolume(getVolumeRequest)
	if err != nil {
		return resources.Volume{}, err
	}

	return volume, nil
}

func (c *Controller) Remove(removeRequest *resources.RemoveVolumeRequest) error {
	c.logger.Println("Controller: remove start")
	defer c.logger.Println("Controller: remove end")
	// forceDelete is set to false to enable deleting just the volume metadata
	err := c.client.RemoveVolume(removeRequest.Name)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) Attach(attachRequest map[string]string) resources.FlexVolumeResponse {
	c.logger.Println("controller-attach-start")
	defer c.logger.Println("controller-attach-end")
	c.logger.Printf("attach-details %#v\n", attachRequest)
	var attachResponse resources.FlexVolumeResponse
	volumeName, exists := attachRequest["volumeName"]
	if !exists {

		attachResponse = resources.FlexVolumeResponse{
			Status:  "Failure",
			Message: fmt.Sprintf("Failed to attach volume: VolumeName not found : #%v", attachRequest),
			Device:  volumeName,
		}
		c.logger.Printf("Failed-to-attach-volume, VolumeName found %#v ", attachRequest)
		return attachResponse

	}

	getVolumeRequest := resources.GetVolumeRequest{Name: volumeName}
	_, err := c.client.GetVolume(getVolumeRequest)

	if err != nil {
		return resources.FlexVolumeResponse{
			Status:  "Failure",
			Message: "Failed checking volume, call create before attach",
			Device:  volumeName}

	}

	return resources.FlexVolumeResponse{
		Status:  "Success",
		Message: "Volume already attached",
		Device:  volumeName,
	}

}

//Detach detaches the volume/ fileset from the pod
func (c *Controller) Detach(detachRequest resources.FlexVolumeDetachRequest) resources.FlexVolumeResponse {
	c.logger.Println("controller-detach-start")
	defer c.logger.Println("controller-detach-end")

	c.logger.Printf("detach-details %#v\n", detachRequest)

	// no-op for now, will change with latest flex api update

	return resources.FlexVolumeResponse{
		Status:  "Success",
		Message: "Volume detached successfully",
		Device:  detachRequest.Name,
	}
}

func (c *Controller) Mount(mountRequest *resources.GenericRequest) *resources.MountResponse {
	c.logger.Println("Controller: mount start")
	defer c.logger.Println("Controller: mount end")

	mountedPath, err := c.client.Attach(mountRequest.Name)
	if err != nil {
		return &resources.MountResponse{Err: err.Error()}
	}

	mountResponse := &resources.MountResponse{Mountpoint: mountedPath}
	return mountResponse
}

func (c *Controller) Unmount(unmountRequest *resources.GenericRequest) *resources.GenericResponse {
	c.logger.Println("Controller: unmount start")
	defer c.logger.Println("Controller: unmount end")

	err := c.client.Detach(unmountRequest.Name)
	if err != nil {
		return &resources.GenericResponse{Err: err.Error()}
	}
	unmountResponse := &resources.GenericResponse{}
	return unmountResponse
}

func (c *Controller) Path(pathRequest *resources.GenericRequest) *resources.MountResponse {
	c.logger.Println("Controller: path start")
	defer c.logger.Println("Controller: path end")
	volume, err := c.client.GetVolumeConfig(pathRequest.Name)
	if err != nil {
		return &resources.MountResponse{Err: err.Error()}
	}
	mountpoint, exists := volume["mountpoint"]
	if exists == false || mountpoint == "" {

		return &resources.MountResponse{Err: "volume not mounted"}
	}
	pathResponse := &resources.MountResponse{Mountpoint: mountpoint.(string)}
	return pathResponse
}

func (c *Controller) Get(getRequest *resources.GenericRequest) *resources.DockerGetResponse {
	c.logger.Println("Controller: get start")
	defer c.logger.Println("Controller: get end")
	volume, err := c.client.GetVolumeConfig(getRequest.Name)
	if err != nil {
		return &resources.DockerGetResponse{Err: err.Error()}
	}
	mountpoint, exists := volume["mountpoint"]
	if exists == false || mountpoint == "" {
		mountpoint = ""
	}

	getResponse := &resources.DockerGetResponse{Volume: resources.Volume{Name: getRequest.Name, Mountpoint: mountpoint.(string)}}
	return getResponse
}

func (c *Controller) List() *resources.ListResponse {
	c.logger.Println("Controller: list start")
	defer c.logger.Println("Controller: list end")
	volumes, err := c.client.ListVolumes()
	if err != nil {
		return &resources.ListResponse{Err: err.Error()}
	}
	listResponse := &resources.ListResponse{Volumes: volumes}
	return listResponse
}
