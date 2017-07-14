package main

import "C"

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.ibm.com/almaden-containers/ubiquity-csi/core"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"flag"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/IBM/ubiquity/resources"
	"github.com/IBM/ubiquity/utils"
	"github.com/IBM/ubiquity/utils/logs"
	"github.com/akutz/csi-examples/gocsi/csi"
)

////////////////////////////////////////////////////////////////////////////////
//                                 CLI                                        //
////////////////////////////////////////////////////////////////////////////////

var configFile = flag.String(
	"config",
	"ubiquity-client.conf",
	"config file with ubiquity client configuration params",
)

func main() {
	l, err := GetCSIEndpointListener()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to listen: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	//init the controller
	flag.Parse()
	var config resources.UbiquityPluginConfig
	fmt.Printf("Starting ubiquity plugin with %s config file\n", *configFile)
	if _, err := toml.DecodeFile(*configFile, &config); err != nil {
		fmt.Println(err)
		return
	}

	defer logs.InitFileLogger(logs.DEBUG, path.Join(config.LogPath, "ubiquity-docker-plugin.log"))()
	logger, logFile := utils.SetupLogger(config.LogPath, "ubiquity-docker-plugin")
	defer utils.CloseLogs(logFile)

	storageAPIURL := fmt.Sprintf("http://%s:%d/ubiquity_storage", config.UbiquityServer.Address, config.UbiquityServer.Port)
	controller, err := core.NewController(logger, storageAPIURL, config)
	if err != nil {
		return nil, err
	}
	s := &sp{controller: controller}
	if err := s.Serve(ctx, l); err != nil {
		fmt.Fprintf(os.Stderr, "error: grpc failed: %v\n", err)
		os.Exit(1)
	}
}

////////////////////////////////////////////////////////////////////////////////
//                              Go Plug-in                                    //
////////////////////////////////////////////////////////////////////////////////

const name = "ubiquity"

var (
	errServerStarted = errors.New("gocsi: the server has been started")
	errServerStopped = errors.New("gocsi: the server has been stopped")
)

// ServiceProviders is an exported symbol that provides a host program
// with a map of the service provider names and constructors.
var ServiceProviders = map[string]func() interface{}{
	name: func() interface{} { return &sp{name: name} },
}

type sp struct {
	sync.Mutex
	name       string
	server     *grpc.Server
	closed     bool
	controller *core.Controller
}

// ServiceProvider.Serve
func (s *sp) Serve(ctx context.Context, li net.Listener) error {
	log.Println(name + ".Serve")
	if err := func() error {
		s.Lock()
		defer s.Unlock()
		if s.closed {
			return errServerStopped
		}
		if s.server != nil {
			return errServerStarted
		}
		s.server = grpc.NewServer()
		return nil
	}(); err != nil {
		return errServerStarted
	}
	csi.RegisterControllerServer(s.server, s)
	csi.RegisterIdentityServer(s.server, s)
	csi.RegisterNodeServer(s.server, s)

	// start the grpc server
	if err := s.server.Serve(li); err != grpc.ErrServerStopped {
		return err
	}
	return errServerStopped
}

//  ServiceProvider.Stop
func (s *sp) Stop(ctx context.Context) {
	log.Println(name + ".Stop")
	s.Lock()
	defer s.Unlock()

	if s.closed || s.server == nil {
		return
	}
	s.server.Stop()
	s.server = nil
	s.closed = true
}

//  ServiceProvider.GracefulStop
func (s *sp) GracefulStop(ctx context.Context) {
	log.Println(name + ".GracefulStop")
	s.Lock()
	defer s.Unlock()

	if s.closed || s.server == nil {
		return
	}
	s.server.GracefulStop()
	s.server = nil
	s.closed = true
}

////////////////////////////////////////////////////////////////////////////////
//                            Controller Service                              //
////////////////////////////////////////////////////////////////////////////////

func (s *sp) CreateVolume(
	ctx context.Context,
	req *csi.CreateVolumeRequest) (
	*csi.CreateVolumeResponse, error) {

	s.Lock()
	defer s.Unlock()

	in := &resources.CreateVolumeRequest{}
	opts := make(map[string]interface{})
	// set the volume size
	if v := req.GetCapacityRange(); v != nil {
		opts["quota"] = v.LimitBytes
	}

	// set additional options
	params := req.GetParameters()
	for k, v := range params {
		opts[k] = v
	}

	in.Name = req.GetName()
	in.Backend = req.Parameters["backend"]

	volume, err := s.controller.Create(in)
	if err != nil {
		return csi.Error_CreateVolumeError{ErrorCode: 1, ErrorDescription: "error creating volume"}, nil
	}

	return &csi.CreateVolumeResponse{
		Reply: &csi.CreateVolumeResponse_Result_{
			Result: &csi.CreateVolumeResponse_Result{
				VolumeInfo: (volume),
			},
		},
	}, nil
}

