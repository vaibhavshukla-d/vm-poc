package handler_impl

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	imagemanager "vm/internal/client/image_manager"
	inframonitor "vm/internal/client/infra_monitor"
	vmmonitor "vm/internal/client/vm_monitor"
	api "vm/internal/gen"
	"vm/internal/modals"
	"vm/internal/service"
	"vm/pkg/constants"
	"vm/pkg/dependency"

	"gorm.io/gorm"
)

// Handler implements the generated API interface
type Handler struct {
	VMService service.VMService
	deps      *dependency.Dependency
}

// NewHandler creates a new Handler instance
func NewHandler(vmService service.VMService, deps *dependency.Dependency) *Handler {
	return &Handler{
		VMService: vmService,
		deps:      deps,
	}
}

// EditVM implements the EditVM operation
func (h *Handler) EditVM(ctx context.Context, req *api.EditVM, params api.EditVMParams) (api.EditVMRes, error) {
	if err := h.validateVMExists(ctx, string(params.VMID), constants.VMReconfigure); err != nil {
		return &api.EditVMInternalServerError{
			Message: err.Error(),
		}, nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal EditVm Request: %v", err)
		return &api.EditVMInternalServerError{
			Message: "Failed to marshal EditVm Request",
		}, nil
	}
	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMReconfigure, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to marshal EditVm Request: %v", vmRequesterr)
		return &api.EditVMInternalServerError{
			Message: vmRequesterr.Error(),
		}, nil
	}

	location := constants.VMRequestBasePath + vmRequest.RequestID
	return &api.EmptyResponseHeaders{
		Location: api.NewOptString(location),
		Response: api.EmptyResponse{},
	}, nil

}

// HCIDeployVM implements the HCIDeployVM operation
func (h *Handler) HCIDeployVM(ctx context.Context, req *api.HCIDeployVM) (api.HCIDeployVMRes, error) {
	h.deps.Logger.Infof("HCIDeployVM handler invoked")

	// Validate image and get image path
	imagePath, err := h.validateImage(ctx, req.ImageSource.Value.ImageId.Value)
	if err != nil {
		return &api.HCIDeployVMInternalServerError{
			Message: err.Error(),
		}, nil
	}
	req.ImageSource.Value.ImageName = api.NewOptString(imagePath)

	// Validate host and cluster
	if err := h.validateHost(ctx, req.Destination.Value.HostId.Value, req.Destination.Value.ClusterId.Value); err != nil {
		return &api.HCIDeployVMInternalServerError{
			Message: err.Error(),
		}, nil
	}

	// Marshal the request to JSON to store as metadata
	metadata, metadataerr := json.Marshal(req)
	if metadataerr != nil {
		h.deps.Logger.Errorf("Failed to marshal HCIDeployVM Request: %v", metadataerr)
		// Handle marshaling error
		return &api.HCIDeployVMInternalServerError{
			Message: "Failed to marshal HCIDeployVM Request",
		}, nil
	}

	// Call the service to create the VM request
	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMDeploy, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VM Deploy request: %v", vmRequesterr)
		return &api.HCIDeployVMInternalServerError{
			Message: vmRequesterr.Error(),
		}, nil
	}

	// The vmRequest.RequestID is now populated by the BeforeCreate hook.
	location := constants.VMRequestBasePath + vmRequest.RequestID

	// Return the async response with the Location header
	return &api.EmptyResponseHeaders{
		Location: api.NewOptString(location),
		Response: api.EmptyResponse{},
	}, nil
}

// VMDelete implements the VMDelete operation
func (h *Handler) VMDelete(ctx context.Context, params api.VMDeleteParams) (api.VMDeleteRes, error) {
	h.deps.Logger.Infof("VMDelete handler invoked")

	if err := h.validateVMExists(ctx, string(params.VMID), constants.VMDelete); err != nil {
		return &api.VMDeleteInternalServerError{
			Message: err.Error(),
		}, nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal VMDelete Request: %v", err)
		return &api.VMDeleteInternalServerError{
			Message: "Failed to marshal VMDelete Request",
		}, nil
	}
	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMDelete, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VMDelete Request: %v", vmRequesterr)
		return &api.VMDeleteInternalServerError{
			Message: vmRequesterr.Error(),
		}, nil
	}

	location := constants.VMRequestBasePath + vmRequest.RequestID
	return &api.EmptyResponseHeaders{
		Location: api.NewOptString(location),
		Response: api.EmptyResponse{},
	}, nil
}

