package middleware

import (
	"auditor/core/utils"
	"auditor/entities"
	"auditor/handlers/logs"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ActivityLog wrap activity log
func ActivityLog() echo.MiddlewareFunc {
	ls := logs.NewService()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cErr := next(c)
			claims := utils.GetJWTClaims(c)
			if claims != nil && claims.Subject != "" {
				uID, err := primitive.ObjectIDFromHex(claims.Subject)
				if err == nil {
					l := entities.Log{
						Uri:            c.Request().RequestURI,
						ClientIP:       c.RealIP(),
						UserAgent:      c.Request().UserAgent(),
						UserID:         &uID,
						HttpStatusCode: c.Response().Status,
					}
					ls.Insert(&l)
				}
			}
			return cErr
		}
	}
}
