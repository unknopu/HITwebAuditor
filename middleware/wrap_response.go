package middleware

import (
	"auditor/core/context"
	"auditor/response"

	"github.com/labstack/echo/v4"
)

// WrapResponse wrap response
func WrapResponse(rr *response.Results) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				cc := c.(*context.Context)
				return rr.GetResponse(err, cc.Localizer)
			}
			return nil
		}
	}
}
