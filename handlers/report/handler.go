package report

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
	GetLatest(c echo.Context) error
	GetHistory(c echo.Context) error
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
	f := &Form{}
	if err := c.(*context.Context).BindAndValidate(f); err != nil {
		return err
	}
	u, err := h.ms.Init(cc, f)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}

func (h *Handler) GetLatest(c echo.Context) error {
	cc := c.(*context.Context)
	f := &GetLatestForm{}
	if err := c.(*context.Context).BindAndValidate(f); err != nil {
		return err
	}
	u, err := h.ms.GetLatest(cc, f)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}

func (h *Handler) GetHistory(c echo.Context) error {
	cc := c.(*context.Context)
	f := &GetLatestForm{}
	if err := c.(*context.Context).BindAndValidate(f); err != nil {
		return err
	}
	u, err := h.ms.GetHistory(cc, f)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}