// VMPowerOff implements the VMPowerOff operation
func (h *Handler) VMPowerOff(ctx context.Context, params api.VMPowerOffParams) (api.VMPowerOffRes, error) {
	h.deps.Logger.Infof("VMPowerOff handler invoked")

	if err := h.validateVMExists(ctx, string(params.VMID), constants.VMPowerOff); err != nil {
		return &api.VMPowerOffInternalServerError{
			Message: err.Error(),
		}, nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal VMPowerOff Request: %v", err)
		return &api.VMPowerOffInternalServerError{
			Message: "Failed to marshal VMPowerOff Request",
		}, nil
	}
	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMPowerOff, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VMPowerOff Request: %v", vmRequesterr)
		return &api.VMPowerOffInternalServerError{
			Message: vmRequesterr.Error(),
		}, nil
	}

	location := constants.VMRequestBasePath + vmRequest.RequestID
	return &api.EmptyResponseHeaders{
		Location: api.NewOptString(location),
		Response: api.EmptyResponse{},
	}, nil

}

// VMPowerOn implements the VMPowerOn operation
func (h *Handler) VMPowerOn(ctx context.Context, params api.VMPowerOnParams) (api.VMPowerOnRes, error) {
	h.deps.Logger.Infof("VMPowerOn handler invoked")

	if err := h.validateVMExists(ctx, string(params.VMID), constants.VMPowerOn); err != nil {
		return &api.VMPowerOnInternalServerError{
			Message: err.Error(),
		}, nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal VMPowerOn params: %v", err)
		return &api.VMPowerOnInternalServerError{
			Message: "Failed to marshal VMPowerOn params",
		}, nil
	}

	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMPowerOn, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VM power on request: %v", vmRequesterr)
		return &api.VMPowerOnInternalServerError{
			Message: vmRequesterr.Error(),
		}, nil
	}

	location := constants.VMRequestBasePath + vmRequest.RequestID
	return &api.EmptyResponseHeaders{
		Location: api.NewOptString(location),
		Response: api.EmptyResponse{},
	}, nil

}

// VMPowerReset implements the VMPowerReset operation
func (h *Handler) VMPowerReset(ctx context.Context, params api.VMPowerResetParams) (api.VMPowerResetRes, error) {
	h.deps.Logger.Infof("VMPowerReset handler invoked")

	if err := h.validateVMExists(ctx, string(params.VMID), constants.VMReset); err != nil {
		return &api.VMPowerResetInternalServerError{
			Message: err.Error(),
		}, nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal VMPowerReset params: %v", err)
		return &api.VMPowerResetInternalServerError{
			Message: "Failed to marshal VMPowerReset params",
		}, nil
	}

	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMReset, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VM power reset request: %v", vmRequesterr)
		return &api.VMPowerResetInternalServerError{
			Message: vmRequesterr.Error(),
		}, nil
	}

	location := constants.VMRequestBasePath + vmRequest.RequestID
	return &api.EmptyResponseHeaders{
		Location: api.NewOptString(location),
		Response: api.EmptyResponse{},
	}, nil

}

// VMRefresh implements the VMRefresh operation
func (h *Handler) VMRefresh(ctx context.Context, params api.VMRefreshParams) (api.VMRefreshRes, error) {
	h.deps.Logger.Infof("VMRefresh handler invoked")

	if err := h.validateVMExists(ctx, string(params.VMID), constants.VMRefresh); err != nil {
		return &api.VMRefreshInternalServerError{
			Message: err.Error(),
		}, nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal VMRefresh params: %v", err)
		return &api.VMRefreshInternalServerError{
			Message: "Failed to marshal VMRefresh params",
		}, nil
	}

	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMRefresh, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VM refresh request: %v", vmRequesterr)
		return &api.VMRefreshInternalServerError{
			Message: vmRequesterr.Error(),
		}, nil
	}

	location := constants.VMRequestBasePath + vmRequest.RequestID
	return &api.EmptyResponseHeaders{
		Location: api.NewOptString(location),
		Response: api.EmptyResponse{},
	}, nil

}

