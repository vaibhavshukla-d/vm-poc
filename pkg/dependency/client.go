package dependency

import (
	"vm/internal/client"
	api "vm/internal/client_gen"
	"vm/pkg/cinterface"
	configmanager "vm/pkg/config-manager"
	"vm/pkg/constants"
)

// ClientDependency holds all the client-side dependencies.
type ClientDependency struct {
	ImageManagerClient *api.Client
}

// SetupClientDependencies initializes all the clients.
func SetupClientDependencies(config *configmanager.Config, logger cinterface.Logger) (*ClientDependency, error) {
	// Initialize the security source for the API client.
	apiSecuritySource := &client.APISecuritySource{}

	// Initialize the ogen client.
	ogenClient, err := api.NewClient(config.App.Application.Application.ImageManagerServiceName, apiSecuritySource)
	if err != nil {
		logger.Error("dependency", "setup", "Failed to create ogen client", map[constants.ExtraKey]interface{}{constants.ErrorMessage: err})
		return nil, err
	}
	logger.Info("dependency", "setup", "Ogen client initialized", nil)

	return &ClientDependency{
		ImageManagerClient: ogenClient,
	}, nil
}
