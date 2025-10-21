package service

import (
	"context"
	"vm/internal/modals"
	"vm/internal/repo"
	"vm/pkg/cinterface"
	"vm/pkg/constants"
)

// VMService defines the interface for VM-related business logic.
type VMService interface {
	DeployVM(ctx context.Context, req *modals.VMRequest) error
}

// vmService implements the VMService interface.
type vmService struct {
	vmRepo repo.VMRepository
	logger cinterface.Logger
}

// NewVMService creates a new VMService.
func NewVMService(vmRepo repo.VMRepository, logger cinterface.Logger) VMService {
	return &vmService{
		vmRepo: vmRepo,
		logger: logger,
	}
}

// DeployVM handles the business logic for deploying a VM.
func (s *vmService) DeployVM(ctx context.Context, req *modals.VMRequest) error {
	err := s.vmRepo.CreateVMRequest(ctx, req)
	if err != nil {
		s.logger.Error(constants.Internal, constants.Api, "Failed to deploy VM", map[constants.ExtraKey]interface{}{
			"error": err.Error(),
		})
		return err
	}

	s.logger.Info(constants.Internal, constants.Api, "Successfully created VM request", nil)
	return nil
}
