package handler_impl_test

import (
	"context"
	"errors"
	"testing"
	"time"
	api "vm/internal/gen"
	"vm/internal/handler_impl"
	"vm/internal/modals"
	configmanager "vm/pkg/config-manager"
	"vm/pkg/constants"
	"vm/pkg/dependency"

	mock_logger "vm/pkg/logger/mock"

	mock_db "vm/pkg/db/mock"

	mock_service "vm/internal/service/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestEditVM_ValidateVMExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVMService := mock_service.NewMockVMService(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	mockDB := mock_db.NewMockDatabase(ctrl)

	baseConfig := &configmanager.Config{
		App: configmanager.ApplicationConfig{
			Application: configmanager.Application{
				ValidateClientRequest: true, // toggle this in different tests
			},
		},
	}

	deps := &dependency.Dependency{
		Ctx:              context.Background(),
		Logger:           mockLogger,
		Database:         mockDB,
		Config:           baseConfig,
		ClientDependency: &dependency.ClientDependency{},
	}

	handler := handler_impl.NewHandler(mockVMService, deps)

	params := api.EditVMParams{VMID: "vm-123"}
	req := &api.EditVM{}

	t.Run("Skip validation when ValidateClientRequest is false", func(t *testing.T) {
		deps.Config.App.Application.ValidateClientRequest = false
		mockVMService.EXPECT().CreateVMRequest(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&modals.VMRequest{RequestID: "req-001"}, nil)

		res, err := handler.EditVM(context.Background(), req, params)
		assert.NoError(t, err)
		assert.IsType(t, &api.EmptyResponseHeaders{}, res)
	})

	t.Run("Failure - CreateVMRequest error", func(t *testing.T) {
		deps.Config.App.Application.ValidateClientRequest = false

		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, errors.New("create failed"))

		res, err := handler.EditVM(context.Background(), req, params)
		assert.NoError(t, err)
		assert.IsType(t, &api.EditVMInternalServerError{}, res)
		assert.Equal(t, "Failed to create VM request", res.(*api.EditVMInternalServerError).Message)
	})

}

