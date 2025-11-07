package handler_impl

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	imagemanager "vm/internal/client/image_manager"
	inframonitor "vm/internal/client/infra_monitor"
	vmmonitor "vm/internal/client/vm_monitor"
	dto "vm/internal/dtos"
	api "vm/internal/gen"
	"vm/internal/service"
	"vm/pkg/constants"
	"vm/pkg/dependency"
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
		res := constants.MapServiceError(*err, constants.VMReconfigure, ctx)
		return res.(api.EditVMRes), nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal EditVm Request: %v", err)
		res := constants.MapServiceError(dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   "Failed to marshal EditVm Request",
		}, constants.VMReconfigure, ctx)
		return res.(api.EditVMRes), nil
	}
	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMReconfigure, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to marshal EditVm Request: %v", vmRequesterr)
		res := constants.MapServiceError(*vmRequesterr, constants.VMReconfigure, ctx)
		return res.(api.EditVMRes), nil
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
		res := constants.MapServiceError(*err, constants.VMDeploy, ctx)
		return res.(api.HCIDeployVMRes), nil
	}
	req.ImageSource.Value.ImageName = api.NewOptString(imagePath)

	// Validate host and cluster
	if err := h.validateHost(ctx, req.Destination.Value.HostId.Value, req.Destination.Value.ClusterId.Value); err != nil {
		res := constants.MapServiceError(*err, constants.VMDeploy, ctx)
		return res.(api.HCIDeployVMRes), nil
	}

	// Marshal the request to JSON to store as metadata
	metadata, metadataerr := json.Marshal(req)
	if metadataerr != nil {
		h.deps.Logger.Errorf("Failed to marshal HCIDeployVM Request: %v", metadataerr)
		// Handle marshaling error
		res := constants.MapServiceError(dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   "Failed to marshal HCIDeployVM Request",
		}, constants.VMDeploy, ctx)
		return res.(api.HCIDeployVMRes), nil
	}

	// Call the service to create the VM request
	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMDeploy, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VM Deploy request: %v", vmRequesterr)
		res := constants.MapServiceError(*vmRequesterr, constants.VMDeploy, ctx)
		return res.(api.HCIDeployVMRes), nil
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
		res := constants.MapServiceError(*err, constants.VMDelete, ctx)
		return res.(api.VMDeleteRes), nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal VMDelete Request: %v", err)
		res := constants.MapServiceError(dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   "Failed to marshal VMDelete Request",
		}, constants.VMDelete, ctx)
		return res.(api.VMDeleteRes), nil
	}
	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMDelete, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VMDelete Request: %v", vmRequesterr)
		res := constants.MapServiceError(*vmRequesterr, constants.VMDelete, ctx)
		return res.(api.VMDeleteRes), nil
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
		res := constants.MapServiceError(*err, constants.VMPowerOff, ctx)
		return res.(api.VMPowerOffRes), nil

	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal VMPowerOff Request: %v", err)
		res := constants.MapServiceError(dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   err.Error(),
		}, constants.VMPowerOff, ctx)
		return res.(api.VMPowerOffRes), nil
	}
	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMPowerOff, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VMPowerOff Request: %v", vmRequesterr)
		res := constants.MapServiceError(*vmRequesterr, constants.VMPowerOff, ctx)
		return res.(api.VMPowerOffRes), nil
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
		res := constants.MapServiceError(*err, constants.VMPowerOn, ctx)
		return res.(api.VMPowerOnRes), nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal VMPowerOn params: %v", err)
		res := constants.MapServiceError(dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   "Failed to marshal VMPowerOn params",
		}, constants.VMPowerOn, ctx)
		return res.(api.VMPowerOnRes), nil
	}

	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMPowerOn, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VM power on request: %v", vmRequesterr)
		res := constants.MapServiceError(*vmRequesterr, constants.VMPowerOn, ctx)
		return res.(api.VMPowerOnRes), nil
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
		res := constants.MapServiceError(*err, constants.VMReset, ctx)
		return res.(api.VMPowerResetRes), nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal VMPowerReset params: %v", err)
		res := constants.MapServiceError(dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   "Failed to marshal VMPowerReset params",
		}, constants.VMReset, ctx)
		return res.(api.VMPowerResetRes), nil
	}

	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMReset, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VM power reset request: %v", vmRequesterr)
		res := constants.MapServiceError(*vmRequesterr, constants.VMReset, ctx)
		return res.(api.VMPowerResetRes), nil
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
		res := constants.MapServiceError(*err, constants.VMRefresh, ctx)
		return res.(api.VMRefreshRes), nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal VMRefresh params: %v", err)
		res := constants.MapServiceError(dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   "Failed to marshal VMRefresh params",
		}, constants.VMRefresh, ctx)
		return res.(api.VMRefreshRes), nil
	}

	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMRefresh, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VM refresh request: %v", vmRequesterr)
		res := constants.MapServiceError(*vmRequesterr, constants.VMRefresh, ctx)
		return res.(api.VMRefreshRes), nil
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
		res := constants.MapServiceError(*err, constants.VMRestartGuestOS, ctx)
		return res.(api.VMRestartGuestOSRes), nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal VMRestartGuestOS params: %v", err)
		res := constants.MapServiceError(dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   "Failed to marshal VMRestartGuestOS params",
		}, constants.VMRestartGuestOS, ctx)
		return res.(api.VMRestartGuestOSRes), nil
	}

	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMRestartGuestOS, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VM restart guest OS request: %v", vmRequesterr)
		res := constants.MapServiceError(*vmRequesterr, constants.VMRestartGuestOS, ctx)
		return res.(api.VMRestartGuestOSRes), nil
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
		res := constants.MapServiceError(*err, constants.VMShutdownGuestOS, ctx)
		return res.(api.VMShutdownGuestOSRes), nil
	}

	metadata, err := json.Marshal(params)
	if err != nil {
		h.deps.Logger.Errorf("Failed to marshal VMShutdownGuestOS params: %v", err)
		res := constants.MapServiceError(dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   "Failed to marshal VMShutdownGuestOS params",
		}, constants.VMShutdownGuestOS, ctx)
		return res.(api.VMShutdownGuestOSRes), nil
	}

	vmRequest, vmRequesterr := h.VMService.CreateVMRequest(ctx, constants.VMShutdownGuestOS, constants.StatusNew, string(metadata))
	if vmRequesterr != nil {
		h.deps.Logger.Errorf("Failed to create VM shutdown guest OS request: %v", vmRequesterr)
		res := constants.MapServiceError(*vmRequesterr, constants.VMShutdownGuestOS, ctx)
		return res.(api.VMShutdownGuestOSRes), nil
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
		h.deps.Logger.Errorf("Failed to get VM request: %v", err)
		res := constants.MapServiceError(*err, constants.VMMachine, ctx)
		return res.(api.GetVirtualMachineRequestRes), nil
	}

	deployInstances, deployInstanceserr := h.VMService.GetVMDeployInstances(ctx, params.RequestID)
	if deployInstanceserr != nil {
		h.deps.Logger.Warnf("Failed to get VM deploy instances, but continuing execution as they are optional: %v", deployInstanceserr)
		res := constants.MapServiceError(*deployInstanceserr, constants.VMMachine, ctx)
		return res.(api.GetVirtualMachineRequestRes), nil
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
		res := constants.MapServiceError(*err, constants.VMMachineList, ctx)
		return res.(api.GetVirtualMachineRequestListRes), nil
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
func (h *Handler) validateImage(ctx context.Context, imageID string) (string, *dto.ApiResponseError) {
	if !h.deps.Config.App.Application.ValidateClientRequest {
		return "", nil
	}

	imageClient := h.deps.ClientDependency.ImageManagerClient
	res, err := imageClient.GetAvailableImages(ctx)
	if err != nil {
		h.deps.Logger.Errorf("Failed to get image %s: %v", imageID, err)
		return "", &dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   err.Error(),
		}
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
		return "", &dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   "image not found for validation",
		}
	}

	h.deps.Logger.Infof("Successfully validated image %s, response: %+v", imageID, matchedImage)
	return matchedImage.ImageURL, nil
}

