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

package remote

import (
	"fmt"
	"log"

	"net/http"

	"github.com/midoblgsm/ubiquity/resources"

	"reflect"

	"github.com/midoblgsm/ubiquity/remote/mounter"
	"github.com/midoblgsm/ubiquity/utils"
)

type remoteClient struct {
	logger            *log.Logger
	isActivated       bool
	isMounted         bool
	httpClient        *http.Client
	storageApiURL     string
	config            resources.UbiquityPluginConfig
	mounterPerBackend map[string]resources.Mounter
}

func NewRemoteClient(logger *log.Logger, storageApiURL string, config resources.UbiquityPluginConfig) (resources.StorageClient, error) {
	return &remoteClient{logger: logger, storageApiURL: storageApiURL, httpClient: &http.Client{}, config: config, mounterPerBackend: make(map[string]resources.Mounter)}, nil
}

func (s *remoteClient) Activate(activateRequest resources.ActivateRequest) resources.ActivateResponse {
	s.logger.Println("remoteClient: Activate start")
	defer s.logger.Println("remoteClient: Activate end")

	if s.isActivated {
		return resources.ActivateResponse{}
	}

	// call remote activate
	activateURL := utils.FormatURL(s.storageApiURL, "activate")
	response, err := utils.HttpExecute(s.httpClient, s.logger, "POST", activateURL, activateRequest)
	if err != nil {
		s.logger.Printf("Error in activate remote call %#v", err)
		return resources.ActivateResponse{Error: fmt.Errorf("Error in activate remote call")}
	}

	if response.StatusCode != http.StatusOK {
		s.logger.Printf("Error in activate remote call %#v\n", response)
		return resources.ActivateResponse{Error: utils.ExtractErrorResponse(response)}
	}
	s.logger.Println("remoteClient: Activate success")
	s.isActivated = true
	return resources.ActivateResponse{}
}

func (s *remoteClient) CreateVolume(createVolumeRequest resources.CreateVolumeRequest) resources.CreateVolumeResponse {
	s.logger.Println("remoteClient: create start")
	defer s.logger.Println("remoteClient: create end")

	createRemoteURL := utils.FormatURL(s.storageApiURL, "volumes")

	if reflect.DeepEqual(s.config.SpectrumNfsRemoteConfig, resources.SpectrumNfsRemoteConfig{}) == false {
		createVolumeRequest.Opts["nfsClientConfig"] = s.config.SpectrumNfsRemoteConfig.ClientConfig
	}

	response, err := utils.HttpExecute(s.httpClient, s.logger, "POST", createRemoteURL, createVolumeRequest)
	if err != nil {
		s.logger.Printf("Error in create volume remote call %s", err.Error())
		return resources.CreateVolumeResponse{Error: fmt.Errorf("Error in create volume remote call(http error)")}
	}

	if response.StatusCode != http.StatusOK {
		s.logger.Printf("Error in create volume remote call %#v", response)
		return resources.CreateVolumeResponse{Error: utils.ExtractErrorResponse(response)}
	}

	return resources.CreateVolumeResponse{}
}

func (s *remoteClient) RemoveVolume(removeVolumeRequest resources.RemoveVolumeRequest) resources.RemoveVolumeResponse {
	s.logger.Println("remoteClient: remove start")
	defer s.logger.Println("remoteClient: remove end")

	removeRemoteURL := utils.FormatURL(s.storageApiURL, "volumes", removeVolumeRequest.Name)

	response, err := utils.HttpExecute(s.httpClient, s.logger, "DELETE", removeRemoteURL, removeVolumeRequest)
	if err != nil {
		s.logger.Printf("Error in remove volume remote call %#v", err)
		return resources.RemoveVolumeResponse{Error: fmt.Errorf("Error in remove volume remote call")}
	}

	if response.StatusCode != http.StatusOK {
		s.logger.Printf("Error in remove volume remote call %#v", response)
		return resources.RemoveVolumeResponse{Error: utils.ExtractErrorResponse(response)}
	}

	return resources.RemoveVolumeResponse{}
}

