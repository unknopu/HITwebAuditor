package middleware

import (
	"auditor/core/context"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// CustomContext custom context
func CustomContext(bundle *i18n.Bundle) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			locale := c.QueryParam("locale")

			if locale != "th" {
				locale = "en"
			} else {
				locale = "th"
			}

			l := i18n.NewLocalizer(bundle, locale)
			session := c.Request().Header.Get("session")
			mobile := c.Request().Header.Get("mobile")
			platform := c.Request().Header.Get("platform")
			cc := &context.Context{Context: c, Localizer: l, Locale: locale, SessionID: &session, Mobile: mobile, Platform: platform}
			return next(cc)
		}
	}
}
