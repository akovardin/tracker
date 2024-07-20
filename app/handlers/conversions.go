package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
)

const (
	networkYandex = "yandex"
	networkVK     = "vk"
)

type Conversions struct {
	app *pocketbase.PocketBase
}

func NewConversions(app *pocketbase.PocketBase) *Conversions {
	return &Conversions{
		app: app,
	}
}

// /downloader/fire?client_id=&yclid=14178035617771290623&install_timestamp=1721291374&appmetrica_device_id=11918280336705624214&click_id=&transaction_id=cpi17188142678303156033&match_type=fingerprint&tracker=appmetrica_821509867285037527&rb_clickid=
func (h *Conversions) Fire(c echo.Context) error {
	h.app.Logger().Info("conversion fire request", "url", c.Request().URL.String())

	name := c.PathParam("name")
	tracker, err := h.app.Dao().FindFirstRecordByFilter(
		"tracker",
		"name = {:name} && enabled = true",
		dbx.Params{"name": name},
	)
	if err != nil {
		return apis.NewNotFoundError("", err)
	}

	conversions, err := h.app.Dao().FindCollectionByNameOrId("conversions")
	if err != nil {
		return apis.NewNotFoundError("", err)
	}

	record := models.NewRecord(conversions)

	yclid := c.QueryParam("yclid")
	rbclickid := c.QueryParam("rb_clickid")
	key := yclid + rbclickid

	if key == "" {
		return nil
	}

	network := networkVK
	if yclid != "" {
		network = networkYandex
	}

	record.Set("yclid", yclid)
	record.Set("rb_clickid", rbclickid)
	record.Set("key", key)
	record.Set("uploaded", false)
	record.Set("network", network)
	record.Set("tracker", tracker.Id)

	if err := h.app.Dao().SaveRecord(record); err != nil {
		h.app.Logger().Error("error on save conversions", "error", err)
		return apis.NewApiError(http.StatusInternalServerError, "error on save", err)
	}

	return nil
}
