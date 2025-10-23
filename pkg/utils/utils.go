package utils

import (
	"context"
	"errors"
)

type contextKey string

const WorkspaceIDKey contextKey = "workspace_id"

func GetWorkspaceIDFromContext(ctx context.Context) (string, error) {
	workspaceIDValue := ctx.Value(WorkspaceIDKey)
	workspaceID, ok := workspaceIDValue.(string)
	if !ok || workspaceID == "" {
		return "", errors.New("workspace_id not found in context")
	}
	return workspaceID, nil
}
