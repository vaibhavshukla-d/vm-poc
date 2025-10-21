package dependency

import (
	"vm/internal/client"
	imagemanager "vm/internal/client/image_manager"
	inframonitor "vm/internal/client/infra_monitor"
	vmmonitor "vm/internal/client/vm_monitor"
	configmanager "vm/pkg/config-manager"
	"vm/pkg/constants"
)

// ClientDependency holds all the client-side dependencies.
type ClientDependency struct {
	ImageManagerClient *imagemanager.Client
	InfraMonitorClient *inframonitor.Client
	VmMonitorClient    *vmmonitor.Client
}

// SetupClientDependencies initializes all the clients.
func SetupClientDependencies(config *configmanager.Config) (*ClientDependency, error) {
	// Initialize the security sources for the API clients.
	imageManagerSecuritySource := &client.ImageManagerSecuritySource{}
	infraMonitorSecuritySource := &client.InfraMonitorSecuritySource{}
	vmMonitorSecuritySource := &client.VmMonitorSecuritySource{}

	// Url from config
	url := config.App.Application.Application

	// logger from config
	logger := config.Logger

	// Initialize the image-manager client.
	imageManagerClient, err := imagemanager.NewClient(url.ImageManagerServiceName, imageManagerSecuritySource)
	if err != nil {
		logger.Error("dependency", "setup", "Failed to create image-manager client", map[constants.ExtraKey]interface{}{constants.ErrorMessage: err})
		return nil, err
	}
	logger.Info("dependency", "setup", "Image-manager client initialized", nil)

	// Initialize the infra-monitor client.
	infraMonitorClient, err := inframonitor.NewClient(url.InfraMonitorServiceName, infraMonitorSecuritySource)
	if err != nil {
		logger.Error("dependency", "setup", "Failed to create infra-monitor client", map[constants.ExtraKey]interface{}{constants.ErrorMessage: err})
		return nil, err
	}
	logger.Info("dependency", "setup", "Infra-monitor client initialized", nil)

	// Initialize the vm-monitor client.
	vmMonitorClient, err := vmmonitor.NewClient(url.VmMonitorServiceName, vmMonitorSecuritySource)
	if err != nil {
		logger.Error("dependency", "setup", "Failed to create vm-monitor client", map[constants.ExtraKey]interface{}{constants.ErrorMessage: err})
		return nil, err
	}
	logger.Info("dependency", "setup", "Vm-monitor client initialized", nil)

	return &ClientDependency{
		ImageManagerClient: imageManagerClient,
		InfraMonitorClient: infraMonitorClient,
		VmMonitorClient:    vmMonitorClient,
	}, nil
}
