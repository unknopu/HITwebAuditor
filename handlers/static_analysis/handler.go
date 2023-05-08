package static_analysis

import (
	"auditor/app"
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
	f := &StaticAnalysisForm{}
	if err := c.Bind(f); err != nil {
		return err
	}
	u, err := h.ms.Init(c, f)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}