func (s *sp) DeleteVolume(
	ctx context.Context,
	req *csi.DeleteVolumeRequest) (
	*csi.DeleteVolumeResponse, error) {

	id, ok := req.GetVolumeId().GetValues()["id"]
	if !ok {
		// INVALID_VOLUME_ID
		return csi.Error_DeleteVolumeError{3, "missing id val"}, nil
	}

	s.Lock()
	defer s.Unlock()

	removeRequest := resources.RemoveVolumeRequest{Name: id}
	err := s.controller.Remove(removeRequest)

	if err != nil {
		// UNDEFINED
		return csi.Error_DeleteVolumeError{ErrorCode: 2, ErrorDescription: "error deleting volume"}
	}

	return &csi.DeleteVolumeResponse{
		Reply: &csi.DeleteVolumeResponse_Result_{
			Result: &csi.DeleteVolumeResponse_Result{},
		},
	}, nil
}

func (s *sp) ControllerPublishVolume(
	ctx context.Context,
	req *csi.ControllerPublishVolumeRequest) (
	*csi.ControllerPublishVolumeResponse, error) {

	id, ok := req.GetVolumeId().GetValues()["id"]
	if !ok {
		// INVALID_VOLUME_ID
		return csi.Error_ControllerPublishVolumeError{ErrorCode: 3, ErrorDescription: "missing id val"}, nil
	}

	nid := req.GetNodeId()
	if nid == nil {
		// INVALID_NODE_ID
		return csi.Error_ControllerPublishVolumeError{ErrorCode: 7, ErrorDescription: "missing node id"}, nil
	}

	s.Lock()
	defer s.Unlock()

	attachRequest := resources.AttachRequest{Name: id, Host: nid}
	s.controller.Attach(attachRequest)
	return &csi.ControllerPublishVolumeResponse{
		Reply: &csi.ControllerPublishVolumeResponse_Result_{
			Result: &csi.ControllerPublishVolumeResponse_Result{
				PublishVolumeInfo: &csi.PublishVolumeInfo{
					Values: map[string]string{},
				},
			},
		},
	}, nil
}

func (s *sp) ControllerUnpublishVolume(
	ctx context.Context,
	req *csi.ControllerUnpublishVolumeRequest) (
	*csi.ControllerUnpublishVolumeResponse, error) {

	id, ok := req.GetVolumeId().GetValues()["id"]
	if !ok {
		// INVALID_VOLUME_ID
		return csi.Error_ControllerUnpublishVolumeError{ErrorCode: 3, ErrorDescription: "missing id val"}, nil
	}

	nid := req.GetNodeId()
	if nid == nil {
		// INVALID_NODE_ID
		return csi.Error_ControllerUnpublishVolumeError{ErrorCode: 7, ErrorDescription: "missing node id"}, nil
	}

	nidv := nid.GetValues()
	if len(nidv) == 0 {
		// INVALID_NODE_ID
		return csi.Error_ControllerUnpublishVolumeError{ErrorCode: 7, "missing node id"}, nil
	}

	nidid, ok := nidv["id"]
	if !ok {
		// NODE_ID_REQUIRED
		return csi.Error_ControllerUnpublishVolumeError{ErrorCode: 9, ErrorDescription: "node id required"}, nil
	}

	_ = id
	_ = nidid

	s.Lock()
	defer s.Unlock()

	return &csi.ControllerUnpublishVolumeResponse{
		Reply: &csi.ControllerUnpublishVolumeResponse_Result_{
			Result: &csi.ControllerUnpublishVolumeResponse_Result{},
		},
	}, nil
}

func (s *sp) ValidateVolumeCapabilities(
	ctx context.Context,
	req *csi.ValidateVolumeCapabilitiesRequest) (
	*csi.ValidateVolumeCapabilitiesResponse, error) {

	return nil, nil
}

func (s *sp) ListVolumes(
	ctx context.Context,
	req *csi.ListVolumesRequest) (
	*csi.ListVolumesResponse, error) {

	s.Lock()
	defer s.Unlock()

	listResponse := s.controller.List()
	entries := make([]*csi.ListVolumesResponse_Result_Entry, len(listResponse.Volumes))
	for x, volume := range listResponse.Volumes {
		entries[x] = &csi.ListVolumesResponse_Result_Entry{
			VolumeInfo: volume,
		}
	}

	return &csi.ListVolumesResponse{
		Reply: &csi.ListVolumesResponse_Result_{
			Result: &csi.ListVolumesResponse_Result{
				Entries: entries,
			},
		},
	}, nil
}

func (s *sp) GetCapacity(
	ctx context.Context,
	req *csi.GetCapacityRequest) (
	*csi.GetCapacityResponse, error) {

	return &csi.GetCapacityResponse{
		Reply: &csi.GetCapacityResponse_Result_{
			Result: &csi.GetCapacityResponse_Result{
				TotalCapacity: tib100,
			},
		},
	}, nil
}

