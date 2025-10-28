package security

import (
	"context"
	imagemanager "vm/internal/client/image_manager"
	inframonitor "vm/internal/client/infra_monitor"
	vmmonitor "vm/internal/client/vm_monitor"
)

// ImageManagerSecuritySource implements the security source for the image manager client.
type ImageManagerSecuritySource struct{}

// Bearer implements the Bearer method for the ImageManagerSecuritySource.
func (s *ImageManagerSecuritySource) Bearer(ctx context.Context, operationName imagemanager.OperationName) (imagemanager.Bearer, error) {
	return imagemanager.Bearer{Token: ""}, nil
}

// InfraMonitorSecuritySource implements the security source for the infra monitor client.
type InfraMonitorSecuritySource struct{}

// Bearer implements the Bearer method for the InfraMonitorSecuritySource.
func (s *InfraMonitorSecuritySource) Bearer(ctx context.Context, operationName inframonitor.OperationName) (inframonitor.Bearer, error) {
	return inframonitor.Bearer{Token: ""}, nil
}

// VmMonitorSecuritySource implements the security source for the vm monitor client.
type VmMonitorSecuritySource struct{}

// Bearer implements the Bearer method for the VmMonitorSecuritySource.
func (s *VmMonitorSecuritySource) Bearer(ctx context.Context, operationName vmmonitor.OperationName) (vmmonitor.Bearer, error) {
	return vmmonitor.Bearer{Token: ""}, nil
}