func TestHandler_HCIDeployVM(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVMService := mock_service.NewMockVMService(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	mockDB := mock_db.NewMockDatabase(ctrl)

	baseConfig := &configmanager.Config{
		App: configmanager.ApplicationConfig{
			Application: configmanager.Application{
				ValidateClientRequest: false, // skip host validation
			},
		},
	}

	deps := &dependency.Dependency{
		Ctx:              context.Background(),
		Logger:           mockLogger,
		Database:         mockDB,
		Config:           baseConfig,
		ClientDependency: &dependency.ClientDependency{},
	}

	handler := handler_impl.NewHandler(mockVMService, deps)

	req := &api.HCIDeployVM{
		Destination: api.OptHCIDeployVMDestination{
			Set: true,
			Value: api.HCIDeployVMDestination{
				ClusterId:      api.OptString{Value: "cluster-uuid-123", Set: true},
				FolderId:       api.OptString{Value: "folder-uuid-456", Set: true},
				HostId:         api.OptString{Value: "host-uuid-789", Set: true},
				ResourcePoolId: api.OptString{Value: "rp-uuid-101", Set: true},
			},
		},
		ImageSource: api.OptHCIDeployVMImageSource{
			Set: true,
			Value: api.HCIDeployVMImageSource{
				ImageId:         api.OptString{Value: "image-uuid-212", Set: true},
				ImageName:       api.OptString{Value: "ubuntu-22.04-template", Set: true},
				ImageSourceType: api.OptHCIDeployVMImageSourceImageSourceType{Value: "HYPERVISOR_IMAGE_LIBRARY", Set: true},
			},
		},
		NetworkConfig: api.OptHCIDeployVMNetworkConfig{
			Set: true,
			Value: api.HCIDeployVMNetworkConfig{
				IpAllocationPolicy: api.OptHCIDeployVMNetworkConfigIpAllocationPolicy{
					Set:   true,
					Value: api.HCIDeployVMNetworkConfigIpAllocationPolicy("DHCP_POLICY"),
				},
				NetworkMapping: []api.HCIDeployVMNetworkConfigNetworkMappingItem{
					{
						Name:    api.OptString{Value: "VM Network", Set: true},
						Network: api.OptString{Value: "network-uuid-313", Set: true},
					},
				},
			},
		},
		StorageConfig: api.HCIDeployVMStorageConfig{
			DefaultDatastoreId: "datastore-uuid-414",
			ProvisioningType: api.OptHCIDeployVMStorageConfigProvisioningType{
				Set:   true,
				Value: api.HCIDeployVMStorageConfigProvisioningType("THIN"),
			},
		},
		VmConfig: api.HCIDeployVMVmConfig{
			AcceptEula:  true,
			Annotation:  api.OptString{Value: "This is a sample VM deployed via API.", Set: true},
			Locale:      api.OptString{Value: "en-US", Set: true},
			Name:        "my-full-config-vm",
			NumberOfVms: api.OptInt{Value: 2, Set: true},
			PowerOn:     api.OptBool{Value: true, Set: true},
			PropertyConfig: []api.HCIDeployVMVmConfigPropertyConfigItem{
				{
					Key:   api.OptString{Value: "guestinfo.hostname", Set: true},
					Value: api.OptString{Value: "my-vm", Set: true},
				},
			},
		},
		VmPolicy: []api.HCIDeployVMVmPolicyItem{
			{ID: api.OptString{Value: "policy-uuid-515", Set: true}, Type: api.OptHCIDeployVMVmPolicyItemType{Value: "VM_PROVISIONING_POLICY", Set: true}},
		},
	}

	t.Run("Success - request created", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMDeploy, constants.StatusNew, gomock.Any()).
			Return(&modals.VMRequest{RequestID: "req-001"}, nil)

		res, err := handler.HCIDeployVM(context.Background(), req)
		assert.NoError(t, err)
		assert.IsType(t, &api.EmptyResponseHeaders{}, res)
		assert.Equal(t, "/virtualization/v1beta1/virtual-machines-request/req-001", res.(*api.EmptyResponseHeaders).Location.Value)
	})

	t.Run("Failure - CreateVMRequest error", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMDeploy, constants.StatusNew, gomock.Any()).
			Return(nil, errors.New("create failed"))

		res, err := handler.HCIDeployVM(context.Background(), req)
		assert.NoError(t, err)
		assert.IsType(t, &api.HCIDeployVMInternalServerError{}, res)
		assert.Equal(t, "Failed to create VM request", res.(*api.HCIDeployVMInternalServerError).Message)
	})
}

func TestHandler_VMDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVMService := mock_service.NewMockVMService(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	mockDB := mock_db.NewMockDatabase(ctrl)

	baseConfig := &configmanager.Config{
		App: configmanager.ApplicationConfig{
			Application: configmanager.Application{
				ValidateClientRequest: false, // skip validateVMExists
			},
		},
	}

	deps := &dependency.Dependency{
		Ctx:              context.Background(),
		Logger:           mockLogger,
		Database:         mockDB,
		Config:           baseConfig,
		ClientDependency: &dependency.ClientDependency{},
	}

	handler := handler_impl.NewHandler(mockVMService, deps)

	params := api.VMDeleteParams{
		VMID: "vm-uuid-999",
	}

	t.Run("Success - VMDelete request created", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMDelete, constants.StatusNew, gomock.Any()).
			Return(&modals.VMRequest{RequestID: "req-001"}, nil)

		res, err := handler.VMDelete(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.EmptyResponseHeaders{}, res)
		assert.Equal(t, "/virtualization/v1beta1/virtual-machines-request/req-001", res.(*api.EmptyResponseHeaders).Location.Value)
	})

	t.Run("Failure - CreateVMRequest error", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMDelete, constants.StatusNew, gomock.Any()).
			Return(nil, errors.New("create failed"))

		res, err := handler.VMDelete(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.VMDeleteInternalServerError{}, res)
		assert.Equal(t, "Failed to create VM request", res.(*api.VMDeleteInternalServerError).Message)
	})
}