// validateHost checks if a host and cluster are active.
func (h *Handler) validateHost(ctx context.Context, hostID, clusterID string) *dto.ApiResponseError {
	if !h.deps.Config.App.Application.ValidateClientRequest {
		return nil
	}

	infraClient := h.deps.ClientDependency.InfraMonitorClient

	// Validate host
	hostRes, err := infraClient.GetHypervisorHosts(ctx)
	if err != nil {
		h.deps.Logger.Errorf("Failed to get host %s: %v", hostID, err)
		return &dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   err.Error(),
		}
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
		return &dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   "host not found for validation",
		}
	}
	if strings.EqualFold(matchedHost.Status, "OK") {
		h.deps.Logger.Warnf("host status %s", matchedHost.Status)
		return &dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   "host is not active",
		}
	}
	h.deps.Logger.Infof("Successfully validated host %s", hostID)

	// Validate cluster
	clusterRes, clustererr := infraClient.GetHypervisorClusters(ctx)
	if clustererr != nil {
		h.deps.Logger.Errorf("Failed to get cluster %s: %v", clusterID, clustererr)
		return &dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   clustererr.Error(),
		}
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
		return &dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   "Cluster not found for validation",
		}
	}
	if strings.EqualFold(matchedcluster.Status, "OK") {
		h.deps.Logger.Warnf("Cluster status %s", matchedcluster.Status)
		return &dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   "Cluster is not active",
		}
	}
	h.deps.Logger.Infof("Successfully validated cluster %s", clusterID)

	return nil
}

// validateVMExists checks if a VM exists using the vm_monitor client.
func (h *Handler) validateVMExists(ctx context.Context, vmID string, vmOperation constants.OperationType) *dto.ApiResponseError {
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
		return &dto.ApiResponseError{
			ErrorCode: constants.InternalServerErrorCode,
			Message:   err.Error(),
		}
	}
	switch vmOperation {
	case constants.VMReconfigure:
		h.deps.Logger.Warnf("VM status: %s", res.Powerstate)
		if strings.EqualFold(string(constants.OperationType(res.Powerstate)), string(constants.VMPowerOff)) {
			h.deps.Logger.Warnf("VM %s is powered off and cannot be reconfigured", vmID)
			return &dto.ApiResponseError{
				ErrorCode: constants.InternalServerErrorCode,
				Message:   "VM is powered off and cannot be reconfigured",
			}
		}
	}

	h.deps.Logger.Infof("Successfully validated VM %s, response: %+v", vmID, res)
	return nil

}
