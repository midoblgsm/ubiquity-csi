package utils

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

// CreateVolume issues a CreateVolume request to a CSI controller.
func CreateVolume(
	ctx context.Context,
	c csi.ControllerClient,
	version *csi.Version,
	name string,
	requiredBytes, limitBytes uint64,
	fsType string, mountFlags []string,
	params map[string]string,
	callOpts ...grpc.CallOption) (volume *csi.VolumeInfo, err error) {

	if version == nil {
		return nil, ErrVersionRequired
	}

	req := &csi.CreateVolumeRequest{
		Name:       name,
		Version:    version,
		Parameters: params,
	}

	if requiredBytes > 0 || limitBytes > 0 {
		req.CapacityRange = &csi.CapacityRange{
			LimitBytes:    limitBytes,
			RequiredBytes: requiredBytes,
		}
	}

	if fsType != "" || len(mountFlags) > 0 {
		cap := &csi.VolumeCapability_MountVolume{}
		cap.FsType = fsType
		if len(mountFlags) > 0 {
			cap.MountFlags = mountFlags
		}
		req.VolumeCapabilities = []*csi.VolumeCapability{}
	}

	res, err := c.CreateVolume(ctx, req, callOpts...)
	if err != nil {
		return nil, err
	}

	// check to see if there is a csi error
	if cerr := res.GetError(); cerr != nil {
		if err := cerr.GetCreateVolumeError(); err != nil {
			return nil, fmt.Errorf(
				"error: CreateVolume failed: %d: %s",
				err.GetErrorCode(),
				err.GetErrorDescription())
		}
		if err := cerr.GetGeneralError(); err != nil {
			return nil, fmt.Errorf(
				"error: CreateVolume failed: %d: %s",
				err.GetErrorCode(),
				err.GetErrorDescription())
		}
		return nil, errors.New(cerr.String())
	}

	result := res.GetResult()
	if result == nil {
		return nil, ErrNilResult
	}

	data := result.GetVolumeInfo()
	if data == nil {
		return nil, ErrNilVolumeInfo
	}

	return data, nil
}

// ControllerPublishVolume issues a
// ControllerPublishVolume request
// to a CSI controller.
func ControllerPublishVolume(
	ctx context.Context,
	c csi.ControllerClient,
	version *csi.Version,
	volumeID *csi.VolumeHandle,
	nodeID *csi.NodeID,
	readonly bool,
	callOpts ...grpc.CallOption) (
	*csi.PublishVolumeInfo, error) {

	if version == nil {
		return nil, ErrVersionRequired
	}

	if volumeID == nil {
		return nil, ErrVolumeIDRequired
	}

	req := &csi.ControllerPublishVolumeRequest{
		Version:      version,
		VolumeHandle: volumeID,
		NodeId:       nodeID,
		Readonly:     readonly,
	}

	res, err := c.ControllerPublishVolume(ctx, req, callOpts...)
	if err != nil {
		return nil, err
	}

	// check to see if there is a csi error
	if cerr := res.GetError(); cerr != nil {
		if err := cerr.GetControllerPublishVolumeError(); err != nil {
			return nil, fmt.Errorf(
				"error: ControllerPublishVolume failed: %d: %s",
				err.GetErrorCode(),
				err.GetErrorDescription())
		}
		if err := cerr.GetGeneralError(); err != nil {
			return nil, fmt.Errorf(
				"error: ControllerPublishVolume failed: %d: %s",
				err.GetErrorCode(),
				err.GetErrorDescription())
		}
		return nil, errors.New(cerr.String())
	}

	result := res.GetResult()
	if result == nil {
		return nil, ErrNilResult
	}

	data := result.GetPublishVolumeInfo()
	if data == nil {
		return nil, ErrNilPublishVolumeInfo
	}

	return data, nil
}

// ControllerUnpublishVolume issues a
// ControllerUnpublishVolume request
// to a CSI controller.
func ControllerUnpublishVolume(
	ctx context.Context,
	c csi.ControllerClient,
	version *csi.Version,
	volumeID *csi.VolumeHandle,
	nodeID *csi.NodeID,
	callOpts ...grpc.CallOption) error {

	if version == nil {
		return ErrVersionRequired
	}

	if volumeID == nil {
		return ErrVolumeIDRequired
	}

	req := &csi.ControllerUnpublishVolumeRequest{
		Version:      version,
		VolumeHandle: volumeID,
		NodeId:       nodeID,
	}

	res, err := c.ControllerUnpublishVolume(ctx, req, callOpts...)
	if err != nil {
		return err
	}

	// check to see if there is a csi error
	if cerr := res.GetError(); cerr != nil {
		if err := cerr.GetControllerUnpublishVolumeError(); err != nil {
			return fmt.Errorf(
				"error: ControllerUnpublishVolume failed: %d: %s",
				err.GetErrorCode(),
				err.GetErrorDescription())
		}
		if err := cerr.GetGeneralError(); err != nil {
			return fmt.Errorf(
				"error: ControllerUnpublishVolume failed: %d: %s",
				err.GetErrorCode(),
				err.GetErrorDescription())
		}
		return errors.New(cerr.String())
	}

	result := res.GetResult()
	if result == nil {
		return ErrNilResult
	}

	return nil
}

