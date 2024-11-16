package sql

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/handlers/common"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
)

// Interface me interface
type Interface interface {
	TestIntruder(c echo.Context) error
	Init(c echo.Context) error
	ErrorBased(c echo.Context) error
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

// TestDB test db
func (h *Handler) TestIntruder(c echo.Context) error {
	cc := c.(*context.Context)
	f := &common.PageQuery{}
	if err := c.(*context.Context).BindAndValidate(f); err != nil {
		return err
	}
	u, err := h.ms.TestIntruder(cc, f)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}

func (h *Handler) Init(c echo.Context) error {
	cc := c.(*context.Context)
	f := &BaseForm{}
	if err := c.(*context.Context).BindAndValidate(f); err != nil {
		return err
	}
	u, err := h.ms.Init(cc, f)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}

func (h *Handler) ErrorBased(c echo.Context) error {
	cc := c.(*context.Context)
	f := &BaseForm{}
	if err := c.(*context.Context).BindAndValidate(f); err != nil {
		return err
	}
	u, err := h.ms.ErrorBased(cc, f)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}