// VMRestartGuestOS implements the VMRestartGuestOS operation
func (h *Handler) VMRestartGuestOS(ctx context.Context, params api.VMRestartGuestOSParams) (api.VMRestartGuestOSRes, error) {
	h.deps.Logger.Infof("VMRestartGuestOS handler invoked")

	if err := h.validateVMExists(ctx, string(params.VMID), constants.VMRestartGuestOS); err != nil {
		return &api.VMRestartGuestOSInternalServerError{
			Message: err.Error(),
		}, nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal VMRestartGuestOS params: %v", err)
		return &api.VMRestartGuestOSInternalServerError{
			Message: "Failed to marshal VMRestartGuestOS params",
		}, nil
	}

	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMRestartGuestOS, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VM restart guest OS request: %v", vmRequesterr)
		return &api.VMRestartGuestOSInternalServerError{
			Message: vmRequesterr.Error(),
		}, nil
	}

	location := constants.VMRequestBasePath + vmRequest.RequestID
	return &api.EmptyResponseHeaders{
		Location: api.NewOptString(location),
		Response: api.EmptyResponse{},
	}, nil

}

// VMShutdownGuestOS implements the VMShutdownGuestOS operation
func (h *Handler) VMShutdownGuestOS(ctx context.Context, params api.VMShutdownGuestOSParams) (api.VMShutdownGuestOSRes, error) {
	h.deps.Logger.Infof("VMShutdownGuestOS handler invoked")

	if err := h.validateVMExists(ctx, string(params.VMID), constants.VMShutdownGuestOS); err != nil {
		return &api.VMShutdownGuestOSInternalServerError{
			Message: err.Error(),
		}, nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal VMShutdownGuestOS params: %v", err)
		return &api.VMShutdownGuestOSInternalServerError{
			Message: "Failed to marshal VMShutdownGuestOS params",
		}, nil
	}

	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMShutdownGuestOS, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VM shutdown guest OS request: %v", vmRequesterr)
		return &api.VMShutdownGuestOSInternalServerError{
			Message: vmRequesterr.Error(),
		}, nil
	}

	location := constants.VMRequestBasePath + vmRequest.RequestID
	return &api.EmptyResponseHeaders{
		Location: api.NewOptString(location),
		Response: api.EmptyResponse{},
	}, nil

}

// GetVirtualMachineRequest implements the GetVirtualMachineRequest operation
func (h *Handler) GetVirtualMachineRequest(ctx context.Context, params api.GetVirtualMachineRequestParams) (api.GetVirtualMachineRequestRes, error) {
	h.deps.Logger.Infof("GetVirtualMachineRequest handler invoked")
	vmRequest, err := h.VMService.GetVMRequest(ctx, params.RequestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &api.GetVirtualMachineRequestNotFound{
				Message: "VM request not found",
			}, nil
		}
		h.deps.Logger.Errorf("Failed to get VM request: %v", err)
		return &api.GetVirtualMachineRequestInternalServerError{
			Message: err.Error(),
		}, nil
	}

	deployInstances, deployInstanceserr := h.VMService.GetVMDeployInstances(ctx, params.RequestID)
	if deployInstanceserr != nil {
		h.deps.Logger.Warnf("Failed to get VM deploy instances, but continuing execution as they are optional: %v", deployInstanceserr)
		if errors.Is(deployInstanceserr, gorm.ErrRecordNotFound) {
			deployInstances = []*modals.VMDeployInstance{} // Initialize with an empty slice
		} else {
			return &api.GetVirtualMachineRequestInternalServerError{
				Message: deployInstanceserr.Error(),
			}, nil
		}
	}

	apiVMRequest := api.VMRequest{
		RequestId:       vmRequest.RequestID,
		Operation:       api.VMRequestOperation(vmRequest.Operation),
		RequestStatus:   api.VMRequestRequestStatus(vmRequest.RequestStatus),
		WorkspaceId:     api.NewOptString(vmRequest.WorkspaceId),
		DatacenterId:    api.NewOptString(vmRequest.DatacenterId),
		CreatedAt:       vmRequest.CreatedAt,
		RequestMetadata: vmRequest.RequestMetadata,
	}
	if vmRequest.CompletedAt != nil {
		apiVMRequest.CompletedAt = api.NewOptNilDateTime(*vmRequest.CompletedAt)
	}

	apiDeployList := make([]api.VMDeployInstance, len(deployInstances))
	for i, inst := range deployInstances {
		apiDeployList[i] = api.VMDeployInstance{
			RequestId:      inst.RequestID,
			VmId:           api.NewOptString(inst.VMID),
			VmName:         inst.VMName,
			VmStatus:       inst.VMStatus,
			VmStateMessage: api.NewOptString(inst.VMStateMessage),
		}
		if inst.CompletedAt != nil {
			apiDeployList[i].CompletedAt = api.NewOptNilDateTime(*inst.CompletedAt)
		}
	}

	return &api.VMRequestWithDeploy{
		VMRequest:    apiVMRequest,
		VMDeployList: apiDeployList,
	}, nil
}