func TestHandler_VMPowerOff(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVMService := mock_service.NewMockVMService(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	mockDB := mock_db.NewMockDatabase(ctrl)

	baseConfig := &configmanager.Config{
		App: configmanager.ApplicationConfig{
			Application: configmanager.Application{
				ValidateClientRequest: false, // skip validateVMExists
			},
		},
	}

	deps := &dependency.Dependency{
		Ctx:              context.Background(),
		Logger:           mockLogger,
		Database:         mockDB,
		Config:           baseConfig,
		ClientDependency: &dependency.ClientDependency{},
	}

	handler := handler_impl.NewHandler(mockVMService, deps)

	params := api.VMPowerOffParams{
		VMID: "vm-uuid-888",
	}

	t.Run("Success - VMPowerOff request created", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMPowerOff, constants.StatusNew, gomock.Any()).
			Return(&modals.VMRequest{RequestID: "req-001"}, nil)

		res, err := handler.VMPowerOff(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.EmptyResponseHeaders{}, res)
		assert.Equal(t, "/virtualization/v1beta1/virtual-machines-request/req-001", res.(*api.EmptyResponseHeaders).Location.Value)
	})

	t.Run("Failure - CreateVMRequest error", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMPowerOff, constants.StatusNew, gomock.Any()).
			Return(nil, errors.New("create failed"))

		res, err := handler.VMPowerOff(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.VMPowerOffInternalServerError{}, res)
		assert.Equal(t, "Failed to create VM request", res.(*api.VMPowerOffInternalServerError).Message)
	})
}

func TestHandler_VMPowerOn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVMService := mock_service.NewMockVMService(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	mockDB := mock_db.NewMockDatabase(ctrl)

	baseConfig := &configmanager.Config{
		App: configmanager.ApplicationConfig{
			Application: configmanager.Application{
				ValidateClientRequest: false, // skip validateVMExists
			},
		},
	}

	deps := &dependency.Dependency{
		Ctx:              context.Background(),
		Logger:           mockLogger,
		Database:         mockDB,
		Config:           baseConfig,
		ClientDependency: &dependency.ClientDependency{},
	}

	handler := handler_impl.NewHandler(mockVMService, deps)

	params := api.VMPowerOnParams{
		VMID: "vm-uuid-777",
	}

	t.Run("Success - VMPowerOn request created", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMPowerOn, constants.StatusNew, gomock.Any()).
			Return(&modals.VMRequest{RequestID: "req-001"}, nil)

		res, err := handler.VMPowerOn(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.EmptyResponseHeaders{}, res)
		assert.Equal(t, "/virtualization/v1beta1/virtual-machines-request/req-001", res.(*api.EmptyResponseHeaders).Location.Value)
	})

	t.Run("Failure - CreateVMRequest error", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMPowerOn, constants.StatusNew, gomock.Any()).
			Return(nil, errors.New("create failed"))

		res, err := handler.VMPowerOn(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.VMPowerOnInternalServerError{}, res)
		assert.Equal(t, "Failed to create VM request", res.(*api.VMPowerOnInternalServerError).Message)
	})
}

func TestHandler_VMPowerReset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVMService := mock_service.NewMockVMService(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	mockDB := mock_db.NewMockDatabase(ctrl)

	baseConfig := &configmanager.Config{
		App: configmanager.ApplicationConfig{
			Application: configmanager.Application{
				ValidateClientRequest: false, // skip validateVMExists
			},
		},
	}

	deps := &dependency.Dependency{
		Ctx:              context.Background(),
		Logger:           mockLogger,
		Database:         mockDB,
		Config:           baseConfig,
		ClientDependency: &dependency.ClientDependency{},
	}

	handler := handler_impl.NewHandler(mockVMService, deps)

	params := api.VMPowerResetParams{
		VMID: "vm-uuid-666",
	}

	t.Run("Success - VMPowerReset request created", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMReset, constants.StatusNew, gomock.Any()).
			Return(&modals.VMRequest{RequestID: "req-001"}, nil)

		res, err := handler.VMPowerReset(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.EmptyResponseHeaders{}, res)
		assert.Equal(t, "/virtualization/v1beta1/virtual-machines-request/req-001", res.(*api.EmptyResponseHeaders).Location.Value)
	})

	t.Run("Failure - CreateVMRequest error", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMReset, constants.StatusNew, gomock.Any()).
			Return(nil, errors.New("create failed"))

		res, err := handler.VMPowerReset(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.VMPowerResetInternalServerError{}, res)
		assert.Equal(t, "Failed to create VM request", res.(*api.VMPowerResetInternalServerError).Message)
	})
}

