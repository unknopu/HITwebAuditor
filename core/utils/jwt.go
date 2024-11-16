package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"auditor/core/context"
)

// GetJWTClaims get jwt data
func GetJWTClaims(c echo.Context) *context.Claims {
	u := c.(*context.Context).Get("user")
	if u != nil {
		user := u.(*jwt.Token)
		if cl, ok := user.Claims.(*context.Claims); ok {
			return cl
		}
	}
	return nil
}
