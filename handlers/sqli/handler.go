package sqli

import (
	"auditor/app"
	"auditor/core/context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
)

// Interface me interface
type Interface interface {
	Init(c echo.Context) error
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

func (h *Handler) Init(c echo.Context) error {
	cc := c.(*context.Context)
	f := &SqliForm{}
	if err := c.(*context.Context).BindAndValidate(f); err != nil {
		return err
	}
	u, err := h.ms.Init(cc, f)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}
