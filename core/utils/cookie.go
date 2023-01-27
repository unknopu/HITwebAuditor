package utils

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/labstack/echo/v4"
)

const (
	jwtSession = "BOOTSSESSIONID"
)

var s *securecookie.SecureCookie

// InitSecureCookie init securecookie
func InitSecureCookie(hashKey, blockKey string) {
	s = securecookie.New([]byte(hashKey), []byte(blockKey))
}

// SaveCookie save cookie
func SaveCookie(c echo.Context, a string) {
	if encoded, err := s.Encode(jwtSession, a); err == nil {
		cookie := &http.Cookie{
			Name:     jwtSession,
			Value:    encoded,
			Path:     "/",
			Secure:   false,
			HttpOnly: false,
			Expires:  time.Now().AddDate(0, 1, 0),
		}
		c.SetCookie(cookie)
	}
}

// GetTokenCookie get jwt cookie
func GetTokenCookie(c echo.Context) string {
	if cookie, err := c.Cookie(jwtSession); err == nil {
		var value string
		if err = s.Decode(jwtSession, cookie.Value, &value); err == nil {
			return value
		}
	}
	return ""
}

// DeleteCookie delete cookie
func DeleteCookie(c echo.Context) {
	cookie := &http.Cookie{
		Name:     jwtSession,
		Value:    "",
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		MaxAge:   -1,
	}
	c.SetCookie(cookie)
}
