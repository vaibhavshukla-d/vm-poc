package repo

import (
	"context"
	"vm/internal/modals"
	"vm/pkg/cinterface"
	"vm/pkg/constants"
	"vm/pkg/db"
)


//go:generate mockgen -source=vm_repository.go -destination=mock/vm_repositoryMock.go
type VMRepository interface {
	CreateVMRequest(ctx context.Context, req *modals.VMRequest) error
	GetVMRequest(ctx context.Context, requestID string) (*modals.VMRequest, error)
	GetVMDeployInstances(ctx context.Context, requestID string) ([]*modals.VMDeployInstance, error)
	CreateVMDeployInstances(ctx context.Context, instances []modals.VMDeployInstance) error
	GetAllVMRequestsWithInstances(ctx context.Context) ([]*modals.VMRequest, []*modals.VMDeployInstance, error)
}

// vmRepository implements the VMRepository interface.
type vmRepository struct {
	db     db.Database
	logger cinterface.Logger
}

// NewVMRepository creates a new VMRepository.
func NewVMRepository(db db.Database, logger cinterface.Logger) VMRepository {
	return &vmRepository{
		db:     db,
		logger: logger,
	}
}

// CreateVMRequest creates a new VMRequest record in the database.
func (r *vmRepository) CreateVMRequest(ctx context.Context, req *modals.VMRequest) error {
	r.logger.Info(constants.MySql, constants.Insert, "CreateVMRequest repository function invoked", nil)
	db := r.db.GetReader()

	result := db.WithContext(ctx).Create(req)
	if result.Error != nil {
		r.logger.Error(constants.MySql, constants.Insert, "Failed to create VMRequest", map[constants.ExtraKey]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}

	r.logger.Info(constants.MySql, constants.Insert, "VMRequest created successfully", map[constants.ExtraKey]interface{}{
		"requestID": req.RequestID,
	})

	return nil
}

// GetVMRequest retrieves a VMRequest record from the database by its ID.
func (r *vmRepository) GetVMRequest(ctx context.Context, requestID string) (*modals.VMRequest, error) {
	r.logger.Info(constants.MySql, constants.Select, "GetVMRequest repository function invoked", nil)
	db := r.db.GetReader()

	var req modals.VMRequest
	result := db.WithContext(ctx).Where("request_id = ?", requestID).First(&req)
	if result.Error != nil {
		r.logger.Error(constants.MySql, constants.Select, "Failed to get VMRequest", map[constants.ExtraKey]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}

	r.logger.Info(constants.MySql, constants.Select, "VMRequest retrieved successfully", map[constants.ExtraKey]interface{}{
		"requestID": req.RequestID,
	})

	return &req, nil
}

// GetVMDeployInstances retrieves all VMDeployInstance records from the database by request ID.
func (r *vmRepository) GetVMDeployInstances(ctx context.Context, requestID string) ([]*modals.VMDeployInstance, error) {
	r.logger.Info(constants.MySql, constants.Select, "GetVMDeployInstances repository function invoked", nil)
	db := r.db.GetReader()

	var instances []*modals.VMDeployInstance
	result := db.WithContext(ctx).Where("request_id = ?", requestID).Find(&instances)
	if result.Error != nil {
		r.logger.Error(constants.MySql, constants.Select, "Failed to get VMDeployInstances", map[constants.ExtraKey]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}

	r.logger.Info(constants.MySql, constants.Select, "VMDeployInstances retrieved successfully", map[constants.ExtraKey]interface{}{
		"requestID": requestID,
		"count":     len(instances),
	})

	return instances, nil
}

func (r *vmRepository) CreateVMDeployInstances(ctx context.Context, instances []modals.VMDeployInstance) error {
	r.logger.Info(constants.MySql, constants.Insert, "CreateVMDeployInstances repository function invoked", map[constants.ExtraKey]interface{}{
		"requestID": instances[0].RequestID,
		"count":     len(instances),
	})

	db := r.db.GetReader()
	result := db.WithContext(ctx).Create(&instances)
	if result.Error != nil {
		r.logger.Error(constants.MySql, constants.Insert, "Failed to create VMDeployInstances", map[constants.ExtraKey]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}

	r.logger.Info(constants.MySql, constants.Insert, "VMDeployInstances created successfully", map[constants.ExtraKey]interface{}{
		"count": len(instances),
	})

	return nil
}

func (r *vmRepository) GetAllVMRequestsWithInstances(ctx context.Context) ([]*modals.VMRequest, []*modals.VMDeployInstance, error) {
	var requests []*modals.VMRequest
	db := r.db.GetReader()

	if err := db.WithContext(ctx).Find(&requests).Error; err != nil {
		return nil, nil, err
	}

	var instances []*modals.VMDeployInstance
	if err := db.WithContext(ctx).Find(&instances).Error; err != nil {
		return nil, nil, err
	}

	return requests, instances, nil
}
