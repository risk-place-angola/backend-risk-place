package main

import (
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/middleware"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/router"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/bootstrap"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	container, err := bootstrap.NewContainer()
	if err != nil {
		log.Fatalf("failed to bootstrap app: %v", err)
	}

	mux := router.SetupRoutes(container)

	handler := middleware.Use(
		mux,
		middleware.CORS,
		middleware.PanicRecovery,
	)

	slog.Info("Starting server on port " + container.Cfg.Port)

	srv := &http.Server{
		Addr:         "0.0.0.0:" + container.Cfg.Port,
		Handler:      handler,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	log.Fatal(srv.ListenAndServe())
}
