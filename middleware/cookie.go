package middleware

import (
	"auditor/core/context"
	"auditor/core/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// Cookie cookie
func Cookie(secret string, required bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			at := utils.GetTokenCookie(c)
			if at == "" {
				h := c.Request().Header.Get("Authorization")
				splitToken := strings.Split(h, "Bearer")
				if len(splitToken) >= 2 {
					at = strings.TrimSpace(splitToken[1])
				}
			}
			t, err := getToken(at, secret)
			if err == nil && t.Valid {
				utils.SaveCookie(c, at)
				c.Set("user", t)
				return next(c)
			}
			if required {
				return &echo.HTTPError{
					Code:     http.StatusUnauthorized,
					Message:  "Unauthorized",
					Internal: err,
				}
			}
			return next(c)
		}
	}
}

func getToken(tokenString, secret string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &context.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}
