package me

import (
	"auditor/app"
	"auditor/core/context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
)

// Interface me interface
type Interface interface {
	GetMe(c echo.Context) error
}

// Handler me handler
type Handler struct {
	ms    ServiceInterface
	cache *cache.Cache
}

// NewHandler new handler
func NewHandler(c *app.Context) Interface {
	return &Handler{
		ms: NewService(c),
	}
}

// GetMe get me
func (h *Handler) GetMe(c echo.Context) error {
	cc := c.(*context.Context)
	// claims := utils.GetJWTClaims(c)
	u, err := h.ms.TEST(cc)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}
