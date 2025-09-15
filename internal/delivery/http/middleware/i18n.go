package middleware

import (
	"github.com/invopop/ctxi18n"
	"github.com/labstack/echo/v4"
)

const defaultLocale = "en"

func I18n() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			acceptLanguage := c.Request().Header.Get("Accept-Language")
			if acceptLanguage == "" {
				acceptLanguage = defaultLocale
			}
			ctx := c.Request().Context()
			ctx, err := ctxi18n.WithLocale(ctx, acceptLanguage)
			if err != nil {
				ctx, _ = ctxi18n.WithLocale(ctx, defaultLocale)
			}
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
