package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	dto "vm/internal/dtos"
	api "vm/internal/gen"
	logger "vm/pkg/cinterface"
	"vm/pkg/constants"

	"github.com/google/uuid"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RecoveryMiddleware(logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error(constants.General, constants.Api, "panic recovered", map[constants.ExtraKey]interface{}{
						"panic": rec,
					})

					err := dto.ApiResponseError{
						ErrorCode: constants.InternalServerErrorCode,
						Message:   "internal server error",
					}
					res := constants.MapServiceError(err, constants.InternalServerErrorCode, r.Context()).(api.ErrorResponse)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(res)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
