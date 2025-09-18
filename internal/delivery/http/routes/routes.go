package routes

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/akfaiz/go-vue-starter-kit/internal/config"
	"github.com/akfaiz/go-vue-starter-kit/internal/delivery/http/handler"
	"github.com/akfaiz/go-vue-starter-kit/internal/delivery/http/handler/dto"
	"github.com/labstack/echo/v4"
	"github.com/oaswrap/spec/adapter/echoopenapi"
	"github.com/oaswrap/spec/option"
	"go.uber.org/fx"
)

type RouteConfig struct {
	fx.In
	Echo   *echo.Echo
	Config config.Config

	AuthMiddleware echo.MiddlewareFunc `name:"auth"`

	AuthHandler        *handler.AuthHandler
	ProfileHandler     *handler.ProfileHandler
	HealthCheckHandler *handler.HealthCheckHandler
	SPAHandler         *handler.SPAHandler
}

func Register(rc RouteConfig) {
	rc.Echo.GET("/health", rc.HealthCheckHandler.HealthCheck)

	rc.registerAPI()
	rc.registerWeb()
}

func (rc RouteConfig) registerWeb() {
	if os.Getenv("DEV") == "1" {
		proxyURL, err := url.Parse(rc.Config.App.FrontendProxyURL)
		if err != nil {
			log.Fatalf("Invalid FRONTEND_PROXY_URL: %v", err)
		}
		proxy := httputil.NewSingleHostReverseProxy(proxyURL)
		rc.Echo.Any("/*", echo.WrapHandler(proxy))
		log.Println("🚀 DEV=1 → proxying frontend to", rc.Config.App.FrontendProxyURL)
	} else {
		rc.Echo.GET("/env.js", rc.SPAHandler.Env)
		rc.Echo.GET("/*", rc.SPAHandler.View)
	}
}

func (rc RouteConfig) registerAPI() {
	r := echoopenapi.NewRouter(rc.Echo,
		option.WithTitle("Go Vue Starter Kit API"),
		option.WithVersion("1.0.0"),
		option.WithDescription("API documentation for Go Vue Starter Kit"),
		option.WithReflectorConfig(
			option.StripDefNamePrefix("Dto"),
		),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
		option.WithServer("https://go-vue-starter-kit.fly.dev/", option.ServerDescription("Production server")),
		option.WithServer(fmt.Sprintf("http://localhost:%d", rc.Config.Server.Port), option.ServerDescription("Local server")),
	)
	v1 := r.Group("/api/v1")

	auth := v1.Group("/auth").With(option.GroupTags("Authentication"))
	auth.POST("/login", rc.AuthHandler.Login).With(
		option.Summary("User Login"),
		option.Description("Authenticate user and return access and refresh tokens"),
		option.Request(new(dto.LoginRequest)),
		option.Response(200, responseOf(dto.TokenResponse{})),
	)
	auth.POST("/register", rc.AuthHandler.Register).With(
		option.Summary("User Registration"),
		option.Description("Register a new user"),
		option.Request(new(dto.RegisterRequest)),
		option.Response(201, responseOf(dto.TokenResponse{})),
	)
	auth.POST("/refresh-token", rc.AuthHandler.RefreshToken).With(
		option.Summary("Refresh Token"),
		option.Description("Refresh access and refresh tokens using a valid refresh token"),
		option.Request(new(dto.RefreshTokenRequest)),
		option.Response(200, responseOf(dto.TokenResponse{})),
	)
	auth.POST("/forgot-password", rc.AuthHandler.SendForgotPasswordEmail).With(
		option.Summary("Send Forgot Password Email"),
		option.Description("Send a password reset email to the user"),
		option.Request(new(dto.SendForgotPasswordEmailRequest)),
		option.Response(200, responseOf[any](nil)),
	)
	auth.POST("/validate-reset-password", rc.AuthHandler.ValidateResetPassword).With(
		option.Summary("Validate Reset Password Token"),
		option.Description("Validate the password reset token and email"),
		option.Request(new(dto.ValidateResetPasswordRequest)),
		option.Response(200, responseOf[any](nil)),
	)
	auth.POST("/reset-password", rc.AuthHandler.ResetPassword).With(
		option.Summary("Reset Password"),
		option.Description("Reset the user's password using a valid token"),
		option.Request(new(dto.ResetPasswordRequest)),
		option.Response(200, responseOf[any](nil)),
	)
	auth.POST("/email/send-verification", rc.AuthHandler.SendVerificationEmail, rc.AuthMiddleware).With(
		option.Summary("Send Verification Email"),
		option.Description("Send an email verification link to the user"),
		option.Response(200, responseOf[any](nil)),
		option.Security("bearerAuth"),
	)
	auth.POST("/email/verify", rc.AuthHandler.VerifyEmail, rc.AuthMiddleware).With(
		option.Summary("Verify Email"),
		option.Description("Verify the user's email using a valid token"),
		option.Request(new(dto.VerifyEmailRequest)),
		option.Response(200, responseOf[any](nil)),
		option.Security("bearerAuth"),
	)

	profile := v1.Group("/profile", rc.AuthMiddleware).With(
		option.GroupTags("Profile"),
		option.GroupSecurity("bearerAuth"),
	)
	profile.GET("", rc.ProfileHandler.GetProfile).With(
		option.Summary("Get User Profile"),
		option.Description("Retrieve the profile information of the authenticated user"),
		option.Response(200, responseOf(dto.ProfileResponse{})),
	)
	profile.PUT("", rc.ProfileHandler.UpdateProfile).With(
		option.Summary("Update User Profile"),
		option.Description("Update the profile information of the authenticated user"),
		option.Request(new(dto.UpdateProfileRequest)),
		option.Response(200, responseOf(dto.ProfileResponse{})),
	)
	profile.DELETE("", rc.ProfileHandler.DeleteAccount).With(
		option.Summary("Delete User Account"),
		option.Description("Delete the authenticated user's account"),
		option.Request(new(dto.DeleteAccountRequest)),
		option.Response(200, responseOf[any](nil)),
	)
	profile.PUT("/password", rc.ProfileHandler.ChangePassword).With(
		option.Summary("Change Password"),
		option.Description("Change the password of the authenticated user"),
		option.Request(new(dto.ChangePasswordRequest)),
		option.Response(200, responseOf[any](nil)),
	)
}

func responseOf[T any](model T) any {
	return struct {
		dto.Response[T]
	}{
		Response: dto.Response[T]{Data: model},
	}
}
