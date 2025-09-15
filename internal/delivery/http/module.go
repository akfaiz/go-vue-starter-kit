package http

import (
	"io/fs"

	"github.com/akfaiz/go-vue-starter-kit/internal/delivery/http/handler"
	"github.com/akfaiz/go-vue-starter-kit/internal/delivery/http/middleware"
	"github.com/akfaiz/go-vue-starter-kit/internal/delivery/http/routes"
	"github.com/akfaiz/go-vue-starter-kit/internal/delivery/http/server"
	"github.com/akfaiz/go-vue-starter-kit/web"
	"go.uber.org/fx"
)

var Module = fx.Module("http",
	fx.Provide(
		server.New,
		middleware.New,
		DistFS,
	),
	fx.Invoke(routes.Register),
	handler.Module,
)

func DistFS() (fs.FS, error) {
	return fs.Sub(web.Dist, "dist")
}
