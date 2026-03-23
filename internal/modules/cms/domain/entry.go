package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	Id            uuid.UUID `json:"id" db:"id"`
	ContentTypeId uuid.UUID `json:"contentTypeId" db:"content_type_id"`

	Slug  string `json:"slug" db:"slug"`
	Title string `json:"title" db:"title"`

	ContentData json.RawMessage `json:"contentData" db:"content_data"`
	Status      EntryStatus     `json:"status" db:"status"` // draft, published, archived

	Version int `json:"version" db:"version"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func (base Entry) IsDifferentTo(cmp Entry) bool {
	if string(base.ContentData) != string(cmp.ContentData) || base.Title != cmp.Title || base.Status != cmp.Status {
		return true
	}
	return false
}

type EntryStatus string

const (
	StatusDraft     EntryStatus = "draft"
	StatusPublished EntryStatus = "published"
	StatusArchived  EntryStatus = "archived"
)

func (es EntryStatus) Value() (driver.Value, error) {
	return string(es), nil
}
func (es *EntryStatus) Scan(src any) error {
	if src == nil {
		return nil
	}

	var source string
	switch t := src.(type) {
	case string:
		source = t
	case []byte:
		source = string(t)
	default:
		return errors.New("incompatible type for EntryStatus")
	}

	validStatus := EntryStatus(source)
	switch validStatus {
	case StatusDraft, StatusPublished, StatusArchived:
		*es = validStatus
		return nil
	default:
		return errors.New("invalid entry status enum value: " + source)
	}
}
