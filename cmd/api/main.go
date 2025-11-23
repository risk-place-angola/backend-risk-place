package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/middleware"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/router"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/bootstrap"
)

const (
	// HTTP server configuration constants
	maxHeaderBytes        = 1 << 20          // 1 MB maximum header size
	readHeaderTimeout     = 10 * time.Second // Timeout for reading request headers
	gracefulShutdownDelay = 30 * time.Second // Time to wait for graceful shutdown
)

func main() {
	os.Exit(run())
}

func run() int {
	container, err := bootstrap.NewContainer()
	if err != nil {
		log.Printf("failed to bootstrap app: %v", err)
		return 1
	}

	mux := router.SetupRoutes(container)

	handler := middleware.Use(
		mux,
		middleware.CORS,
		middleware.PanicRecovery,
	)

	// Production-grade HTTP server configuration for high concurrency
	srv := &http.Server{
		Addr:         "0.0.0.0:" + container.Cfg.Port,
		Handler:      handler,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
		// Limits for handling thousands of concurrent connections
		MaxHeaderBytes:    maxHeaderBytes,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Channel to listen for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Channel for server errors
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		slog.Info("Starting server",
			"port", container.Cfg.Port,
			"env", container.Cfg.AppEnv,
			"maxHeaderBytes", srv.MaxHeaderBytes,
			"readTimeout", srv.ReadTimeout,
			"writeTimeout", srv.WriteTimeout,
			"idleTimeout", srv.IdleTimeout,
		)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErrors:
		slog.Error("Server failed to start", "error", err)
		return 1
	case <-shutdown:
		slog.Info("Shutdown signal received, starting graceful shutdown...")
	}

	// Create context with timeout for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gracefulShutdownDelay)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		return 1
	}

	slog.Info("Server stopped gracefully")
	return 0
}
