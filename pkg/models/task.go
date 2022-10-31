package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"personal/go-proxy-service/pkg/utilities"
)

const (
	StatusDone      = "Done"
	StatusInProcess = "In process"
	StatusError     = "Error"
	StatusNew       = "New"
	Response        = "Response"
	Request         = "Request"
)

type Task struct {
	Id              *utilities.JsonNullInt64  `json:"id,omitempty"`         // serial id
	CreatedOn       *utilities.JsonNullTime   `json:"created_on,omitempty"` // default now() timestamp
	UpdatedOn       *utilities.JsonNullTime   `json:"updated_on,omitempty"`
	Url             *utilities.JsonNullString `json:"url,omitempty"`
	RequestHeaders  *Header                   `json:"request_headers,omitempty"`
	Method          *utilities.JsonNullString `json:"method,omitempty"`
	Body            BodyType                  `json:"body,omitempty"`
	Status          *utilities.JsonNullString `json:"status,omitempty"`
	HttpStatusCode  *utilities.JsonNullInt32  `json:"http_status_code,omitempty"`
	ResponseHeaders ResHeaders                `json:"response_headers,omitempty"`
	Length          *utilities.JsonNullInt64  `json:"length,omitempty"`
}

type BodyType map[string]interface{}

func (b BodyType) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *BodyType) Scan(value interface{}) error {
	by, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(by, &b)
}

type Header struct {
	Authorization *utilities.JsonNullString `json:"authorization,omitempty"`
	ContentType   *utilities.JsonNullString `json:"content_type,omitempty"`
}

func (h Header) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *Header) Scan(value interface{}) error {

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &h)
}

type ResHeaders map[string][]string

func (rh ResHeaders) Value() (driver.Value, error) {
	return json.Marshal(rh)
}

func (rh *ResHeaders) Scan(value interface{}) error {

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &rh)
}
