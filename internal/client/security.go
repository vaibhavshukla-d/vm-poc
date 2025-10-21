package client

import (
	"context"

	api "vm/internal/client_gen"
)

// APISecuritySource implements the api.SecuritySource interface.
type APISecuritySource struct{}

// Bearer returns a bearer token.
func (s *APISecuritySource) Bearer(ctx context.Context, operationName api.OperationName) (api.Bearer, error) {
	// In a real application, you would fetch a valid token from a secure source.
	return api.Bearer{Token: "dummy-token"}, nil
}