// ListVolumes issues a ListVolumes request to a CSI controller.
func ListVolumes(
	ctx context.Context,
	c csi.ControllerClient,
	version *csi.Version,
	maxEntries uint32,
	startingToken string,
	callOpts ...grpc.CallOption) (
	volumes []*csi.VolumeInfo, nextToken string, err error) {

	if version == nil {
		return nil, "", ErrVersionRequired
	}

	req := &csi.ListVolumesRequest{
		MaxEntries:    maxEntries,
		StartingToken: startingToken,
		Version:       version,
	}

	res, err := c.ListVolumes(ctx, req, callOpts...)
	if err != nil {
		return nil, "", err
	}

	// check to see if there is a csi error
	if cerr := res.GetError(); cerr != nil {
		if err := cerr.GetGeneralError(); err != nil {
			return nil, "", fmt.Errorf(
				"error: ListVolumes failed: %d: %s",
				err.GetErrorCode(),
				err.GetErrorDescription())
		}
		return nil, "", errors.New(cerr.String())
	}

	result := res.GetResult()
	if result == nil {
		return nil, "", ErrNilResult
	}

	nextToken = result.GetNextToken()
	entries := result.GetEntries()

	// check to see if there are zero entries
	if len(entries) == 0 {
		return nil, nextToken, nil
	}

	volumes = make([]*csi.VolumeInfo, len(entries))

	for x, e := range entries {
		if volumes[x] = e.GetVolumeInfo(); volumes[x] == nil {
			return nil, "", ErrNilVolumeInfo
		}
	}

	return volumes, nextToken, nil
}

// GetNodeID issues a
// GetNodeID request
// to a CSI controller.
func GetNodeID(
	ctx context.Context,
	c csi.NodeClient,
	version *csi.Version,
	callOpts ...grpc.CallOption) (*csi.NodeID, error) {

	if version == nil {
		return nil, ErrVersionRequired
	}

	req := &csi.GetNodeIDRequest{
		Version: version,
	}

	res, err := c.GetNodeID(ctx, req, callOpts...)
	if err != nil {
		return nil, err
	}

	// check to see if there is a csi error
	if cerr := res.GetError(); cerr != nil {
		if err := cerr.GetGetNodeIdError(); err != nil {
			return nil, fmt.Errorf(
				"error: GetNodeID failed: %d: %s",
				err.GetErrorCode(),
				err.GetErrorDescription())
		}
		if err := cerr.GetGeneralError(); err != nil {
			return nil, fmt.Errorf(
				"error: GetNodeID failed: %d: %s",
				err.GetErrorCode(),
				err.GetErrorDescription())
		}
		return nil, errors.New(cerr.String())
	}

	result := res.GetResult()
	if result == nil {
		return nil, ErrNilResult
	}

	data := result.GetNodeId()
	if data == nil {
		return nil, ErrNilNodeID
	}

	return data, nil
}

// NodePublishVolume issues a
// NodePublishVolume request
// to a CSI controller.
func NodePublishVolume(
	ctx context.Context,
	c csi.NodeClient,
	version *csi.Version,
	volumeID *csi.VolumeHandle,
	publishVolumeInfo *csi.PublishVolumeInfo,
	targetPath string,
	volumeCapability *csi.VolumeCapability,
	readonly bool,
	callOpts ...grpc.CallOption) error {

	if version == nil {
		return ErrVersionRequired
	}

	if volumeID == nil {
		return ErrVolumeIDRequired
	}

	if volumeCapability == nil {
		return ErrVolumeCapabilityRequired
	}

	if targetPath == "" {
		return ErrInvalidTargetPath
	}

	req := &csi.NodePublishVolumeRequest{
		Version:           version,
		VolumeHandle:      volumeID,
		PublishVolumeInfo: publishVolumeInfo,
		TargetPath:        targetPath,
		VolumeCapability:  volumeCapability,
		Readonly:          readonly,
	}

	res, err := c.NodePublishVolume(ctx, req, callOpts...)
	if err != nil {
		return err
	}

	// check to see if there is a csi error
	if cerr := res.GetError(); cerr != nil {
		if err := cerr.GetNodePublishVolumeError(); err != nil {
			return fmt.Errorf(
				"error: NodePublishVolume failed: %d: %s",
				err.GetErrorCode(),
				err.GetErrorDescription())
		}
		if err := cerr.GetGeneralError(); err != nil {
			return fmt.Errorf(
				"error: NodePublishVolume failed: %d: %s",
				err.GetErrorCode(),
				err.GetErrorDescription())
		}
		return errors.New(cerr.String())
	}

	result := res.GetResult()
	if result == nil {
		return ErrNilResult
	}

	return nil
}

// NodeUnpublishVolume issues a
// NodeUnpublishVolume request
// to a CSI controller.
func NodeUnpublishVolume(
	ctx context.Context,
	c csi.NodeClient,
	version *csi.Version,
	volumeID *csi.VolumeHandle,
	targetPath string,
	callOpts ...grpc.CallOption) error {

	if version == nil {
		return ErrVersionRequired
	}

	if volumeID == nil {
		return ErrVolumeIDRequired
	}

	if targetPath == "" {
		return ErrInvalidTargetPath
	}

	req := &csi.NodeUnpublishVolumeRequest{
		Version:      version,
		VolumeHandle: volumeID,
		TargetPath:   targetPath,
	}

	res, err := c.NodeUnpublishVolume(ctx, req, callOpts...)
	if err != nil {
		return err
	}

	// check to see if there is a csi error
	if cerr := res.GetError(); cerr != nil {
		if err := cerr.GetNodeUnpublishVolumeError(); err != nil {
			return fmt.Errorf(
				"error: NodeUnpublishVolume failed: %d: %s",
				err.GetErrorCode(),
				err.GetErrorDescription())
		}
		if err := cerr.GetGeneralError(); err != nil {
			return fmt.Errorf(
				"error: NodeUnpublishVolume failed: %d: %s",
				err.GetErrorCode(),
				err.GetErrorDescription())
		}
		return errors.New(cerr.String())
	}

	result := res.GetResult()
	if result == nil {
		return ErrNilResult
	}

	return nil
}
