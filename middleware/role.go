package middleware

import (
	"auditor/core/context"
	"auditor/core/utils"
	"auditor/entities"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RequiredRoles required admin
func RequiredRoles(roles ...entities.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := utils.GetJWTClaims(c.(*context.Context))
			for _, role := range roles {
				if entities.Role(claims.Role) == role {
					return next(c)
				}
			}
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			}
		}
	}
}
