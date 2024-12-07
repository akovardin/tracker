package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/tools/template"

	"gohome.4gophers.ru/kovardin/vzor/views"
)

type Home struct {
	registry *template.Registry
}

func NewHome(registry *template.Registry) *Home {
	return &Home{
		registry: registry,
	}
}

func (h *Home) Home(c echo.Context) error {
	html, err := h.registry.LoadFS(views.FS,
		"layout.html",
		"home/page.html",
	).Render(map[string]any{
		"slug": "",
	})

	if err != nil {
		return apis.NewNotFoundError("", err)
	}

	return c.HTML(http.StatusOK, html)
}
