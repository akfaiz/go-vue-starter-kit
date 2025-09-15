package middleware

import (
	"github.com/akfaiz/go-vue-starter-kit/internal/delivery/http/middleware/auth"
	"github.com/akfaiz/go-vue-starter-kit/internal/domain"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type Middleware struct {
	fx.Out

	Auth echo.MiddlewareFunc `name:"auth"`
}

type MiddlewareConfig struct {
	fx.In

	JWTManager domain.JWTManager
}

func New(cfg MiddlewareConfig) Middleware {
	return Middleware{
		Auth: auth.New(cfg.JWTManager),
	}
}
