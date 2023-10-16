// Package app configures and runs application.
package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/almiluk/sipacks/config"
	"github.com/almiluk/sipacks/internal/adapter/repo"
	v1 "github.com/almiluk/sipacks/internal/controller/http/v1"
	"github.com/almiluk/sipacks/internal/usecase"
	"github.com/almiluk/sipacks/pkg/httpserver"
	"github.com/almiluk/sipacks/pkg/logger"
	"github.com/almiluk/sipacks/pkg/postgres"
	"github.com/labstack/echo/v4"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := repo.NewPGRepo(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use case
	translationUseCase, err := usecase.New(pg, cfg.FileStoragePath)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - usecase.New: %w", err))
	}

	// HTTP Server
	handler := echo.New()
	v1.NewRouter(handler, cfg.HTTP, l, translationUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	log.Printf("Config: %+v\n", cfg)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