func (s *remoteClient) GetVolume(getVolumeRequest resources.GetVolumeRequest) resources.GetVolumeResponse {
	s.logger.Println("remoteClient: get start")
	defer s.logger.Println("remoteClient: get finish")

	getRemoteURL := utils.FormatURL(s.storageApiURL, "volumes", getVolumeRequest.Name)
	response, err := utils.HttpExecute(s.httpClient, s.logger, "GET", getRemoteURL, getVolumeRequest)
	if err != nil {
		s.logger.Printf("Error in get volume remote call %#v", err)
		return resources.GetVolumeResponse{Error: fmt.Errorf("Error in get volume remote call")}
	}

	if response.StatusCode != http.StatusOK {
		s.logger.Printf("Error in get volume remote call %#v", response)
		return resources.GetVolumeResponse{Error: utils.ExtractErrorResponse(response)}
	}

	getResponse := resources.GetVolumeResponse{}
	err = utils.UnmarshalResponse(response, &getResponse)
	if err != nil {
		s.logger.Printf("Error in unmarshalling response for get remote call %#v for response %#v", err, response)
		return resources.GetVolumeResponse{Error: fmt.Errorf("Error in unmarshalling response for get remote call")}
	}

	return getResponse
}

func (s *remoteClient) GetVolumeConfig(getVolumeConfigRequest resources.GetVolumeConfigRequest) resources.GetVolumeConfigResponse {
	s.logger.Println("remoteClient: GetVolumeConfig start")
	defer s.logger.Println("remoteClient: GetVolumeConfig finish")

	getRemoteURL := utils.FormatURL(s.storageApiURL, "volumes", getVolumeConfigRequest.Name, "config")
	response, err := utils.HttpExecute(s.httpClient, s.logger, "GET", getRemoteURL, getVolumeConfigRequest)
	if err != nil {
		s.logger.Printf("Error in get volume remote call %#v", err)
		return resources.GetVolumeConfigResponse{Error: fmt.Errorf("Error in get volume remote call")}
	}

	if response.StatusCode != http.StatusOK {
		s.logger.Printf("Error in get volume remote call %#v", response)
		return resources.GetVolumeConfigResponse{Error: utils.ExtractErrorResponse(response)}
	}

	getVolumeConfigResponse := resources.GetVolumeConfigResponse{}
	err = utils.UnmarshalResponse(response, &getVolumeConfigResponse)
	if err != nil {
		s.logger.Printf("Error in unmarshalling response for get remote call %#v for response %#v", err, response)
		return resources.GetVolumeConfigResponse{Error: fmt.Errorf("Error in unmarshalling response for get remote call")}
	}

	return getVolumeConfigResponse
}

func (s *remoteClient) Attach(attachRequest resources.AttachRequest) resources.AttachResponse {
	s.logger.Println("remoteClient: attach start")
	defer s.logger.Println("remoteClient: attach end")

	attachRemoteURL := utils.FormatURL(s.storageApiURL, "volumes", attachRequest.Name, "attach")
	response, err := utils.HttpExecute(s.httpClient, s.logger, "PUT", attachRemoteURL, attachRequest)
	if err != nil {
		s.logger.Printf("Error in attach volume remote call %#v", err)
		return resources.AttachResponse{Error: fmt.Errorf("Error in attach volume remote call")}
	}

	if response.StatusCode != http.StatusOK {
		s.logger.Printf("Error in attach volume remote call %#v", response)

		return resources.AttachResponse{Error: utils.ExtractErrorResponse(response)}
	}

	attachResponse := resources.AttachResponse{}
	err = utils.UnmarshalResponse(response, &attachResponse)
	if err != nil {
		return resources.AttachResponse{Error: fmt.Errorf("Error in unmarshalling response for attach remote call")}
	}
	getVolumeConfigRequest := resources.GetVolumeConfigRequest{Name: attachRequest.Name}
	getVolumeConfigResponse := s.GetVolumeConfig(getVolumeConfigRequest)
	if getVolumeConfigResponse.Error != nil {
		return resources.AttachResponse{Error: getVolumeConfigResponse.Error}
	}
	getVolumeRequest := resources.GetVolumeRequest{Name: attachRequest.Name}
	getVolumeResponse := s.GetVolume(getVolumeRequest)
	if getVolumeResponse.Error != nil {
		return resources.AttachResponse{Error: getVolumeConfigResponse.Error}
	}

	mounter, err := s.getMounterForBackend(getVolumeResponse.Volume.Backend)
	if err != nil {
		return resources.AttachResponse{Error: fmt.Errorf("Error determining mounter for volume: %s", err.Error())}
	}
	mountRequest := resources.MountRequest{Mountpoint: attachResponse.Mountpoint, VolumeConfig: getVolumeConfigResponse.VolumeConfig}
	mountResponse := mounter.Mount(mountRequest)
	if mountResponse.Error != nil {
		return resources.AttachResponse{Error: mountResponse.Error}
	}

	return resources.AttachResponse{Mountpoint: mountResponse.Mountpoint}
}

