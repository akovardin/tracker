package uploader

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
)

func (t *Task) vk(tracker *models.Record) error {
	records, err := t.app.Dao().FindRecordsByFilter(
		"conversions",
		"uploaded = false && network = 'vk' && landing = {:landing}",
		"-created",
		10, // limit
		0,
		dbx.Params{"landing": tracker.Id},
	)

	_ = records

	return err
}