// GetVirtualMachineRequestList implements the GetVirtualMachineRequestList operation.
func (h *Handler) GetVirtualMachineRequestList(ctx context.Context) (api.GetVirtualMachineRequestListRes, error) {
	h.deps.Logger.Infof("GetVirtualMachineRequestList handler invoked")

	// Call the service to get the list of VM requests and their instances
	vmRequests, deployInstances, reqCount, instCount, err := h.VMService.GetAllVMRequestsWithInstances(ctx)
	if err != nil {
		h.deps.Logger.Errorf("Failed to get VM request list: %v", err)
		return &api.GetVirtualMachineRequestListInternalServerError{
			Message: err.Error(),
		}, nil
	}

	// Create the response structure for VM requests
	apiVMRequests := make([]api.VMRequest, len(vmRequests))
	for i, vmRequest := range vmRequests {
		apiVMRequests[i] = api.VMRequest{
			RequestId:       vmRequest.RequestID,
			Operation:       api.VMRequestOperation(vmRequest.Operation),
			RequestStatus:   api.VMRequestRequestStatus(vmRequest.RequestStatus),
			WorkspaceId:     api.NewOptString(vmRequest.WorkspaceId),
			DatacenterId:    api.NewOptString(vmRequest.DatacenterId),
			CreatedAt:       vmRequest.CreatedAt,
			RequestMetadata: vmRequest.RequestMetadata,
		}
		if vmRequest.CompletedAt != nil {
			apiVMRequests[i].CompletedAt = api.NewOptNilDateTime(*vmRequest.CompletedAt)
		}
	}

	// Create the response structure for VM deploy instances
	apiDeployList := make([]api.VMDeployInstance, len(deployInstances))
	for i, inst := range deployInstances {
		apiDeployList[i] = api.VMDeployInstance{
			RequestId:      inst.RequestID,
			VmId:           api.NewOptString(inst.VMID),
			VmName:         inst.VMName,
			VmStatus:       inst.VMStatus,
			VmStateMessage: api.NewOptString(inst.VMStateMessage),
		}
		if inst.CompletedAt != nil {
			apiDeployList[i].CompletedAt = api.NewOptNilDateTime(*inst.CompletedAt)
		}
	}

	return &api.VMRequestsList{
		VMMinusRequestsListCount: api.NewOptInt(reqCount),
		VMMinusDeployListCount:   api.NewOptInt(instCount),
		Items: api.NewOptVMRequestsListItems(api.VMRequestsListItems{
			VMRequetsList: apiVMRequests,
			VMDeployList:  apiDeployList,
		}),
	}, nil
}