func (s *remoteClient) Detach(detachRequest resources.DetachRequest) resources.DetachResponse {
	s.logger.Println("remoteClient: detach start")
	defer s.logger.Println("remoteClient: detach end")

	getVolumeRequest := resources.GetVolumeRequest{Name: detachRequest.Name}
	getVolumeResponse := s.GetVolume(getVolumeRequest)
	if getVolumeResponse.Error != nil {
		return resources.DetachResponse{Error: getVolumeResponse.Error}
	}
	mounter, err := s.getMounterForBackend(getVolumeResponse.Volume.Backend)
	if err != nil {
		return resources.DetachResponse{Error: fmt.Errorf("Volume not found")}
	}

	getVolumeConfigRequest := resources.GetVolumeConfigRequest{Name: detachRequest.Name}
	getVolumeConfigResponse := s.GetVolumeConfig(getVolumeConfigRequest)
	if getVolumeConfigResponse.Error != nil {
		return resources.DetachResponse{Error: getVolumeConfigResponse.Error}
	}
	unmountRequest := resources.UnmountRequest{VolumeConfig: getVolumeConfigResponse.VolumeConfig}
	unmountResponse := mounter.Unmount(unmountRequest)
	if unmountResponse.Error != nil {
		return resources.DetachResponse{Error: unmountResponse.Error}
	}

	detachRemoteURL := utils.FormatURL(s.storageApiURL, "volumes", detachRequest.Name, "detach")
	response, err := utils.HttpExecute(s.httpClient, s.logger, "PUT", detachRemoteURL, detachRequest)
	if err != nil {
		s.logger.Printf("Error in detach volume remote call %#v", err)
		return resources.DetachResponse{Error: fmt.Errorf("Error in detach volume remote call")}
	}

	if response.StatusCode != http.StatusOK {
		s.logger.Printf("Error in detach volume remote call %#v", response)
		return resources.DetachResponse{Error: utils.ExtractErrorResponse(response)}
	}

	afterDetachRequest := resources.AfterDetachRequest{VolumeConfig: getVolumeConfigResponse.VolumeConfig}
	if afterDetachResponse := mounter.ActionAfterDetach(afterDetachRequest); afterDetachResponse.Error != nil {
		s.logger.Printf(fmt.Sprintf("Error execute action after detaching the volume : %#v", err))
		return resources.DetachResponse{Error: err}
	}
	return resources.DetachResponse{}

}

func (s *remoteClient) ListVolumes(listVolumesRequest resources.ListVolumesRequest) resources.ListVolumesResponse {
	s.logger.Println("remoteClient: list start")
	defer s.logger.Println("remoteClient: list end")

	listRemoteURL := utils.FormatURL(s.storageApiURL, "volumes")
	response, err := utils.HttpExecute(s.httpClient, s.logger, "GET", listRemoteURL, listVolumesRequest)
	if err != nil {
		s.logger.Printf("Error in list volume remote call %#v", err)
		return resources.ListVolumesResponse{Error: fmt.Errorf("Error in list volume remote call")}
	}

	if response.StatusCode != http.StatusOK {
		s.logger.Printf("Error in list volume remote call %#v", err)
		return resources.ListVolumesResponse{Error: utils.ExtractErrorResponse(response)}
	}

	listVolumesResponse := resources.ListVolumesResponse{}
	err = utils.UnmarshalResponse(response, &listVolumesResponse)
	if err != nil {
		s.logger.Printf("Error in unmarshalling response for get remote call %#v for response %#v", err, response)
		return resources.ListVolumesResponse{Error: err}
	}

	return listVolumesResponse

}

// Return the mounter object. If mounter object already used(in the map mounterPerBackend) then just reuse it
func (s *remoteClient) getMounterForBackend(backend string) (resources.Mounter, error) {
	s.logger.Println("remoteClient: getMounterForVolume start")
	defer s.logger.Println("remoteClient: getMounterForVolume end")
	mounterInst, ok := s.mounterPerBackend[backend]
	if ok {
		s.logger.Printf("getMounterForVolume reuse existing mounter for backend " + backend)
		return mounterInst, nil
	} else if backend == resources.SpectrumScale {
		s.mounterPerBackend[backend] = mounter.NewSpectrumScaleMounter(s.logger)
	} else if backend == resources.SoftlayerNFS || backend == resources.SpectrumScaleNFS {
		s.mounterPerBackend[backend] = mounter.NewNfsMounter(s.logger)
	} else if backend == resources.SCBE {
		s.mounterPerBackend[backend] = mounter.NewScbeMounter(s.config.ScbeRemoteConfig)
	} else {
		return nil, fmt.Errorf("Mounter not found for backend: %s", backend)
	}
	return s.mounterPerBackend[backend], nil
}
