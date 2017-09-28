package main

import "C"

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
	ubiquity_csi_core "github.com/midoblgsm/ubiquity-csi/core"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"flag"
	"path"

	csi_utils "github.com/midoblgsm/ubiquity-csi/utils"
	"github.com/midoblgsm/ubiquity/csi"
	"github.com/midoblgsm/ubiquity/resources"
	"github.com/midoblgsm/ubiquity/utils"
	"github.com/midoblgsm/ubiquity/utils/logs"
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
	l, err := csi_utils.GetCSIEndpointListener()
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

	defer logs.InitFileLogger(logs.DEBUG, path.Join(config.LogPath, "ubiquity-csi.log"))()
	logger, logFile := utils.SetupLogger(config.LogPath, "ubiquity-csi")
	defer utils.CloseLogs(logFile)

	storageAPIURL := fmt.Sprintf("http://%s:%d/ubiquity_storage", config.UbiquityServer.Address, config.UbiquityServer.Port)
	controller, err := ubiquity_csi_core.NewController(logger, "ubiquity", storageAPIURL, config)
	if err != nil {
		logger.Printf("error-creating-controller %#v\n", err)
		panic(fmt.Sprintf("error-creating-controller", err))
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
	controller *ubiquity_csi_core.Controller
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

func (s *sp) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {

	s.Lock()
	defer s.Unlock()

	createVolumeResponse, err := s.controller.CreateVolume(*req)
	if err != nil {
		return nil, err
	}

	return &createVolumeResponse, nil
}

func (s *sp) DeleteVolume(
	ctx context.Context,
	req *csi.DeleteVolumeRequest) (
	*csi.DeleteVolumeResponse, error) {

	s.Lock()
	defer s.Unlock()

	response, err := s.controller.DeleteVolume(*req)

	if err != nil {
		// UNDEFINED
		return nil, err
	}

	return &response, nil
}

func (s *sp) ControllerPublishVolume(
	ctx context.Context,
	req *csi.ControllerPublishVolumeRequest) (
	*csi.ControllerPublishVolumeResponse, error) {

	s.Lock()
	defer s.Unlock()

	response, err := s.controller.Attach(*req)
	if err != nil {
		// UNDEFINED
		return nil, err
	}
	return &response, nil
}

func (s *sp) ControllerUnpublishVolume(
	ctx context.Context,
	req *csi.ControllerUnpublishVolumeRequest) (
	*csi.ControllerUnpublishVolumeResponse, error) {

	s.Lock()
	defer s.Unlock()

	detachResponse, err := s.controller.Detach(*req)
	if err != nil {
		// UNDEFINED
		return nil, err
	}
	return &detachResponse, nil
}

func (s *sp) ValidateVolumeCapabilities(
	ctx context.Context,
	req *csi.ValidateVolumeCapabilitiesRequest) (
	*csi.ValidateVolumeCapabilitiesResponse, error) {
	resp, err := s.controller.ValidateCapabilities(*req)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *sp) ListVolumes(
	ctx context.Context,
	req *csi.ListVolumesRequest) (
	*csi.ListVolumesResponse, error) {

	s.Lock()
	defer s.Unlock()

	listResponse, err := s.controller.ListVolumes(*req)
	if err != nil {
		// UNDEFINED
		return nil, err
	}

	return &listResponse, nil
}

func (s *sp) GetCapacity(
	ctx context.Context,
	req *csi.GetCapacityRequest) (
	*csi.GetCapacityResponse, error) {

	response, err := s.controller.GetCapacity(*req)
	if err != nil {
		// UNDEFINED
		return nil, err
	}
	return &response, nil
}

func (s *sp) ControllerGetCapabilities(
	ctx context.Context,
	req *csi.ControllerGetCapabilitiesRequest) (
	*csi.ControllerGetCapabilitiesResponse, error) {

	response, err := s.controller.ControllerGetCapabilities(*req)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

////////////////////////////////////////////////////////////////////////////////
//                             Identity Service                               //
////////////////////////////////////////////////////////////////////////////////
// Server API for Identity service
//type IdentityServer interface {
//	GetSupportedVersions(context.Context, *GetSupportedVersionsRequest) (*GetSupportedVersionsResponse, error)
//	GetPluginInfo(context.Context, *GetPluginInfoRequest) (*GetPluginInfoResponse, error)
//}
func (s *sp) GetSupportedVersions(
	ctx context.Context,
	req *csi.GetSupportedVersionsRequest) (
	*csi.GetSupportedVersionsResponse, error) {

	response, err := s.controller.GetSupportedVersions(*req)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *sp) GetPluginInfo(
	ctx context.Context,
	req *csi.GetPluginInfoRequest) (
	*csi.GetPluginInfoResponse, error) {

	response, err := s.controller.GetPluginInfos(*req)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

////////////////////////////////////////////////////////////////////////////////
//                                Node Service                                //
////////////////////////////////////////////////////////////////////////////////
// Server API for Node service
//
//type NodeServer interface {
//	NodePublishVolume(context.Context, *NodePublishVolumeRequest) (*NodePublishVolumeResponse, error)
//	NodeUnpublishVolume(context.Context, *NodeUnpublishVolumeRequest) (*NodeUnpublishVolumeResponse, error)
//	GetNodeID(context.Context, *GetNodeIDRequest) (*GetNodeIDResponse, error)
//	ProbeNode(context.Context, *ProbeNodeRequest) (*ProbeNodeResponse, error)
//	NodeGetCapabilities(context.Context, *NodeGetCapabilitiesRequest) (*NodeGetCapabilitiesResponse, error)
//}
func (s *sp) NodePublishVolume(
	ctx context.Context,
	req *csi.NodePublishVolumeRequest) (
	*csi.NodePublishVolumeResponse, error) {

	s.Lock()
	defer s.Unlock()

	response, err := s.controller.Mount(*req)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *sp) NodeUnpublishVolume(
	ctx context.Context,
	req *csi.NodeUnpublishVolumeRequest) (
	*csi.NodeUnpublishVolumeResponse, error) {

	s.Lock()
	defer s.Unlock()
	response, err := s.controller.Unount(*req)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *sp) GetNodeID(
	ctx context.Context,
	req *csi.GetNodeIDRequest) (
	*csi.GetNodeIDResponse, error) {

	response, err := s.controller.GetNodeID(*req)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *sp) ProbeNode(
	ctx context.Context,
	req *csi.ProbeNodeRequest) (
	*csi.ProbeNodeResponse, error) {
	response, err := s.controller.ProbeNode(*req)
	if err != nil {
		return nil, err
	}
	return &response, nil

}

func (s *sp) NodeGetCapabilities(
	ctx context.Context,
	req *csi.NodeGetCapabilitiesRequest) (
	*csi.NodeGetCapabilitiesResponse, error) {

	response, err := s.controller.GetNodeCapabilities(*req)
	if err != nil {
		return nil, err
	}
	return &response, nil
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