func TestHandler_VMRefresh(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVMService := mock_service.NewMockVMService(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	mockDB := mock_db.NewMockDatabase(ctrl)

	baseConfig := &configmanager.Config{
		App: configmanager.ApplicationConfig{
			Application: configmanager.Application{
				ValidateClientRequest: false, // skip validateVMExists
			},
		},
	}

	deps := &dependency.Dependency{
		Ctx:              context.Background(),
		Logger:           mockLogger,
		Database:         mockDB,
		Config:           baseConfig,
		ClientDependency: &dependency.ClientDependency{},
	}

	handler := handler_impl.NewHandler(mockVMService, deps)

	params := api.VMRefreshParams{
		VMID: "vm-uuid-555",
	}

	t.Run("Success - VMRefresh request created", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMRefresh, constants.StatusNew, gomock.Any()).
			Return(&modals.VMRequest{RequestID: "req-001"}, nil)

		res, err := handler.VMRefresh(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.EmptyResponseHeaders{}, res)
		assert.Equal(t, "/virtualization/v1beta1/virtual-machines-request/req-001", res.(*api.EmptyResponseHeaders).Location.Value)
	})

	t.Run("Failure - CreateVMRequest error", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMRefresh, constants.StatusNew, gomock.Any()).
			Return(nil, errors.New("create failed"))

		res, err := handler.VMRefresh(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.VMRefreshInternalServerError{}, res)
		assert.Equal(t, "Failed to create VM request", res.(*api.VMRefreshInternalServerError).Message)
	})
}

func TestHandler_VMRestartGuestOS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVMService := mock_service.NewMockVMService(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	mockDB := mock_db.NewMockDatabase(ctrl)

	baseConfig := &configmanager.Config{
		App: configmanager.ApplicationConfig{
			Application: configmanager.Application{
				ValidateClientRequest: false, // skip validateVMExists
			},
		},
	}

	deps := &dependency.Dependency{
		Ctx:              context.Background(),
		Logger:           mockLogger,
		Database:         mockDB,
		Config:           baseConfig,
		ClientDependency: &dependency.ClientDependency{},
	}

	handler := handler_impl.NewHandler(mockVMService, deps)

	params := api.VMRestartGuestOSParams{
		VMID: "vm-uuid-444",
	}

	t.Run("Success - VMRestartGuestOS request created", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMRestartGuestOS, constants.StatusNew, gomock.Any()).
			Return(&modals.VMRequest{RequestID: "req-001"}, nil)

		res, err := handler.VMRestartGuestOS(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.EmptyResponseHeaders{}, res)
		assert.Equal(t, "/virtualization/v1beta1/virtual-machines-request/req-001", res.(*api.EmptyResponseHeaders).Location.Value)
	})

	t.Run("Failure - CreateVMRequest error", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMRestartGuestOS, constants.StatusNew, gomock.Any()).
			Return(nil, errors.New("create failed"))

		res, err := handler.VMRestartGuestOS(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.VMRestartGuestOSInternalServerError{}, res)
		assert.Equal(t, "Failed to create VM request", res.(*api.VMRestartGuestOSInternalServerError).Message)
	})
}