func (s *sp) ControllerGetCapabilities(
	ctx context.Context,
	req *csi.ControllerGetCapabilitiesRequest) (
	*csi.ControllerGetCapabilitiesResponse, error) {

	return &csi.ControllerGetCapabilitiesResponse{
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

////////////////////////////////////////////////////////////////////////////////
//                             Identity Service                               //
////////////////////////////////////////////////////////////////////////////////

func (s *sp) GetSupportedVersions(
	ctx context.Context,
	req *csi.GetSupportedVersionsRequest) (
	*csi.GetSupportedVersionsResponse, error) {

	return &csi.GetSupportedVersionsResponse{
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

func (s *sp) GetPluginInfo(
	ctx context.Context,
	req *csi.GetPluginInfoRequest) (
	*csi.GetPluginInfoResponse, error) {

	return &csi.GetPluginInfoResponse{
		Reply: &csi.GetPluginInfoResponse_Result_{
			Result: &csi.GetPluginInfoResponse_Result{
				Name:          s.name,
				VendorVersion: "0.1.0",
				Manifest:      nil,
			},
		},
	}, nil
}

////////////////////////////////////////////////////////////////////////////////
//                                Node Service                                //
////////////////////////////////////////////////////////////////////////////////

func (s *sp) NodePublishVolume(
	ctx context.Context,
	req *csi.NodePublishVolumeRequest) (
	*csi.NodePublishVolumeResponse, error) {

	id, ok := req.GetVolumeId().GetValues()["id"]
	if !ok {
		// MISSING_REQUIRED_FIELD
		return csi.Error_NodePublishVolumeError{ErrorCode: 3, ErrorDescription: "missing id val"}, nil
	}

	s.Lock()
	defer s.Unlock()

	_ = id

	return &csi.NodePublishVolumeResponse{
		Reply: &csi.NodePublishVolumeResponse_Result_{
			Result: &csi.NodePublishVolumeResponse_Result{},
		},
	}, nil
}

func (s *sp) NodeUnpublishVolume(
	ctx context.Context,
	req *csi.NodeUnpublishVolumeRequest) (
	*csi.NodeUnpublishVolumeResponse, error) {

	s.Lock()
	defer s.Unlock()

	id, ok := req.GetVolumeId().GetValues()["id"]
	if !ok {
		// VOLUME_DOES_NOT_EXIST
		return csi.Error_NodePublishVolumeError{ErrorCode: 2, ErrorDescription: "missing id val"}, nil
	}

	_ = id

	return &csi.NodeUnpublishVolumeResponse{
		Reply: &csi.NodeUnpublishVolumeResponse_Result_{
			Result: &csi.NodeUnpublishVolumeResponse_Result{},
		},
	}, nil
}

func (s *sp) GetNodeID(
	ctx context.Context,
	req *csi.GetNodeIDRequest) (
	*csi.GetNodeIDResponse, error) {

	return &csi.GetNodeIDResponse{
		Reply: &csi.GetNodeIDResponse_Result_{
			Result: &csi.GetNodeIDResponse_Result{
				NodeId: &csi.NodeID{
					Values: map[string]string{
						"instanceID": os.Hostname(),
					},
				},
			},
		},
	}, nil
}

func (s *sp) ProbeNode(
	ctx context.Context,
	req *csi.ProbeNodeRequest) (
	*csi.ProbeNodeResponse, error) {

	return &csi.ProbeNodeResponse{
		Reply: &csi.ProbeNodeResponse_Result_{
			Result: &csi.ProbeNodeResponse_Result{},
		},
	}, nil
}

func (s *sp) NodeGetCapabilities(
	ctx context.Context,
	req *csi.NodeGetCapabilitiesRequest) (
	*csi.NodeGetCapabilitiesResponse, error) {

	return &csi.NodeGetCapabilitiesResponse{
		Reply: &csi.NodeGetCapabilitiesResponse_Result_{
			Result: &csi.NodeGetCapabilitiesResponse_Result{
				Capabilities: []*csi.NodeServiceCapability{
					{
						Type: &csi.NodeServiceCapability_VolumeCapability{
							VolumeCapability: &csi.VolumeCapability{
								Value: &csi.VolumeCapability_Mount{
									Mount: &csi.VolumeCapability_MountVolume{
										FsType: "ext4",
										MountFlags: []string{
											"norootsquash",
											"uid=500",
											"gid=500",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

////////////////////////////////////////////////////////////////////////////////
//                                  Utils                                     //
////////////////////////////////////////////////////////////////////////////////

const (
	kib    uint64 = 1024
	mib    uint64 = kib * 1024
	gib    uint64 = mib * 1024
	gib100 uint64 = gib * 100
	tib    uint64 = gib * 1024
	tib100 uint64 = tib * 100
)
