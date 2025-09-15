package serve

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/akfaiz/go-vue-starter-kit/internal/config"
	"github.com/akfaiz/go-vue-starter-kit/internal/db"
	deliveryhttp "github.com/akfaiz/go-vue-starter-kit/internal/delivery/http"
	"github.com/akfaiz/go-vue-starter-kit/internal/gateway"
	"github.com/akfaiz/go-vue-starter-kit/internal/hash"
	"github.com/akfaiz/go-vue-starter-kit/internal/lang"
	"github.com/akfaiz/go-vue-starter-kit/internal/logger"
	"github.com/akfaiz/go-vue-starter-kit/internal/repository"
	"github.com/akfaiz/go-vue-starter-kit/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var Command = &cli.Command{
	Name:  "serve",
	Usage: "Start the web server",
	Action: func(ctx context.Context, c *cli.Command) error {
		app, err := newApp()
		if err != nil {
			return err
		}
		app.Run()
		return nil
	},
}

func newApp() (*fx.App, error) {
	cfg := config.Load()
	options := appOptions(cfg)
	lang.Init()
	logger.Init(cfg.App)
	if err := fx.ValidateApp(options...); err != nil {
		return nil, err
	}

	return fx.New(options...), nil
}

func appOptions(cfg config.Config) []fx.Option {
	return []fx.Option{
		fx.WithLogger(func() fxevent.Logger {
			slogLogger := &fxevent.SlogLogger{Logger: slog.Default()}
			slogLogger.UseLogLevel(slog.LevelDebug)
			return slogLogger
		}),
		fx.Supply(cfg, cfg.Auth, cfg.Auth.JWT, cfg.Database),
		fx.Provide(
			db.NewDatabase,
		),
		repository.Module,
		hash.Module,
		gateway.Module,
		service.Module,
		deliveryhttp.Module,
		fx.Invoke(httpServerLifecycle),
	}
}

func httpServerLifecycle(lc fx.Lifecycle, e *echo.Echo, cfg config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := e.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
					log.Fatal(err)
				}
			}()
			log.Printf("🚀 Server started at http://localhost:%d", cfg.Server.Port)
			log.Printf("📚 OpenAPI docs at http://localhost:%d/docs", cfg.Server.Port)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})
}
