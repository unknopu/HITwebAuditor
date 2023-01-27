package logs

import (
	"auditor/app"
	"auditor/core/context"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Interface user interface
type Interface interface {
	GetAllWithAdmin(c echo.Context) error
}

// Handler user handler
type Handler struct {
	us ServiceInterface
}

// NewHandler new handler
func NewHandler(c *app.Context) Interface {
	return &Handler{
		us: NewService(),
	}
}

// GetAllWithAdmin get all users with admin
func (h *Handler) GetAllWithAdmin(c echo.Context) error {
	form := &GetAllWithAdminForm{}
	if err := c.(*context.Context).BindAndValidate(form); err != nil {
		return err
	}
	p, err := h.us.GetAllWithAdmin(form)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, p)
}
