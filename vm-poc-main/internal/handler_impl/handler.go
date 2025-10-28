package handler_impl
 
import (
    "context"
    "encoding/json"
    "errors"
    "net/url"
    "time"
 
    imagemanager "vm/internal/client/image_manager"
    inframonitor "vm/internal/client/infra_monitor"
    vmmonitor "vm/internal/client/vm_monitor"
    api "vm/internal/gen"
    "vm/internal/modals"
    "vm/internal/service"
    "vm/pkg/constants"
    "vm/pkg/dependency"
 
    "github.com/google/uuid"
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
    // TODO: Implement edit VM logic
 
    if err := h.validateVMExists(ctx, string(params.VMID)); err != nil {
        return &api.EditVMInternalServerError{
            Message: err.Error(),
        }, nil
    }
 
    metadata, err := json.Marshal(params)
    if err != nil {
        h.deps.Logger.Errorf("Failed to marshal EditVm Request: %v", err)
        return &api.EditVMInternalServerError{
            Message: "Failed to process request",
        }, nil
    }
    vmRequest, err := h.VMService.CreateVMRequest(ctx, constants.VMReconfigure, constants.StatusNew, string(metadata))
    if err != nil {
        h.deps.Logger.Errorf("Failed to marshal EditVm Request: %v", err)
        return &api.EditVMInternalServerError{
            Message: "Failed to create VM request",
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
    metadata, err := json.Marshal(req)
    if err != nil {
        h.deps.Logger.Errorf("Failed to marshal HCIDeployVM Request: %v", err)
        // Handle marshaling error
        return &api.HCIDeployVMInternalServerError{
            Message: "Failed to process request",
        }, nil
    }
 
    // Call the service to create the VM request
    vmRequest, err := h.VMService.CreateVMRequest(ctx, constants.VMDeploy, constants.StatusNew, string(metadata))
    if err != nil {
        h.deps.Logger.Errorf("Failed to create VM Deploy request: %v", err)
        return &api.HCIDeployVMInternalServerError{
            Message: "Failed to create VM request",
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
 
    if err := h.validateVMExists(ctx, string(params.VMID)); err != nil {
        return &api.VMDeleteInternalServerError{
            Message: err.Error(),
        }, nil
    }
 
    metadata, err := json.Marshal(params)
    if err != nil {
        h.deps.Logger.Errorf("Failed to marshal VMDelete Request: %v", err)
        return &api.VMDeleteInternalServerError{
            Message: "Failed to process request",
        }, nil
    }
    vmRequest, err := h.VMService.CreateVMRequest(ctx, constants.VMDelete, constants.StatusNew, string(metadata))
    if err != nil {
        h.deps.Logger.Errorf("Failed to marshal VMDelete Request: %v", err)
        return &api.VMDeleteInternalServerError{
            Message: "Failed to create VM request",
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
 
    if err := h.validateVMExists(ctx, string(params.VMID)); err != nil {
        return &api.VMPowerOffInternalServerError{
            Message: err.Error(),
        }, nil
    }
 
    metadata, err := json.Marshal(params)
    if err != nil {
        h.deps.Logger.Errorf("Failed to marshal VMPowerOff Request: %v", err)
        return &api.VMPowerOffInternalServerError{
            Message: "Failed to process request",
        }, nil
    }
    vmRequest, err := h.VMService.CreateVMRequest(ctx, constants.VMPowerOff, constants.StatusNew, string(metadata))
    if err != nil {
        h.deps.Logger.Errorf("Failed to marshal VMPowerOff Request: %v", err)
        return &api.VMPowerOffInternalServerError{
            Message: "Failed to create VM request",
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
 
    if err := h.validateVMExists(ctx, string(params.VMID)); err != nil {
        return &api.VMPowerOnInternalServerError{
            Message: err.Error(),
        }, nil
    }
 
    metadata, err := json.Marshal(params)
    if err != nil {
        h.deps.Logger.Errorf("Failed to marshal VMPowerOn params: %v", err)
        return &api.VMPowerOnInternalServerError{
            Message: "Failed to process request",
        }, nil
    }
 
    vmRequest, err := h.VMService.CreateVMRequest(ctx, constants.VMPowerOn, constants.StatusNew, string(metadata))
    if err != nil {
        h.deps.Logger.Errorf("Failed to create VM power on request: %v", err)
        return &api.VMPowerOnInternalServerError{
            Message: "Failed to create VM request",
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
 
    if err := h.validateVMExists(ctx, string(params.VMID)); err != nil {
        return &api.VMPowerResetInternalServerError{
            Message: err.Error(),
        }, nil
    }
 
    metadata, err := json.Marshal(params)
    if err != nil {
        h.deps.Logger.Errorf("Failed to marshal VMPowerReset params: %v", err)
        return &api.VMPowerResetInternalServerError{
            Message: "Failed to process request",
        }, nil
    }
 
    vmRequest, err := h.VMService.CreateVMRequest(ctx, constants.VMReset, constants.StatusNew, string(metadata))
    if err != nil {
        h.deps.Logger.Errorf("Failed to create VM power reset request: %v", err)
        return &api.VMPowerResetInternalServerError{
            Message: "Failed to create VM request",
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
 
    if err := h.validateVMExists(ctx, string(params.VMID)); err != nil {
        return &api.VMRefreshInternalServerError{
            Message: err.Error(),
        }, nil
    }
 
    metadata, err := json.Marshal(params)
    if err != nil {
        h.deps.Logger.Errorf("Failed to marshal VMRefresh params: %v", err)
        return &api.VMRefreshInternalServerError{
            Message: "Failed to process request",
        }, nil
    }
 
    vmRequest, err := h.VMService.CreateVMRequest(ctx, constants.VMRefresh, constants.StatusNew, string(metadata))
    if err != nil {
        h.deps.Logger.Errorf("Failed to create VM refresh request: %v", err)
        return &api.VMRefreshInternalServerError{
            Message: "Failed to create VM request",
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
 
    if err := h.validateVMExists(ctx, string(params.VMID)); err != nil {
        return &api.VMRestartGuestOSInternalServerError{
            Message: err.Error(),
        }, nil
    }
 
    metadata, err := json.Marshal(params)
    if err != nil {
        h.deps.Logger.Errorf("Failed to marshal VMRestartGuestOS params: %v", err)
        return &api.VMRestartGuestOSInternalServerError{
            Message: "Failed to process request",
        }, nil
    }
 
    vmRequest, err := h.VMService.CreateVMRequest(ctx, constants.VMRestartGuestOS, constants.StatusNew, string(metadata))
    if err != nil {
        h.deps.Logger.Errorf("Failed to create VM restart guest OS request: %v", err)
        return &api.VMRestartGuestOSInternalServerError{
            Message: "Failed to create VM request",
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
 
    if err := h.validateVMExists(ctx, string(params.VMID)); err != nil {
        return &api.VMShutdownGuestOSInternalServerError{
            Message: err.Error(),
        }, nil
    }
 
    metadata, err := json.Marshal(params)
    if err != nil {
        h.deps.Logger.Errorf("Failed to marshal VMShutdownGuestOS params: %v", err)
        return &api.VMShutdownGuestOSInternalServerError{
            Message: "Failed to process request",
        }, nil
    }
 
    vmRequest, err := h.VMService.CreateVMRequest(ctx, constants.VMShutdownGuestOS, constants.StatusNew, string(metadata))
    if err != nil {
        h.deps.Logger.Errorf("Failed to create VM shutdown guest OS request: %v", err)
        return &api.VMShutdownGuestOSInternalServerError{
            Message: "Failed to create VM request",
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
 
    requestIDStr := params.RequestID.String()
 
    vmRequest, err := h.VMService.GetVMRequest(ctx, requestIDStr)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return &api.GetVirtualMachineRequestNotFound{
                Message: "VM request not found",
            }, nil
        }
        h.deps.Logger.Errorf("Failed to get VM request: %v", err)
        return &api.GetVirtualMachineRequestInternalServerError{
            Message: "Failed to get VM request",
        }, nil
    }
 
    deployInstances, err := h.VMService.GetVMDeployInstances(ctx, requestIDStr)
    if err != nil {
        h.deps.Logger.Warnf("Failed to get VM deploy instances, but continuing execution as they are optional: %v", err)
        deployInstances = []*modals.VMDeployInstance{} // Initialize with an empty slice
    }
 
    requestID, err := uuid.Parse(vmRequest.RequestID)
    if err != nil {
        h.deps.Logger.Errorf("Failed to parse request ID: %v", err)
        return &api.GetVirtualMachineRequestInternalServerError{
            Message: "Failed to parse request ID",
        }, nil
    }
 
    apiVMRequest := api.VMRequest{
        RequestId:       requestID,
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
        instanceRequestID, err := uuid.Parse(inst.RequestID)
        if err != nil {
            h.deps.Logger.Errorf("Failed to parse instance request ID: %v", err)
            return &api.GetVirtualMachineRequestInternalServerError{
                Message: "Failed to parse instance request ID",
            }, nil
        }
        apiDeployList[i] = api.VMDeployInstance{
            RequestId:      instanceRequestID,
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
            Message: "Failed to get VM request list",
        }, nil
    }
 
    // Create the response structure for VM requests
    apiVMRequests := make([]api.VMRequest, len(vmRequests))
    for i, vmRequest := range vmRequests {
        requestID, err := uuid.Parse(vmRequest.RequestID)
        if err != nil {
            h.deps.Logger.Errorf("Failed to parse request ID: %v", err)
            return &api.GetVirtualMachineRequestListInternalServerError{
                Message: "Failed to parse request ID",
            }, nil
        }
 
        apiVMRequests[i] = api.VMRequest{
            RequestId:       requestID,
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
        instanceRequestID, err := uuid.Parse(inst.RequestID)
        if err != nil {
            h.deps.Logger.Errorf("Failed to parse instance request ID: %v", err)
            return &api.GetVirtualMachineRequestListInternalServerError{
                Message: "Failed to parse instance request ID",
            }, nil
        }
        apiDeployList[i] = api.VMDeployInstance{
            RequestId:      instanceRequestID,
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
    res, err := imageClient.GetImage(ctx, imagemanager.GetImageParams{ImageID: imageID})
    if err != nil {
        h.deps.Logger.Errorf("Failed to get image %s: %v", imageID, err)
        return "", errors.New("failed to validate image")
    }
 
    image, ok := res.(*imagemanager.Image)
    if !ok {
        return "", errors.New("unexpected response type for image validation")
    }
 
    h.deps.Logger.Infof("Successfully validated image %s, response: %+v", imageID, image)
    return image.Name, nil
}
 
// validateHost checks if a host and cluster are active.
func (h *Handler) validateHost(ctx context.Context, hostID, clusterID string) error {
    if !h.deps.Config.App.Application.ValidateClientRequest {
        return nil
    }
 
    infraClient := h.deps.ClientDependency.InfraMonitorClient
 
    // Validate host
    hostRes, err := infraClient.HypervisorHost(ctx, inframonitor.HypervisorHostParams{HostID: hostID})
    if err != nil {
        h.deps.Logger.Errorf("Failed to get host %s: %v", hostID, err)
        return errors.New("failed to validate host")
    }
    host, ok := hostRes.(*inframonitor.HypervisorHost)
    if !ok {
        return errors.New("unexpected response type for host validation")
    }
    if host.Status.Value != "OK" {
        return errors.New("host is not active")
    }
    h.deps.Logger.Infof("Successfully validated host %s", hostID)
 
    // Validate cluster
    clusterRes, err := infraClient.HypervisorCluster(ctx, inframonitor.HypervisorClusterParams{ClusterID: clusterID})
    if err != nil {
        h.deps.Logger.Errorf("Failed to get cluster %s: %v", clusterID, err)
        return errors.New("failed to validate cluster")
    }
    cluster, ok := clusterRes.(*inframonitor.HypervisorCluster)
    if !ok {
        return errors.New("unexpected response type for cluster validation")
    }
    if cluster.Status.Value != "OK" {
        return errors.New("cluster is not active")
    }
    h.deps.Logger.Infof("Successfully validated cluster %s", clusterID)
 
    return nil
}
 
// validateVMExists checks if a VM exists using the vm_monitor client.
func (h *Handler) validateVMExists(ctx context.Context, vmID string) error {
    if !h.deps.Config.App.Application.ValidateClientRequest {
        h.deps.Logger.Infof("validate client request", h.deps.Config.App.Application.ValidateClientRequest)
        return nil
    }
 
    vmClient := h.deps.ClientDependency.VmMonitorClient
    timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
    defer cancel()
 
    res, err := vmClient.GetVm(timeoutCtx, vmmonitor.GetVmParams{VMID: vmID})
    if err != nil {
        // This prevents transient issues with the vm_monitor from blocking the request.
        var notFoundErr *vmmonitor.GetVmNotFound
        var urlErr *url.Error
 
        if errors.As(err, &notFoundErr) {
            h.deps.Logger.Errorf("VM %s not found in vm-monitor.", vmID)
        } else if errors.As(err, &urlErr) {
            if urlErr.Timeout() {
                h.deps.Logger.Errorf("Connection to vm-monitor timed out for VM %s: %v", vmID, err)
            } else {
                h.deps.Logger.Errorf("Network error while validating VM %s: %v", vmID, err)
            }
        } else {
            h.deps.Logger.Errorf("Failed to validate VM %s: %v", vmID, err)
        }
        return errors.New("vm validation failed")
    }
 
    h.deps.Logger.Infof("Successfully validated VM %s, response: %+v", vmID, res)
    return nil
}