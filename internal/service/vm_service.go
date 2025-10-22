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
	CreateVMRequest(ctx context.Context, operation constants.OperationType, status constants.RequestStatus, metadata string) (*modals.VMRequest, error)
	GetVMRequest(ctx context.Context, requestID string) (*modals.VMRequest, error)
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
func (s *vmService) CreateVMRequest(ctx context.Context, operation constants.OperationType, status constants.RequestStatus, metadata string) (*modals.VMRequest, error) {
	s.logger.Info(constants.Internal, constants.Api, "CreateVMRequest service function invoked", nil)

	s.logger.Info(constants.Internal, constants.Api, "VMRequest payload log", map[constants.ExtraKey]interface{}{
		"operation": operation,
		"status":    status,
		"metadata":  metadata,
	})

	vmRequest := &modals.VMRequest{
		Operation:       string(operation),
		RequestStatus:   string(status),
		RequestMetadata: metadata,
	}

	err := s.vmRepo.CreateVMRequest(ctx, vmRequest)
	if err != nil {
		s.logger.Error(constants.Internal, constants.Api, "Failed to deploy VM", map[constants.ExtraKey]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	s.logger.Info(constants.Internal, constants.Api, "Successfully created VM request", nil)

	return vmRequest, nil
}

// GetVMRequest handles the business logic for retrieving a VM request.
func (s *vmService) GetVMRequest(ctx context.Context, requestID string) (*modals.VMRequest, error) {
	s.logger.Info(constants.Internal, constants.Api, "GetVMRequest service function invoked", nil)

	vmRequest, err := s.vmRepo.GetVMRequest(ctx, requestID)
	if err != nil {
		s.logger.Error(constants.Internal, constants.Api, "Failed to get VM request", map[constants.ExtraKey]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	s.logger.Info(constants.Internal, constants.Api, "Successfully retrieved VM request", nil)

	return vmRequest, nil
}
