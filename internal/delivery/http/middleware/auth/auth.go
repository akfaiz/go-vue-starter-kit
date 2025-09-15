package auth

import (
	"context"

	"github.com/akfaiz/go-vue-starter-kit/internal/domain"
	"github.com/akfaiz/go-vue-starter-kit/internal/errdefs"
	"github.com/labstack/echo/v4"
)

type contextKey string

const userContextKey contextKey = "user"
const userKey = "user"

func New(jwtManager domain.JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return errdefs.ErrUnauthorized("missing Authorization header")
			}
			if len(token) < 7 || token[:7] != "Bearer " {
				return errdefs.ErrUnauthorized("Authorization header must start with 'Bearer '")
			}
			token = token[7:] // Remove "Bearer " prefix
			if token == "" {
				return errdefs.ErrUnauthorized("missing token in Authorization header")
			}

			claims, err := jwtManager.VerifyAccessToken(token)
			if err != nil {
				return err
			}

			c.Set(userKey, claims) // Set user in context for Echo

			req := c.Request()
			ctx := context.WithValue(req.Context(), userContextKey, claims)
			c.SetRequest(req.WithContext(ctx)) // Update request context

			return next(c)
		}
	}
}

func GetUser(c echo.Context) *domain.JWTClaims {
	claims, ok := c.Get(userKey).(*domain.JWTClaims)
	if !ok {
		return nil
	}
	return claims
}

func GetUserFromContext(ctx context.Context) *domain.JWTClaims {
	claims, ok := ctx.Value(userContextKey).(*domain.JWTClaims)
	if !ok {
		return nil
	}
	return claims
}
