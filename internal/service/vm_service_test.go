package service

import (
	"context"
	"errors"
	"testing"
	"vm/internal/modals"
	"vm/internal/repo/mocks"
	"vm/pkg/constants"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogger is a mock implementation of the cinterface.Logger for testing.
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Init() {
	m.Called()
}

func (m *MockLogger) Debug(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {
	m.Called(cat, sub, msg, extra)
}

func (m *MockLogger) Debugf(templateName string, args ...interface{}) {
	m.Called(templateName, args)
}

func (m *MockLogger) Info(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {
	m.Called(cat, sub, msg, extra)
}

func (m *MockLogger) Infof(templateName string, args ...interface{}) {
	m.Called(templateName, args)
}

func (m *MockLogger) Warn(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {
	m.Called(cat, sub, msg, extra)
}

func (m *MockLogger) Warnf(templateName string, args ...interface{}) {
	m.Called(templateName, args)
}

func (m *MockLogger) Error(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {
	m.Called(cat, sub, msg, extra)
}

func (m *MockLogger) Errorf(templateName string, args ...interface{}) {
	m.Called(templateName, args)
}

func (m *MockLogger) Fatal(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {
	m.Called(cat, sub, msg, extra)
}

func (m *MockLogger) Fatalf(templateName string, args ...interface{}) {
	m.Called(templateName, args)
}

func TestGetAllVMRequestsWithInstances(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mocks.VMRepository)
		mockLogger := new(MockLogger)
		vmService := NewVMService(mockRepo, mockLogger)

		expectedVMs := []*modals.VMRequest{{RequestID: "req1"}}
		expectedInstances := []*modals.VMDeployInstance{{RequestID: "req1", VMName: "vm1"}}

		mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		mockRepo.On("GetAllVMRequestsWithInstances", mock.Anything).Return(expectedVMs, expectedInstances, nil)

		vms, instances, reqCount, instCount, err := vmService.GetAllVMRequestsWithInstances(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedVMs, vms)
		assert.Equal(t, len(expectedVMs), reqCount)
		assert.Equal(t, len(expectedInstances), instCount)
		assert.Equal(t, expectedInstances, instances)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DB Error", func(t *testing.T) {
		mockRepo := new(mocks.VMRepository)
		mockLogger := new(MockLogger)
		vmService := NewVMService(mockRepo, mockLogger)

		dbError := errors.New("database error")
		mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		mockRepo.On("GetAllVMRequestsWithInstances", mock.Anything).Return(nil, nil, dbError)

		vms, instances, reqCount, instCount, err := vmService.GetAllVMRequestsWithInstances(context.Background())

		assert.Error(t, err)
		assert.Nil(t, vms)
		assert.Equal(t, 0, reqCount)
		assert.Equal(t, 0, instCount)
		assert.Nil(t, instances)
		assert.Equal(t, dbError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateVMRequest_Hybrid(t *testing.T) {
	metadata := `{
        "imageSource": {"imageId": "some-image-id"},
        "destination": {"hostId":"some-host-id","clusterId":"some-cluster-id"},
        "vmConfig":{"name":"test-vm","numberOfVms":1,"acceptEula":true},
        "storageConfig":{"defaultDatastoreId":"datastore-1"}
    }`

	// --- Unique test case kept separate ---
	t.Run("Invalid Metadata", func(t *testing.T) {
		mockRepo := new(mocks.VMRepository)
		mockLogger := new(MockLogger)
		vmService := NewVMService(mockRepo, mockLogger)

		mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

		invalidMetadata := `{"invalid-json": }`
		vmRequest, err := vmService.CreateVMRequest(context.Background(), constants.VMDeploy, constants.StatusPending, invalidMetadata)

		assert.Nil(t, vmRequest)
		assert.Error(t, err)
	})

	// --- Table-driven for cases with same flow ---
	tests := []struct {
		name           string
		setupMocks     func(mockRepo *mocks.VMRepository, mockLogger *MockLogger) string
		expectError    bool
		expectVMRequest bool
	}{
		{
			name: "Success",
			setupMocks: func(mockRepo *mocks.VMRepository, mockLogger *MockLogger) string {
				expectedID := uuid.New().String()
				mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

				mockRepo.On("CreateVMRequest", mock.Anything, mock.AnythingOfType("*modals.VMRequest")).
					Run(func(args mock.Arguments) {
						arg := args.Get(1).(*modals.VMRequest)
						arg.RequestID = expectedID
					}).Return(nil)

				mockRepo.On("CreateVMDeployInstances", mock.Anything, expectedID, "test-vm", 1).Return(nil)
				return expectedID
			},
			expectError: false,
			expectVMRequest: true,
		},
		{
			name: "DB Error on CreateVMRequest",
			setupMocks: func(mockRepo *mocks.VMRepository, mockLogger *MockLogger) string {
				mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

				mockRepo.On("CreateVMRequest", mock.Anything, mock.AnythingOfType("*modals.VMRequest")).
					Return(errors.New("DB insert failed"))
				return ""
			},
			expectError: true,
			expectVMRequest: false,
		},
		{
			name: "DB Error on CreateVMDeployInstances",
			setupMocks: func(mockRepo *mocks.VMRepository, mockLogger *MockLogger) string {
				expectedID := uuid.New().String()
				mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

				mockRepo.On("CreateVMRequest", mock.Anything, mock.AnythingOfType("*modals.VMRequest")).
					Run(func(args mock.Arguments) {
						arg := args.Get(1).(*modals.VMRequest)
						arg.RequestID = expectedID
					}).Return(nil)

				mockRepo.On("CreateVMDeployInstances", mock.Anything, expectedID, "test-vm", 1).
					Return(errors.New("DB insert failed"))

				return expectedID
			},
			expectError: true,
			expectVMRequest: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.VMRepository)
			mockLogger := new(MockLogger)
			vmService := NewVMService(mockRepo, mockLogger)

			tt.setupMocks(mockRepo, mockLogger)

			vmRequest, err := vmService.CreateVMRequest(context.Background(), constants.VMDeploy, constants.StatusPending, metadata)

			if tt.expectVMRequest {
				assert.NotNil(t, vmRequest)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, vmRequest)
				assert.Error(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}
