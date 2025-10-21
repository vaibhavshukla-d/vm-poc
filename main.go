//go:generate sh -c "mkdir -p internal/gen && go run github.com/ogen-go/ogen/cmd/ogen --clean --target internal/gen --config doc/ogen.yaml doc/openapi.yaml"
//go:generate sh -c "mkdir -p internal/client_gen && go run github.com/ogen-go/ogen/cmd/ogen --clean --target internal/client_gen --config doc/client.ogen.yaml doc/client.yaml"
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "vm/internal/gen"
	"vm/internal/handler_impl"
	"vm/internal/repo"
	"vm/internal/service"
	"vm/pkg/constants"
	"vm/pkg/dependency"
)

func main() {
	// Create a root context
	ctx := context.Background()

	// Setup dependencies
	deps, err := dependency.Setup(ctx)
	if err != nil {
		panic(err)
	}

	// Initialize repository and service
	vmRepo := repo.NewVMRepository(deps.Database, deps.Logger)
	vmService := service.NewVMService(vmRepo, deps.Logger)

	// Initialize handlers
	handler := handler_impl.NewHandler(vmService, deps)
	securityHandler := handler_impl.NewSecurityHandler()

	// Create new server
	server, err := api.NewServer(handler, securityHandler)
	if err != nil {
		deps.Logger.Fatal(constants.General, constants.Startup, "failed to create server", map[constants.ExtraKey]interface{}{"error": err})
	}

	// Start server
	addr := ":" + deps.Config.App.Application.Application.Port
	httpServer := &http.Server{Addr: addr, Handler: server}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		deps.Logger.Info(constants.General, constants.Startup, "Server listening on "+addr, nil)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			deps.Logger.Fatal(constants.General, constants.Startup, "server error", map[constants.ExtraKey]interface{}{"error": err})
		}
	}()

	<-stop
	deps.Logger.Info(constants.General, constants.Startup, "Shutting down server...", nil)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		deps.Logger.Fatal(constants.General, constants.Startup, "server shutdown error", map[constants.ExtraKey]interface{}{"error": err})
	}

	deps.Logger.Info(constants.General, constants.Startup, "Server gracefully stopped", nil)
}
