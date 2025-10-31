package repo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"vm/internal/modals"
	"vm/internal/repo"
	mock_db "vm/pkg/db/mock"

	mock_logger "vm/pkg/logger/mock"

	"github.com/golang/mock/gomock"
)

func TestCreateVMRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_db.NewMockDatabase(ctrl)
	mockLogger := &mock_logger.StubLogger{}

	ctx := context.Background()

	t.Run("Successful insert", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()

		gormDB, _ := gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		mockDB.EXPECT().GetReader().Return(gormDB.Session(&gorm.Session{SkipHooks: true}))

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `vm_requests`").
			WithArgs("req-123", "vmDeploy", "New", "workspace-001", "dc-001", sqlmock.AnyArg(), nil, `{"key":"value"}`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := repo.NewVMRepository(mockDB, mockLogger)

		req := &modals.VMRequest{
			RequestID:       "req-123",
			Operation:       "vmDeploy",
			RequestStatus:   "New",
			WorkspaceId:     "workspace-001",
			DatacenterId:    "dc-001",
			CreatedAt:       time.Now(),
			CompletedAt:     nil,
			RequestMetadata: `{"key":"value"}`,
		}

		err := repo.CreateVMRequest(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("Insert failure", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()

		gormDB, _ := gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		mockDB.EXPECT().GetReader().Return(gormDB.Session(&gorm.Session{SkipHooks: true}))

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `vm_requests`").
			WithArgs("req-123", "vmDeploy", "New", "workspace-001", "dc-001", sqlmock.AnyArg(), nil, `{"key":"value"}`).
			WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		repo := repo.NewVMRepository(mockDB, mockLogger)

		req := &modals.VMRequest{
			RequestID:       "req-123",
			Operation:       "vmDeploy",
			RequestStatus:   "New",
			WorkspaceId:     "workspace-001",
			DatacenterId:    "dc-001",
			CreatedAt:       time.Now(),
			CompletedAt:     nil,
			RequestMetadata: `{"key":"value"}`,
		}

		err := repo.CreateVMRequest(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, "insert error", err.Error())
	})
}

func TestGetVMRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_db.NewMockDatabase(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	ctx := context.Background()
	requestID := "req-123"

	t.Run("Successful retrieval", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()

		gormDB, _ := gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		mockDB.EXPECT().GetReader().Return(gormDB)

		mock.ExpectQuery("SELECT .* FROM `vm_requests` WHERE request_id = ?").
			WithArgs(requestID, 1).
			WillReturnRows(sqlmock.NewRows([]string{
				"request_id", "operation", "request_status", "workspace_id", "datacenter_id", "created_at", "completed_at", "request_metadata",
			}).AddRow(
				requestID, "vmDeploy", "New", "workspace-001", "dc-001", time.Now(), nil, `{"key":"value"}`,
			))

		repo := repo.NewVMRepository(mockDB, mockLogger)
		result, err := repo.GetVMRequest(ctx, requestID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, requestID, result.RequestID)
	})

	t.Run("Query error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()

		gormDB, _ := gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		mockDB.EXPECT().GetReader().Return(gormDB)

		mock.ExpectQuery("SELECT .* FROM `vm_requests` WHERE request_id = ?").
			WithArgs(requestID, 1).
			WillReturnError(errors.New("query failed"))

		repo := repo.NewVMRepository(mockDB, mockLogger)
		result, err := repo.GetVMRequest(ctx, requestID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetVMDeployInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_db.NewMockDatabase(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	ctx := context.Background()
	requestID := "req-123"

	t.Run("Successful retrieval", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()

		gormDB, _ := gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		mockDB.EXPECT().GetReader().Return(gormDB)

		mock.ExpectQuery("SELECT .* FROM `vm_deploy_instances` WHERE request_id = ?").
			WithArgs(requestID).
			WillReturnRows(sqlmock.NewRows([]string{
				"request_id", "vm_name", "vm_id", "vm_status", "vm_state_message", "completed_at",
			}).AddRow(
				requestID, "vm-1", "vmid-001", "INIT", "Starting up", nil,
			).AddRow(
				requestID, "vm-2", "vmid-002", "RUNNING", "Running smoothly", nil,
			))

		repo := repo.NewVMRepository(mockDB, mockLogger)
		result, err := repo.GetVMDeployInstances(ctx, requestID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "vm-1", result[0].VMName)
		assert.Equal(t, "vm-2", result[1].VMName)
	})

	t.Run("No records found", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()

		gormDB, _ := gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		mockDB.EXPECT().GetReader().Return(gormDB)

		mock.ExpectQuery("SELECT .* FROM `vm_deploy_instances` WHERE request_id = ?").
			WithArgs(requestID).
			WillReturnRows(sqlmock.NewRows([]string{
				"request_id", "vm_name", "vm_id", "vm_status", "vm_state_message", "completed_at",
			})) // no rows

		repo := repo.NewVMRepository(mockDB, mockLogger)
		result, err := repo.GetVMDeployInstances(ctx, requestID)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("Query error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()

		gormDB, _ := gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		mockDB.EXPECT().GetReader().Return(gormDB)

		mock.ExpectQuery("SELECT .* FROM `vm_deploy_instances` WHERE request_id = ?").
			WithArgs(requestID).
			WillReturnError(errors.New("query failed"))

		repo := repo.NewVMRepository(mockDB, mockLogger)
		result, err := repo.GetVMDeployInstances(ctx, requestID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestCreateVMDeployInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_db.NewMockDatabase(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	ctx := context.Background()

	t.Run("Successful insert", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()

		gormDB, _ := gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		mockDB.EXPECT().GetReader().Return(gormDB)

		instances := []modals.VMDeployInstance{
			{
				RequestID:      "req-123",
				VMName:         "vm-1",
				VMID:           "vmid-001",
				VMStatus:       "INIT",
				VMStateMessage: "Starting",
				CompletedAt:    nil,
			},
			{
				RequestID:      "req-123",
				VMName:         "vm-2",
				VMID:           "vmid-002",
				VMStatus:       "RUNNING",
				VMStateMessage: "Running",
				CompletedAt:    nil,
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `vm_deploy_instances`").
			WithArgs(
				instances[0].RequestID, instances[0].VMName, instances[0].VMID, instances[0].VMStatus, instances[0].VMStateMessage, nil,
				instances[1].RequestID, instances[1].VMName, instances[1].VMID, instances[1].VMStatus, instances[1].VMStateMessage, nil,
			).
			WillReturnResult(sqlmock.NewResult(1, 2))
		mock.ExpectCommit()

		repo := repo.NewVMRepository(mockDB, mockLogger)
		err := repo.CreateVMDeployInstances(ctx, instances)

		assert.NoError(t, err)
	})

	t.Run("Insert failure", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()

		gormDB, _ := gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		mockDB.EXPECT().GetReader().Return(gormDB)

		instances := []modals.VMDeployInstance{
			{
				RequestID:      "req-456",
				VMName:         "vm-3",
				VMID:           "vmid-003",
				VMStatus:       "FAILED",
				VMStateMessage: "Error",
				CompletedAt:    nil,
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `vm_deploy_instances`").
			WithArgs(
				instances[0].RequestID, instances[0].VMName, instances[0].VMID, instances[0].VMStatus, instances[0].VMStateMessage, nil,
			).
			WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		repo := repo.NewVMRepository(mockDB, mockLogger)
		err := repo.CreateVMDeployInstances(ctx, instances)

		assert.Error(t, err)
		assert.Equal(t, "insert error", err.Error())
	})
}

func TestGetAllVMRequestsWithInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_db.NewMockDatabase(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	ctx := context.Background()

	t.Run("Successful retrieval", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()

		gormDB, _ := gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		mockDB.EXPECT().GetReader().Return(gormDB).Times(1)

		mock.ExpectQuery("SELECT .* FROM `vm_requests`").
			WillReturnRows(sqlmock.NewRows([]string{
				"request_id", "operation", "request_status", "workspace_id", "datacenter_id", "created_at", "completed_at", "request_metadata",
			}).AddRow(
				"req-001", "vmDeploy", "New", "workspace-001", "dc-001", time.Now(), nil, `{"key":"value"}`,
			))

		mock.ExpectQuery("SELECT .* FROM `vm_deploy_instances`").
			WillReturnRows(sqlmock.NewRows([]string{
				"request_id", "vm_name", "vm_id", "vm_status", "vm_state_message", "completed_at",
			}).AddRow(
				"req-001", "vm-1", "vmid-001", "INIT", "Starting", nil,
			))

		repo := repo.NewVMRepository(mockDB, mockLogger)
		requests, instances, err := repo.GetAllVMRequestsWithInstances(ctx)

		assert.NoError(t, err)
		assert.Len(t, requests, 1)
		assert.Len(t, instances, 1)
		assert.Equal(t, "req-001", requests[0].RequestID)
		assert.Equal(t, "vm-1", instances[0].VMName)
	})

	t.Run("Error fetching requests", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()

		gormDB, _ := gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		mockDB.EXPECT().GetReader().Return(gormDB)

		mock.ExpectQuery("SELECT .* FROM `vm_requests`").
			WillReturnError(errors.New("query error"))

		repo := repo.NewVMRepository(mockDB, mockLogger)
		requests, instances, err := repo.GetAllVMRequestsWithInstances(ctx)

		assert.Error(t, err)
		assert.Nil(t, requests)
		assert.Nil(t, instances)
	})

	t.Run("Error fetching instances", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()

		gormDB, _ := gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		mockDB.EXPECT().GetReader().Return(gormDB).Times(1)

		mock.ExpectQuery("SELECT .* FROM `vm_requests`").
			WillReturnRows(sqlmock.NewRows([]string{
				"request_id", "operation", "request_status", "workspace_id", "datacenter_id", "created_at", "completed_at", "request_metadata",
			}).AddRow(
				"req-001", "vmDeploy", "New", "workspace-001", "dc-001", time.Now(), nil, `{"key":"value"}`,
			))

		mock.ExpectQuery("SELECT .* FROM `vm_deploy_instances`").
			WillReturnError(errors.New("instance query error"))

		repo := repo.NewVMRepository(mockDB, mockLogger)
		requests, instances, err := repo.GetAllVMRequestsWithInstances(ctx)

		assert.Error(t, err)
		assert.Nil(t, requests)
		assert.Nil(t, instances)
	})

	t.Run("RecordNotFound for requests", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()

		gormDB, _ := gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		mockDB.EXPECT().GetReader().Return(gormDB)

		mock.ExpectQuery("SELECT .* FROM `vm_requests`").
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectQuery("SELECT .* FROM `vm_deploy_instances`").
			WillReturnRows(sqlmock.NewRows([]string{
				"request_id", "vm_name", "vm_id", "vm_status", "vm_state_message", "completed_at",
			}).AddRow(
				"req-001", "vm-1", "vmid-001", "INIT", "Starting", nil,
			))

		repo := repo.NewVMRepository(mockDB, mockLogger)
		requests, instances, err := repo.GetAllVMRequestsWithInstances(ctx)

		assert.NoError(t, err)
		assert.Empty(t, requests)
		assert.Len(t, instances, 1)
	})

	t.Run("RecordNotFound for instances", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()

		gormDB, _ := gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		mockDB.EXPECT().GetReader().Return(gormDB).Times(1)

		mock.ExpectQuery("SELECT .* FROM `vm_requests`").
			WillReturnRows(sqlmock.NewRows([]string{
				"request_id", "operation", "request_status", "workspace_id", "datacenter_id", "created_at", "completed_at", "request_metadata",
			}).AddRow(
				"req-001", "vmDeploy", "New", "workspace-001", "dc-001", time.Now(), nil, `{"key":"value"}`,
			))

		mock.ExpectQuery("SELECT .* FROM `vm_deploy_instances`").
			WillReturnError(gorm.ErrRecordNotFound)

		repo := repo.NewVMRepository(mockDB, mockLogger)
		requests, instances, err := repo.GetAllVMRequestsWithInstances(ctx)

		assert.NoError(t, err)
		assert.Len(t, requests, 1)
		assert.Empty(t, instances)
	})


}
