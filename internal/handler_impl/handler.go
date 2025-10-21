package handler_impl

import (
	"context"
	"encoding/json"
	"time"
	imagemanager "vm/internal/client/image_manager"
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
	// TODO: Implement edit VM logic
	panic("not implemented")
}

// HCIDeployVM implements the HCIDeployVM operation
func (h *Handler) HCIDeployVM(ctx context.Context, req *api.HCIDeployVM) (api.HCIDeployVMRes, error) {

	h.deps.Logger.Infof("HCIDeployVM handler invoked")

	// Image Manager Ogen client
	imageClient := h.deps.ClientDependency.ImageManagerClient

	// Create a new context with a 10-millisecond timeout.
	timeoutCtx, cancel := context.WithTimeout(h.deps.Ctx, 10*time.Millisecond)
	defer cancel()

	// Attempt to get the image, expecting a timeout.
	res, err := imageClient.GetImage(timeoutCtx, imagemanager.GetImageParams{ImageID: req.ImageSource.Value.ImageId.Value})
	if err != nil {
		// Check if the error is a timeout.
		if timeoutCtx.Err() == context.DeadlineExceeded {
			h.deps.Logger.Warnf("Timeout occurred while getting image %s, but continuing execution as expected.", req.ImageSource.Value.ImageId.Value)
		} else {
			h.deps.Logger.Errorf("Failed to get image %s: %v", req.ImageSource.Value.ImageId.Value, err)
		}
	} else {
		h.deps.Logger.Infof("Successfully validated image %s, response: %+v", req.ImageSource.Value.ImageId.Value, res)
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

	//TODO Implement image ID validation
	// res, err := imageClient.GetImage(ctx, client.GetImageParams{ImageID: req.ImageSource.Value.ImageId.Value})
	// if err != nil {
	// 	h.deps.Logger.Errorf("Failed to get image %s: %v", req.ImageSource.Value.ImageId.Value, err)
	// } else {
	// 	h.deps.Logger.Infof("Successfully validated image %s, response: %+v", req.ImageSource.Value.ImageId.Value, res)
	// }

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
	// TODO: Implement get virtual machine request logic
	panic("not implemented")
}
