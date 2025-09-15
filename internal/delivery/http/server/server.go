package server

import (
	"log/slog"

	"github.com/akfaiz/go-vue-starter-kit/internal/delivery/http/middleware"
	"github.com/akfaiz/go-vue-starter-kit/internal/validator"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	e.Validator = validator.New()
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.HideBanner = true
	e.HidePort = true

	// logCfg := zap.NewProductionConfig()
	// logger := zap.Must(logCfg.Build())

	// Middleware
	// e.Use(middleware.Logger(logger))

	e.Pre(echomiddleware.RemoveTrailingSlash())
	e.Use(middleware.Logger(slog.Default()))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.CORS())
	e.Use(middleware.I18n())

	return e
}
