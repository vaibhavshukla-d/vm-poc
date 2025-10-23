package service

import (
	"context"
	"encoding/json"
	api "vm/internal/gen"
	"vm/internal/modals"
	"vm/internal/repo"
	"vm/pkg/cinterface"
	"vm/pkg/constants"
	// utils "vm/pkg/utils"
)

// VMService defines the interface for VM-related business logic.
type VMService interface {
	CreateVMRequest(ctx context.Context, operation constants.OperationType, status constants.RequestStatus, metadata string) (*modals.VMRequest, error)
	GetVMRequest(ctx context.Context, requestID string) (*modals.VMRequest, error)
	GetVMDeployInstances(ctx context.Context, requestID string) ([]*modals.VMDeployInstance, error)
	GetAllVMRequestsWithInstances(ctx context.Context) ([]*modals.VMRequest, []*modals.VMDeployInstance, int, int, error)
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

	// workspaceID, errUtlis := utils.GetWorkspaceIDFromContext(ctx)
	// if errUtlis != nil {
	// 	s.logger.Error(constants.Internal, constants.Api, "Missing or invalid workspace_id in context", map[constants.ExtraKey]interface{}{
	// 		"error": errUtlis.Error(),
	// 	})
	// 	return nil, errUtlis
	// }

	vmRequest := &modals.VMRequest{
		Operation:       string(operation),
		RequestStatus:   string(status),
		RequestMetadata: metadata,
		// WorkspaceId:     workspaceID,	
	}
	err := s.vmRepo.CreateVMRequest(ctx, vmRequest)
	if err != nil {
		s.logger.Error(constants.Internal, constants.Api, "Failed to deploy VM", map[constants.ExtraKey]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	if operation == constants.VMDeploy {
		var deployReq api.HCIDeployVM
		if err := json.Unmarshal([]byte(metadata), &deployReq); err != nil {
			s.logger.Error(constants.Internal, constants.Api, "Failed to unmarshal deploy metadata", map[constants.ExtraKey]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		numVMs := deployReq.VmConfig.NumberOfVms.Value
		vmName := deployReq.VmConfig.Name

		if numVMs > 0 {
			err := s.vmRepo.CreateVMDeployInstances(ctx, vmRequest.RequestID, vmName, numVMs)
			if err != nil {
				s.logger.Error(constants.Internal, constants.Api, "Failed to create VM deploy instances", map[constants.ExtraKey]interface{}{
					"error": err.Error(),
				})
				return nil, err
			}
		}
	}

	s.logger.Info(constants.Internal, constants.Api, "Successfully created VM request", nil)

	return vmRequest, nil
}

// GetVMDeployInstances handles the business logic for retrieving VM deploy instances.
func (s *vmService) GetVMDeployInstances(ctx context.Context, requestID string) ([]*modals.VMDeployInstance, error) {
	s.logger.Info(constants.Internal, constants.Api, "GetVMDeployInstances service function invoked", nil)

	instances, err := s.vmRepo.GetVMDeployInstances(ctx, requestID)
	if err != nil {
		s.logger.Error(constants.Internal, constants.Api, "Failed to get VM deploy instances", map[constants.ExtraKey]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	s.logger.Info(constants.Internal, constants.Api, "Successfully retrieved VM deploy instances", nil)

	return instances, nil
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

func (s *vmService) GetAllVMRequestsWithInstances(ctx context.Context) ([]*modals.VMRequest, []*modals.VMDeployInstance, int, int, error) {
	s.logger.Info(constants.Internal, constants.Api, "GetAllVMRequestsWithInstances service function invoked", nil)

	vmRequests, vmInstances, err := s.vmRepo.GetAllVMRequestsWithInstances(ctx)
	if err != nil {
		s.logger.Error(constants.Internal, constants.Api, "Failed to get all VM requests and instances", map[constants.ExtraKey]interface{}{
			"error": err.Error(),
		})
		return nil, nil, 0, 0, err
	}

	s.logger.Info(constants.Internal, constants.Api, "Successfully retrieved all VM requests and instances", map[constants.ExtraKey]interface{}{
		"request_count":  len(vmRequests),
		"instance_count": len(vmInstances),
	})

	return vmRequests, vmInstances, len(vmRequests), len(vmInstances), nil
}
