package handler_impl

import (
	"context"
	"encoding/json"
	client_gen "vm/internal/client_gen"
	api "vm/internal/gen"
	"vm/internal/modals"
	"vm/internal/service"
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

	// Image Manager Ogen client
	imageClient := h.deps.ClientDependency.ImageManagerClient

	// Marshal the request to JSON to store as metadata
	metadata, err := json.Marshal(req)
	if err != nil {
		// Handle marshaling error
		return &api.HCIDeployVMInternalServerError{
			Message: "Failed to process request",
		}, nil
	}

	//TODO Implement image ID validation
	res, err := imageClient.GetImage(ctx, client_gen.GetImageParams{ImageID: req.ImageSource.Value.ImageId.Value})
	if err != nil {
		h.deps.Logger.Errorf("Failed to get image %s: %v", req.ImageSource.Value.ImageId.Value, err)
	} else {
		h.deps.Logger.Infof("Successfully validated image %s, response: %+v", req.ImageSource.Value.ImageId.Value, res)
	}

	// Create a new VMRequest modal
	vmRequest := &modals.VMRequest{
		Operation:       "VMDeploy",
		RequestStatus:   "NEW",
		RequestMetadata: string(metadata),
	}

	// Call the service to create the VM request
	if err = h.VMService.DeployVM(ctx, vmRequest); err != nil {
		// Handle service error
		return &api.HCIDeployVMInternalServerError{
			Message: "Failed to create VM request",
		}, nil
	}

	// The vmRequest.RequestID is now populated by the BeforeCreate hook.
	location := "/virtualization/v1beta1/virtual-machines-request/" + vmRequest.RequestID

	// Return the async response with the Location header
	return &api.EmptyResponseHeaders{
		Location: api.NewOptString(location),
		Response: api.EmptyResponse{},
	}, nil
}

// VMDelete implements the VMDelete operation
func (h *Handler) VMDelete(ctx context.Context, params api.VMDeleteParams) (api.VMDeleteRes, error) {
	// TODO: Implement VM delete logic
	panic("not implemented")
}

// VMPowerOff implements the VMPowerOff operation
func (h *Handler) VMPowerOff(ctx context.Context, params api.VMPowerOffParams) (api.VMPowerOffRes, error) {
	// TODO: Implement VM power off logic
	panic("not implemented")
}

// VMPowerOn implements the VMPowerOn operation
func (h *Handler) VMPowerOn(ctx context.Context, params api.VMPowerOnParams) (api.VMPowerOnRes, error) {
	// TODO: Implement VM power on logic
	panic("not implemented")
}

// VMPowerReset implements the VMPowerReset operation
func (h *Handler) VMPowerReset(ctx context.Context, params api.VMPowerResetParams) (api.VMPowerResetRes, error) {
	// TODO: Implement VM power reset logic
	panic("not implemented")
}

// VMRefresh implements the VMRefresh operation
func (h *Handler) VMRefresh(ctx context.Context, params api.VMRefreshParams) (api.VMRefreshRes, error) {
	// TODO: Implement VM refresh logic
	panic("not implemented")
}

// VMRestartGuestOS implements the VMRestartGuestOS operation
func (h *Handler) VMRestartGuestOS(ctx context.Context, params api.VMRestartGuestOSParams) (api.VMRestartGuestOSRes, error) {
	// TODO: Implement VM restart guest OS logic
	panic("not implemented")
}

// VMShutdownGuestOS implements the VMShutdownGuestOS operation
func (h *Handler) VMShutdownGuestOS(ctx context.Context, params api.VMShutdownGuestOSParams) (api.VMShutdownGuestOSRes, error) {
	// TODO: Implement VM shutdown guest OS logic
	panic("not implemented")
}

// GetVirtualMachineRequest implements the GetVirtualMachineRequest operation
func (h *Handler) GetVirtualMachineRequest(ctx context.Context, params api.GetVirtualMachineRequestParams) (api.GetVirtualMachineRequestRes, error) {
	// TODO: Implement get virtual machine request logic
	panic("not implemented")
}
