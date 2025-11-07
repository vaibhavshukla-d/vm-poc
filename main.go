//go:generate sh -c "mkdir -p internal/gen && go run -mod=mod github.com/ogen-go/ogen/cmd/ogen --clean --target internal/gen --config doc/ogen.yaml doc/openapi.yaml"
//go:generate sh -c "mkdir -p internal/client/image_manager && go run -mod=mod github.com/ogen-go/ogen/cmd/ogen --clean --target internal/client/image_manager --package client --config doc/client.ogen.yaml doc/image-manager.yaml"
//go:generate sh -c "mkdir -p internal/client/vm_monitor && go run -mod=mod github.com/ogen-go/ogen/cmd/ogen --clean --target internal/client/vm_monitor --package client --config doc/client.ogen.yaml doc/vm-monitor.yaml"
//go:generate sh -c "mkdir -p internal/client/infra_monitor && go run -mod=mod github.com/ogen-go/ogen/cmd/ogen --clean --target internal/client/infra_monitor --package client --config doc/client.ogen.yaml doc/infra-monitor.yaml"
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	promexporter "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"

	api "vm/internal/gen"
	"vm/internal/handler_impl"
	"vm/internal/repo"
	"vm/internal/service"
	"vm/pkg/constants"
	"vm/pkg/dependency"
	"vm/pkg/middleware"
)

func main() {
	// Create a root context
	ctx := context.Background()

	// Setup dependencies
	deps, err := dependency.Setup(ctx)
	if err != nil {
		panic(err)
	}

	// Set up Prometheus exporter BEFORE creating the server
	exporter, err := promexporter.New()
	if err != nil {
		deps.Logger.Fatal(constants.General, constants.Startup, "failed to create prometheus exporter", map[constants.ExtraKey]interface{}{"error": err})
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	otel.SetMeterProvider(provider)

	// Initialize repository and service
	vmRepo := repo.NewVMRepository(deps.Database, deps.Logger)
	vmService := service.NewVMService(vmRepo, deps.Logger)

	// Initialize handlers
	handler := handler_impl.NewHandler(vmService, deps)
	securityHandler := handler_impl.NewSecurityHandler(deps.Logger)

	// Create new server with OTel support
	server, err := api.NewServer(
		handler,
		securityHandler,
		api.WithTracerProvider(otel.GetTracerProvider()),
		api.WithMeterProvider(otel.GetMeterProvider()), // Add meter provider
	)
	if err != nil {
		deps.Logger.Fatal(constants.General, constants.Startup, "failed to create server", map[constants.ExtraKey]interface{}{"error": err})
	}

	// Start main application server
	addr := ":" + deps.Config.App.Application.Port
	wrappedHandler := middleware.RecoveryMiddleware(deps.Logger)(middleware.RequestIDMiddleware(server))
	httpServer := &http.Server{Addr: addr, Handler: wrappedHandler}

	// Start metrics server on a separate port
	metricsServer := &http.Server{
		Addr:    ":9090",
		Handler: promhttp.Handler(),
	}

	go func() {
		deps.Logger.Info(constants.General, constants.Startup, "Metrics server listening on :9090", nil)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			deps.Logger.Fatal(constants.General, constants.Startup, "metrics server error", map[constants.ExtraKey]interface{}{"error": err})
		}
	}()

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
	deps.Logger.Info(constants.General, constants.Startup, "Shutting down servers...", nil)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown both servers
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		deps.Logger.Error(constants.General, constants.Startup, "main server shutdown error", map[constants.ExtraKey]interface{}{"error": err})
	}

	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		deps.Logger.Error(constants.General, constants.Startup, "metrics server shutdown error", map[constants.ExtraKey]interface{}{"error": err})
	}

	// Shutdown meter provider
	if err := provider.Shutdown(shutdownCtx); err != nil {
		deps.Logger.Error(constants.General, constants.Startup, "meter provider shutdown error", map[constants.ExtraKey]interface{}{"error": err})
	}

	deps.Logger.Info(constants.General, constants.Startup, "Servers gracefully stopped", nil)
}