func TestHandler_VMShutdownGuestOS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVMService := mock_service.NewMockVMService(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	mockDB := mock_db.NewMockDatabase(ctrl)

	baseConfig := &configmanager.Config{
		App: configmanager.ApplicationConfig{
			Application: configmanager.Application{
				ValidateClientRequest: false, // skip validateVMExists
			},
		},
	}

	deps := &dependency.Dependency{
		Ctx:              context.Background(),
		Logger:           mockLogger,
		Database:         mockDB,
		Config:           baseConfig,
		ClientDependency: &dependency.ClientDependency{},
	}

	handler := handler_impl.NewHandler(mockVMService, deps)

	params := api.VMShutdownGuestOSParams{
		VMID: "vm-uuid-333",
	}

	t.Run("Success - VMShutdownGuestOS request created", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMShutdownGuestOS, constants.StatusNew, gomock.Any()).
			Return(&modals.VMRequest{RequestID: "req-001"}, nil)

		res, err := handler.VMShutdownGuestOS(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.EmptyResponseHeaders{}, res)
		assert.Equal(t, "/virtualization/v1beta1/virtual-machines-request/req-001", res.(*api.EmptyResponseHeaders).Location.Value)
	})

	t.Run("Failure - CreateVMRequest error", func(t *testing.T) {
		mockVMService.EXPECT().
			CreateVMRequest(gomock.Any(), constants.VMShutdownGuestOS, constants.StatusNew, gomock.Any()).
			Return(nil, errors.New("create failed"))

		res, err := handler.VMShutdownGuestOS(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.VMShutdownGuestOSInternalServerError{}, res)
		assert.Equal(t, "Failed to create VM request", res.(*api.VMShutdownGuestOSInternalServerError).Message)
	})
}

func TestHandler_GetVirtualMachineRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVMService := mock_service.NewMockVMService(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	mockDB := mock_db.NewMockDatabase(ctrl)

	baseConfig := &configmanager.Config{
		App: configmanager.ApplicationConfig{
			Application: configmanager.Application{
				ValidateClientRequest: false,
			},
		},
	}

	deps := &dependency.Dependency{
		Ctx:              context.Background(),
		Logger:           mockLogger,
		Database:         mockDB,
		Config:           baseConfig,
		ClientDependency: &dependency.ClientDependency{},
	}

	handler := handler_impl.NewHandler(mockVMService, deps)

	requestID := "req-123"
	params := api.GetVirtualMachineRequestParams{
		RequestID: requestID,
	}

	t.Run("Success - VM request and deploy instances found", func(t *testing.T) {
		mockVMService.EXPECT().
			GetVMRequest(gomock.Any(), requestID).
			Return(&modals.VMRequest{
				RequestID:       requestID,
				Operation:       "DEPLOY",
				RequestStatus:   "NEW",
				WorkspaceId:     "ws-001",
				DatacenterId:    "dc-001",
				CreatedAt:       time.Now(),
				RequestMetadata: `{"key":"value"}`,
			}, nil)

		mockVMService.EXPECT().
			GetVMDeployInstances(gomock.Any(), requestID).
			Return([]*modals.VMDeployInstance{
				{
					RequestID:      requestID,
					VMID:           "vm-001",
					VMName:         "test-vm",
					VMStatus:       "DEPLOYED",
					VMStateMessage: "Running",
				},
			}, nil)

		res, err := handler.GetVirtualMachineRequest(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.VMRequestWithDeploy{}, res)
		assert.Equal(t, requestID, res.(*api.VMRequestWithDeploy).VMRequest.RequestId)
		assert.Len(t, res.(*api.VMRequestWithDeploy).VMDeployList, 1)
	})

	t.Run("Failure - VM request not found", func(t *testing.T) {
		mockVMService.EXPECT().
			GetVMRequest(gomock.Any(), requestID).
			Return(nil, gorm.ErrRecordNotFound)

		res, err := handler.GetVirtualMachineRequest(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.GetVirtualMachineRequestNotFound{}, res)
		assert.Equal(t, "VM request not found", res.(*api.GetVirtualMachineRequestNotFound).Message)
	})

	t.Run("Failure - unexpected error from GetVMRequest", func(t *testing.T) {
		mockVMService.EXPECT().
			GetVMRequest(gomock.Any(), requestID).
			Return(nil, errors.New("db failure"))

		res, err := handler.GetVirtualMachineRequest(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.GetVirtualMachineRequestInternalServerError{}, res)
		assert.Equal(t, "db failure", res.(*api.GetVirtualMachineRequestInternalServerError).Message)
	})

	t.Run("Success - deploy instances not found (optional)", func(t *testing.T) {
		mockVMService.EXPECT().
			GetVMRequest(gomock.Any(), requestID).
			Return(&modals.VMRequest{
				RequestID:       requestID,
				Operation:       "DEPLOY",
				RequestStatus:   "NEW",
				WorkspaceId:     "ws-001",
				DatacenterId:    "dc-001",
				CreatedAt:       time.Now(),
				RequestMetadata: `{"key":"value"}`}, nil)

		mockVMService.EXPECT().
			GetVMDeployInstances(gomock.Any(), requestID).
			Return(nil, gorm.ErrRecordNotFound)

		res, err := handler.GetVirtualMachineRequest(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.VMRequestWithDeploy{}, res)
		assert.Len(t, res.(*api.VMRequestWithDeploy).VMDeployList, 0)
	})

	t.Run("Failure - unexpected error from GetVMDeployInstances", func(t *testing.T) {
		mockVMService.EXPECT().
			GetVMRequest(gomock.Any(), requestID).
			Return(&modals.VMRequest{
				RequestID:       requestID,
				Operation:       "DEPLOY",
				RequestStatus:   "NEW",
				WorkspaceId:     "ws-001",
				DatacenterId:    "dc-001",
				CreatedAt:       time.Now(),
				RequestMetadata: `{"key":"value"}`,
			}, nil)

		mockVMService.EXPECT().
			GetVMDeployInstances(gomock.Any(), requestID).
			Return(nil, errors.New("deploy lookup failed"))

		res, err := handler.GetVirtualMachineRequest(context.Background(), params)
		assert.NoError(t, err)
		assert.IsType(t, &api.GetVirtualMachineRequestInternalServerError{}, res)
		assert.Equal(t, "deploy lookup failed", res.(*api.GetVirtualMachineRequestInternalServerError).Message)
	})
}

