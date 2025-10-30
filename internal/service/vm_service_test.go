package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	api "vm/internal/gen"
	"vm/internal/modals"
	"vm/internal/service"

	mock_repo "vm/internal/repo/mock"
	"vm/pkg/constants"
	mock_logger "vm/pkg/logger/mock"
)

func TestCreateVMRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &mock_logger.StubLogger{}

	mockRepo := mock_repo.NewMockVMRepository(ctrl)
	vmSvc := service.NewVMService(mockRepo, logger)

	ctx := context.Background()

	// Sample deploy metadata
	deployReq := api.HCIDeployVM{
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

	metadataBytes, _ := json.Marshal(deployReq)
	metadata := string(metadataBytes)

	t.Run("Successful VM deploy", func(t *testing.T) {
		expectedVM := &modals.VMRequest{
			Operation:       string(constants.VMDeploy),
			RequestStatus:   string(constants.StatusNew),
			RequestMetadata: metadata,
		}

		mockRepo.EXPECT().
			CreateVMRequest(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, req *modals.VMRequest) error {
				req.RequestID = "req-123"
				return nil
			})

		vmName := deployReq.VmConfig.Name
		expectedInstances := []modals.VMDeployInstance{
			{RequestID: "req-123", VMName: fmt.Sprintf("%s_1", vmName), VMStatus: string(constants.VMINIT)},
			{RequestID: "req-123", VMName: fmt.Sprintf("%s_2", vmName), VMStatus: string(constants.VMINIT)},
		}

		mockRepo.EXPECT().
			CreateVMDeployInstances(ctx, expectedInstances).
			Return(nil)

		result, err := vmSvc.CreateVMRequest(ctx, constants.VMDeploy, constants.StatusNew, metadata)

		assert.NoError(t, err)
		assert.Equal(t, "req-123", result.RequestID)
		assert.Equal(t, expectedVM.Operation, result.Operation)
		assert.Equal(t, expectedVM.RequestStatus, result.RequestStatus)
		assert.Equal(t, expectedVM.RequestMetadata, result.RequestMetadata)
	})

	t.Run("Unmarshal failure", func(t *testing.T) {
		badMetadata := `{"invalid_json":}`

		mockRepo.EXPECT().
			CreateVMRequest(ctx, gomock.Any()).Return(nil)

		result, err := vmSvc.CreateVMRequest(ctx, constants.VMDeploy, constants.StatusNew, badMetadata)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("CreateVMRequest fails", func(t *testing.T) {
		mockRepo.EXPECT().
			CreateVMRequest(ctx, gomock.Any()).
			Return(errors.New("db error"))

		result, err := vmSvc.CreateVMRequest(ctx, constants.VMDeploy, constants.StatusNew, metadata)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("CreateVMDeployInstances fails", func(t *testing.T) {
		mockRepo.EXPECT().
			CreateVMRequest(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, req *modals.VMRequest) error {
				req.RequestID = "req-456"
				return nil
			})

		vmName := deployReq.VmConfig.Name
		expectedInstances := []modals.VMDeployInstance{
			{RequestID: "req-456", VMName: fmt.Sprintf("%s_1", vmName), VMStatus: string(constants.VMINIT)},
			{RequestID: "req-456", VMName: fmt.Sprintf("%s_2", vmName), VMStatus: string(constants.VMINIT)},
		}
		mockRepo.EXPECT().
			CreateVMDeployInstances(ctx, expectedInstances).
			Return(errors.New("deploy error"))

		result, err := vmSvc.CreateVMRequest(ctx, constants.VMDeploy, constants.StatusNew, metadata)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetVMDeployInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &mock_logger.StubLogger{}
	mockRepo := mock_repo.NewMockVMRepository(ctrl)
	vmSvc := service.NewVMService(mockRepo, logger)

	ctx := context.Background()
	requestID := "req-123"

	t.Run("Successful retrieval of VM deploy instances", func(t *testing.T) {
		expectedInstances := []*modals.VMDeployInstance{
			{RequestID: requestID, VMName: "vm-1", VMStatus: string(constants.VMINIT)},
			{RequestID: requestID, VMName: "vm-2", VMStatus: string(constants.VMINIT)},
		}

		mockRepo.EXPECT().
			GetVMDeployInstances(ctx, requestID).
			Return(expectedInstances, nil)

		result, err := vmSvc.GetVMDeployInstances(ctx, requestID)

		assert.NoError(t, err)
		assert.Equal(t, expectedInstances, result)
	})

	t.Run("Failure to retrieve VM deploy instances", func(t *testing.T) {
		mockRepo.EXPECT().
			GetVMDeployInstances(ctx, requestID).
			Return(nil, errors.New("db error"))

		result, err := vmSvc.GetVMDeployInstances(ctx, requestID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetVMRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &mock_logger.StubLogger{}
	mockRepo := mock_repo.NewMockVMRepository(ctrl)
	vmSvc := service.NewVMService(mockRepo, logger)

	ctx := context.Background()
	requestID := "req-123"

	t.Run("Successful retrieval of VM request", func(t *testing.T) {
		expectedVMRequest := &modals.VMRequest{
			RequestID:       requestID,
			Operation:       string(constants.VMDeploy),
			RequestStatus:   string(constants.StatusNew),
			WorkspaceId:     "workspace-001",
			DatacenterId:    "dc-001",
			CreatedAt:       time.Now(),
			CompletedAt:     nil,
			RequestMetadata: `{"key":"value"}`,
		}

		mockRepo.EXPECT().
			GetVMRequest(ctx, requestID).
			Return(expectedVMRequest, nil)

		result, err := vmSvc.GetVMRequest(ctx, requestID)

		assert.NoError(t, err)
		assert.Equal(t, expectedVMRequest, result)
	})

	t.Run("Failure to retrieve VM request", func(t *testing.T) {
		mockRepo.EXPECT().
			GetVMRequest(ctx, requestID).
			Return(nil, errors.New("db error"))

		result, err := vmSvc.GetVMRequest(ctx, requestID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetAllVMRequestsWithInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &mock_logger.StubLogger{}
	mockRepo := mock_repo.NewMockVMRepository(ctrl)
	vmSvc := service.NewVMService(mockRepo, logger)

	ctx := context.Background()

	t.Run("Successful retrieval of requests and instances", func(t *testing.T) {
		expectedRequests := []*modals.VMRequest{
			{
				RequestID:       "req-1",
				Operation:       string(constants.VMDeploy),
				RequestStatus:   string(constants.StatusNew),
				WorkspaceId:     "workspace-001",
				DatacenterId:    "dc-001",
				RequestMetadata: `{"key":"value"}`,
			},
			{
				RequestID:       "req-2",
				Operation:       string(constants.VMDelete),
				RequestStatus:   string(constants.StatusDone),
				WorkspaceId:     "workspace-002",
				DatacenterId:    "dc-002",
				RequestMetadata: `{"key":"value"}`,
			},
		}

		expectedInstances := []*modals.VMDeployInstance{
			{RequestID: "req-1", VMName: "vm-1", VMStatus: string(constants.VMINIT)},
			{RequestID: "req-2", VMName: "vm-2", VMStatus: string(constants.VMCLOSE)},
		}

		mockRepo.EXPECT().
			GetAllVMRequestsWithInstances(ctx).
			Return(expectedRequests, expectedInstances, nil)

		reqs, insts, reqCount, instCount, err := vmSvc.GetAllVMRequestsWithInstances(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedRequests, reqs)
		assert.Equal(t, expectedInstances, insts)
		assert.Equal(t, 2, reqCount)
		assert.Equal(t, 2, instCount)
	})

	t.Run("Failure to retrieve requests and instances", func(t *testing.T) {
		mockRepo.EXPECT().
			GetAllVMRequestsWithInstances(ctx).
			Return(nil, nil, errors.New("db error"))

		reqs, insts, reqCount, instCount, err := vmSvc.GetAllVMRequestsWithInstances(ctx)

		assert.Error(t, err)
		assert.Nil(t, reqs)
		assert.Nil(t, insts)
		assert.Equal(t, 0, reqCount)
		assert.Equal(t, 0, instCount)
	})
}
