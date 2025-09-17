package config

import "github.com/akfaiz/go-vue-starter-kit/pkg/env"

type App struct {
	Name             string
	Key              string
	Env              string
	ApiBaseURL       string
	FrontendBaseURL  string
	FrontendProxyURL string // development only
	LogLevel         string
	LogFormat        string
}

func loadAppConfig() App {
	return App{
		Name:             env.GetString("APP_NAME", "go-vue-starter-kit"),
		Key:              env.MustGetString("APP_KEY"),
		Env:              env.GetString("APP_ENV", "development"),
		ApiBaseURL:       env.GetString("API_BASE_URL", "http://localhost:8080/api"),
		FrontendBaseURL:  env.GetString("FRONTEND_BASE_URL", "http://localhost:8080"),
		FrontendProxyURL: env.GetString("FRONTEND_PROXY_URL", "http://localhost:5173"),
		LogLevel:         env.GetString("LOG_LEVEL", "debug"),
		LogFormat:        env.GetString("LOG_FORMAT", "json"),
	}
}