func TestHandler_GetVirtualMachineRequestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVMService := mock_service.NewMockVMService(ctrl)
	mockLogger := &mock_logger.StubLogger{}
	mockDB := mock_db.NewMockDatabase(ctrl)

	baseConfig := &configmanager.Config{
		App: configmanager.ApplicationConfig{
			Application: configmanager.Application{
				ValidateClientRequest: false,
			},
		},
	}

	deps := &dependency.Dependency{
		Ctx:              context.Background(),
		Logger:           mockLogger,
		Database:         mockDB,
		Config:           baseConfig,
		ClientDependency: &dependency.ClientDependency{},
	}

	handler := handler_impl.NewHandler(mockVMService, deps)

	t.Run("Success - returns VM request list", func(t *testing.T) {
		mockVMService.EXPECT().
			GetAllVMRequestsWithInstances(gomock.Any()).
			Return([]*modals.VMRequest{
				{
					RequestID:      "req-001",
					Operation:      "DEPLOY",
					RequestStatus:  "NEW",
					WorkspaceId:    "ws-001",
					DatacenterId:   "dc-001",
					CreatedAt:      time.Now(),
					RequestMetadata: `{"key":"value"}`,
				},
			}, []*modals.VMDeployInstance{
				{
					RequestID:      "req-001",
					VMID:           "vm-001",
					VMName:         "test-vm",
					VMStatus:       "DEPLOYED",
					VMStateMessage: "Running",
				},
			}, 1, 1, nil)

		res, err := handler.GetVirtualMachineRequestList(context.Background())
		assert.NoError(t, err)
		assert.IsType(t, &api.VMRequestsList{}, res)

		list := res.(*api.VMRequestsList)
		assert.Equal(t, 1, list.VMMinusRequestsListCount.Value)
		assert.Equal(t, 1, list.VMMinusDeployListCount.Value)
		assert.Len(t, list.Items.Value.VMRequetsList, 1)
		assert.Len(t, list.Items.Value.VMDeployList, 1)
	})

	t.Run("Failure - service returns error", func(t *testing.T) {
		mockVMService.EXPECT().
			GetAllVMRequestsWithInstances(gomock.Any()).
			Return(nil, nil, 0, 0, errors.New("db failure"))

		res, err := handler.GetVirtualMachineRequestList(context.Background())
		assert.NoError(t, err)
		assert.IsType(t, &api.GetVirtualMachineRequestListInternalServerError{}, res)
		assert.Equal(t, "db failure", res.(*api.GetVirtualMachineRequestListInternalServerError).Message)
	})
}