package handler_impl

import (
	"context"
	"errors"
	api "vm/internal/gen"
	"vm/pkg/cinterface"
	"vm/pkg/constants"
	"vm/pkg/utils"

	"github.com/golang-jwt/jwt/v4"
)

type SecurityHandler struct {
	logger cinterface.Logger
}

func NewSecurityHandler(logger cinterface.Logger) *SecurityHandler {
	return &SecurityHandler{
		logger: logger,
	}
}

func (h *SecurityHandler) HandleBearer(ctx context.Context, operationName api.OperationName, t api.Bearer) (context.Context, error) {
	if t.Token == "" {
		h.logger.Error(constants.General, constants.Api, "Missing Bearer token", nil)
		return nil, errors.New("Missing Bearer token")
	}
	ctx = context.WithValue(ctx, constants.BearerTokenKey, t.Token)

	// Parse token without verifying signature (for testing only)
	token, _, err := new(jwt.Parser).ParseUnverified(t.Token, jwt.MapClaims{})
	if err != nil {
		h.logger.Error(constants.General, constants.Api, "Failed to parse token", map[constants.ExtraKey]interface{}{
			"error": err,
		})
		return ctx, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		h.logger.Error(constants.General, constants.Api, "Invalid token claims", nil)
		return ctx, nil
	}

	// Extract "id" from claims
	if id, ok := claims["id"].(string); ok && id != "" {
		ctx = context.WithValue(ctx, utils.WorkspaceIDKey, id)
	} else {
		h.logger.Info(constants.General, constants.Api, "Token does not contain 'id' claim", nil)
	}

	return ctx, nil
}
