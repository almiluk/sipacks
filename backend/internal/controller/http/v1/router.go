// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/almiluk/sipacks/config"
	"github.com/almiluk/sipacks/internal/controller"

	// Swagger docs.
	_ "github.com/almiluk/sipacks/docs"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/almiluk/sipacks/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       	API for managing question packs for the 'SIGame' game
// @version     	0.0.0
// @contact.name   	almiluk
// @contact.email  	almiluk@gmail.com
// @host        	localhost:8080
// @BasePath    	/
func NewRouter(handler *echo.Echo, cfg config.HTTP, l logger.Interface, uc controller.IPacksUC) {
	handler.Debug = cfg.Debug

	// Options
	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())

	// Swagger
	if cfg.EnableSwagger {
		handler.GET("/swagger/*", echoSwagger.WrapHandler)
	}
}
