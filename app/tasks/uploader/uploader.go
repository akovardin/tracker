package uploader

import (
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase"
)

type Task struct {
	app    *pocketbase.PocketBase
	client *http.Client
	loc    *time.Location
}

func New(app *pocketbase.PocketBase) *Task {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}

	return &Task{
		app:    app,
		client: &http.Client{},
		loc:    location,
	}
}

func (t *Task) Do() {
	// upload fired conversions
	t.app.Logger().Info("start upload conversions")

	trackers, err := t.app.Dao().FindRecordsByFilter(
		"tracker",
		"enabled = true && network = 'yandex'",
		"-created",
		10, // limit
		0,
	)

	if err != nil {
		t.app.Logger().Error("error on get trackers")

		return
	}

	for _, tracker := range trackers {
		if err := t.yandex(tracker); err != nil {
			t.app.Logger().Error("error on upload yandex conversions", "error", err)
		}

		if err := t.vk(tracker); err != nil {
			t.app.Logger().Error("error on upload vk conversions", "error", err)
		}
	}

	t.app.Logger().Info("finish upload conversions")
}
