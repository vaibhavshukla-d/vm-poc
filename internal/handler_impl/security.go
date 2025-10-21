package handler_impl

import (
	"context"
	"errors"
	api "vm/internal/gen"
)

type SecurityHandler struct {
	// Add your security configuration here
}

func NewSecurityHandler() *SecurityHandler {
	return &SecurityHandler{}
}
	
func (h *SecurityHandler) HandleBearer(ctx context.Context, operationName api.OperationName, t api.Bearer) (context.Context, error) {
	// For testing purposes:
	// 1. If no token is provided, return an error
	if t.Token == "" {
		return nil, errors.New("Missing Bearer token")
	}

	// 2. For testing, accept any non-empty token
	// In production, you would validate the token here
	return ctx, nil
}
