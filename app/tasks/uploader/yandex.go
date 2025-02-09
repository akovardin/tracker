package uploader

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
)

type UploadError struct {
	ErrorType string `json:"error_type"`
	Message   string `json:"message"`
}
type UploadResponse struct {
	Errors  []UploadError `json:"errors"`
	Code    int           `json:"code"`
	Message string        `json:"message"`
}

func (t *Task) yandex(tracker *models.Record) error {
	to := time.Now().Add(-time.Hour * 4).Format("2006-01-02 15:04:05")
	records, err := t.app.Dao().FindRecordsByFilter(
		"conversions",
		"uploaded = false && network = 'yandex' && tracker = {:tracker} && created < {:to}",
		"-created",
		100, // limit
		0,
		dbx.Params{"tracker": tracker.Id, "to": to},
	)

	if err != nil {
		return err
	}

	if len(records) == 0 {
		t.app.Logger().Warn("empty yandex conversion")
		return nil
	}

	yaurl := tracker.GetString("yaurl")
	yatoken := tracker.GetString("yatoken")

	// create conversions file data
	data := &bytes.Buffer{}
	// create writer for conversions file data
	file := csv.NewWriter(data)
	if err := file.Write([]string{
		//"ClientId",
		"Yclid",
		"Target",
		"DateTime",
	}); err != nil {
		return err
	}

	for _, record := range records {
		if err := file.Write([]string{
			//item.ClientId,
			record.GetString("yclid"),
			"app_install",
			strconv.Itoa(int(record.Created.Time().In(t.loc).Unix())),
		}); err != nil {
			t.app.Logger().Warn("error on write csv row", "error", err)
			continue
		}
	}

	file.Flush()

	// create all body
	body := &bytes.Buffer{}
	// create writer for body
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("file", "file.csv")
	if _, err := io.Copy(part, data); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	request, _ := http.NewRequest("POST", yaurl, body)
	request.Header.Add("Authorization", "OAuth "+yatoken)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		return err
	}

	_ = dump // debug here

	resp, err := t.client.Do(request)
	if err != nil {
		return err
	}

	dump, err = httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}

	_ = dump // debug here

	result := UploadResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if len(result.Errors) != 0 {
		return fmt.Errorf("error on upload file: %s", result.Message)
	}

	for _, record := range records {
		record.Set("uploaded", true)

		if err := t.app.Dao().Save(record); err != nil {
			t.app.Logger().Warn("error save uploaded conversions", "error", err)
		}
	}

	return nil
}
