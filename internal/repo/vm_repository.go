package repo

import (
	"context"
	"vm/internal/modals"
	"vm/pkg/cinterface"
	"vm/pkg/constants"
	"vm/pkg/db"
)

// VMRepository defines the interface for VM database operations.
type VMRepository interface {
	CreateVMRequest(ctx context.Context, req *modals.VMRequest) error
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
