package client

import (
	"context"
	imagemanager "vm/internal/client/image_manager"
	inframonitor "vm/internal/client/infra_monitor"
	vmmonitor "vm/internal/client/vm_monitor"
	"vm/pkg/constants"
)

// ImageManagerSecuritySource implements security for the image-manager client.
type ImageManagerSecuritySource struct{}

// Bearer returns a bearer token for image-manager.
func (s *ImageManagerSecuritySource) Bearer(ctx context.Context, operationName imagemanager.OperationName) (imagemanager.Bearer, error) {
	token, _ := ctx.Value(constants.BearerTokenKey).(string)
	return imagemanager.Bearer{Token: token}, nil
}

// InfraMonitorSecuritySource implements security for the infra-monitor client.
type InfraMonitorSecuritySource struct{}

// Bearer returns a bearer token for infra-monitor.
func (s *InfraMonitorSecuritySource) Bearer(ctx context.Context, operationName inframonitor.OperationName) (inframonitor.Bearer, error) {
	token, _ := ctx.Value(constants.BearerTokenKey).(string)
	return inframonitor.Bearer{Token: token}, nil
}

// VmMonitorSecuritySource implements security for the vm-monitor client.
type VmMonitorSecuritySource struct{}

// Bearer returns a bearer token for vm-monitor.
func (s *VmMonitorSecuritySource) Bearer(ctx context.Context, operationName vmmonitor.OperationName) (vmmonitor.Bearer, error) {
	token, _ := ctx.Value(constants.BearerTokenKey).(string)
	return vmmonitor.Bearer{Token: token}, nil
}