// validateImage checks if an image exists and returns its path.
func (h *Handler) validateImage(ctx context.Context, imageID string) (string, error) {
	if !h.deps.Config.App.Application.ValidateClientRequest {
		return "", nil
	}

	imageClient := h.deps.ClientDependency.ImageManagerClient
	res, err := imageClient.GetAvailableImages(ctx)
	if err != nil {
		h.deps.Logger.Errorf("Failed to get image %s: %v", imageID, err)
		return "", err
	}

	var matchedImage *imagemanager.HypervisorImage
	for _, img := range res {
		if img.ID == imageID {
			matchedImage = &img
			break
		}
	}

	if matchedImage == nil {
		h.deps.Logger.Warnf("Image with ID %s not found", imageID)
		return "", errors.New("image not found for validation")
	}

	h.deps.Logger.Infof("Successfully validated image %s, response: %+v", imageID, matchedImage)
	return matchedImage.Filename, nil

}

// validateHost checks if a host and cluster are active.
func (h *Handler) validateHost(ctx context.Context, hostID, clusterID string) error {
	if !h.deps.Config.App.Application.ValidateClientRequest {
		return nil
	}

	infraClient := h.deps.ClientDependency.InfraMonitorClient

	// Validate host
	hostRes, err := infraClient.GetHypervisorHosts(ctx)
	if err != nil {
		h.deps.Logger.Errorf("Failed to get host %s: %v", hostID, err)
		return err
	}

	var matchedHost *inframonitor.HypervisorHost
	for _, host := range hostRes {
		if host.HostName == hostID {
			matchedHost = &host
			break
		}
	}

	if matchedHost == nil {
		h.deps.Logger.Warnf("host with ID %s not found", hostID)
		return errors.New("host not found for validation")
	}
	if strings.EqualFold(matchedHost.Status, "OK") {
		h.deps.Logger.Warnf("host status %s", matchedHost.Status)
		return errors.New("host is not active")
	}
	h.deps.Logger.Infof("Successfully validated host %s", hostID)

	// Validate cluster
	clusterRes, clustererr := infraClient.GetHypervisorClusters(ctx)
	if clustererr != nil {
		h.deps.Logger.Errorf("Failed to get cluster %s: %v", clusterID, clustererr)
		return clustererr
	}
	var matchedcluster *inframonitor.HypervisorCluster
	for _, cluster := range clusterRes {
		if cluster.ClusterName == clusterID {
			matchedcluster = &cluster
			break
		}
	}

	if matchedcluster == nil {
		h.deps.Logger.Warnf("Cluster with ID %s not found", clusterID)
		return errors.New("Cluster not found for validation")
	}
	if strings.EqualFold(matchedcluster.Status, "OK") {
		h.deps.Logger.Warnf("Cluster status %s", matchedcluster.Status)
		return errors.New("Cluster is not active")
	}
	h.deps.Logger.Infof("Successfully validated cluster %s", clusterID)

	return nil
}

// validateVMExists checks if a VM exists using the vm_monitor client.
func (h *Handler) validateVMExists(ctx context.Context, vmID string, vmOperation constants.OperationType) error {
	if !h.deps.Config.App.Application.ValidateClientRequest {
		h.deps.Logger.Infof("validate client request", h.deps.Config.App.Application.ValidateClientRequest)
		return nil
	}

	vmClient := h.deps.ClientDependency.VmMonitorClient
	timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	res, err := vmClient.GetVmMetrics(timeoutCtx, vmmonitor.GetVmMetricsParams{VMID: vmID})
	// This prevents transient issues with the vm_monitor from blocking the request.
	if err != nil {
		h.deps.Logger.Errorf("Error validating VM %s: %v", vmID, err)
		return err
	}
	switch vmOperation {
	case constants.VMReconfigure:
		h.deps.Logger.Warnf("VM status: %s", res.Powerstate)
		if strings.EqualFold(string(constants.OperationType(res.Powerstate)), string(constants.VMPowerOff)) {
			h.deps.Logger.Warnf("VM %s is powered off and cannot be reconfigured", vmID)
			return errors.New("VM is powered off and cannot be reconfigured")
		}
	}

	h.deps.Logger.Infof("Successfully validated VM %s, response: %+v", vmID, res)
	return nil

}
