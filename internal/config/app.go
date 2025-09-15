package config

import "github.com/akfaiz/go-vue-starter-kit/pkg/env"

type App struct {
	Name             string
	Key              string
	Env              string
	FrontendBaseURL  string
	FrontendProxyURL string // development only
	LogLevel         string
	LogFormat        string
}

func loadAppConfig() App {
	return App{
		Name:             env.GetString("APP_NAME", "gova"),
		Key:              env.MustGetString("APP_KEY"),
		Env:              env.GetString("APP_ENV", "development"),
		FrontendBaseURL:  env.GetString("FRONTEND_BASE_URL", "http://localhost:3000"),
		FrontendProxyURL: env.GetString("FRONTEND_PROXY_URL", "http://localhost:5173"),
		LogLevel:         env.GetString("LOG_LEVEL", "debug"),
		LogFormat:        env.GetString("LOG_FORMAT", "json"),
	}
}
