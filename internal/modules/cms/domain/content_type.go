package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SchemaDefinitionMap map[string]SchemaDefinition

type ContentType struct {
	Id        uuid.UUID `json:"id" db:"id"`
	AccountId uuid.UUID `json:"accountId" db:"account_id"`

	Name             string              `json:"name" db:"name"`
	Slug             string              `json:"slug" db:"slug"`
	Description      string              `json:"description" db:"description"`
	SchemaDefinition SchemaDefinitionMap `json:"schemaDefinition" db:"schema_definition"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type ColumnType string

const (
	ShortTextColumn ColumnType = "short-text"
	LongTextColumn  ColumnType = "long-text"
	MarkdownColumn  ColumnType = "markdown"
	NumberColumn    ColumnType = "number"
	BooleanColumn   ColumnType = "bool"

	FileColumn     ColumnType = "file"
	DateColumn     ColumnType = "date"
	DateTimeColumn ColumnType = "date-time"

	JsonColumn      ColumnType = "json"
	ReferenceColumn ColumnType = "reference"
)

const SHORT_TEXT_MAX_LENGTH = 255

func (c ColumnType) IsValid() bool {
	switch c {
	case ShortTextColumn,
		LongTextColumn,
		MarkdownColumn,
		NumberColumn,
		BooleanColumn,
		FileColumn,
		DateColumn,
		DateTimeColumn,
		JsonColumn,
		ReferenceColumn:

		return true
	}
	return false
}

type ColumnDefinition int

const (
	SingleValuedColumn ColumnDefinition = iota
	ListValuedColumn
	columnDefinitionCount // intenal count
)

func (cd ColumnDefinition) IsValid() bool {
	return cd >= SingleValuedColumn && cd < columnDefinitionCount
}

type SchemaDefinition struct {
	ColumnType       ColumnType       `json:"columnType"`
	ColumnDefinition ColumnDefinition `json:"columnDefinition"` // SingleValuedColumn or ListValuedColumn
	DefaultValue     any              `json:"defaultValue,omitempty"`
	Required         bool             `json:"required"`
	Description      string           `json:"description,omitempty"`
	// Metadata         any              `json:"metadata,omitempty"`
}

func (sd *SchemaDefinition) IsValid() bool {
	if !sd.ColumnType.IsValid() {
		return false
	}

	if !sd.ColumnDefinition.IsValid() {
		return false
	}

	if sd.DefaultValue != nil {
		if err := sd.ValidateAny(sd.DefaultValue); err != nil {
			return false
		}
	}

	return true
}

func (sd *SchemaDefinition) ValidateAny(v any) error {
	var value any
	isMulti := sd.ColumnDefinition == ListValuedColumn

	if isMulti {
		arr, ok := v.([]any)
		if !ok {
			return errors.New("needs to be an array")
		}

		// If its empty, value is an empty array which is valid
		if len(arr) == 0 {
			return nil
		} else {
			value = arr[0]
		}
	} else {
		value = v
	}

	switch sd.ColumnType {
	case BooleanColumn:
		_, ok := value.(bool)
		if !ok {
			return errors.New("needs to be a boolean")
		}

	case ShortTextColumn:
		if str, ok := value.(string); !ok {
			return errors.New("needs to be a string")
		} else {
			if len(str) > SHORT_TEXT_MAX_LENGTH {
				return fmt.Errorf("should be less than %d characters", SHORT_TEXT_MAX_LENGTH)
			}
		}
	case LongTextColumn:
		if _, ok := value.(string); !ok {
			return errors.New("needs to be a string")
		}
	case MarkdownColumn:
		if _, ok := value.(string); !ok {
			return errors.New("needs to be a string")
		}
	case NumberColumn:
		if _, ok := value.(json.Number); !ok {
			return errors.New("needs to be a number")
		}
	// TODO: For now, i'm thinking i'll store both UUIDs and strings in this column.
	// Strings will be urls and uuids will be references to a files table (not yet in schemas)
	case FileColumn:
		if _, ok := value.(string); !ok {
			return errors.New("needs to be a string")
		}
	case DateColumn:
		if _, ok := value.(time.Time); !ok {
			return errors.New("needs to be a time")
		}
	case DateTimeColumn:
		if _, ok := value.(time.Time); !ok {
			return errors.New("needs to be a time")
		}
	// JSON Marshallable Bytes
	case JsonColumn:
		if _, err := json.Marshal(value); err != nil {
			return errors.New("needs to be valid json: " + err.Error())
		}

	case ReferenceColumn:
		if _, ok := v.(uuid.UUID); !ok {
			if str, ok := v.(string); !ok {
				return errors.New("needs to be of type UUID or string")
			} else {
				if _, err := uuid.Parse(str); err != nil {
					return errors.New("needs to be valid UUID")
				}
			}
		}
	default:
		return errors.New("invalid column type: " + string(sd.ColumnType))
	}

	return nil
}

func (ct *ContentType) Validate() error {
	errFields := []string{}
	for i, v := range ct.SchemaDefinition {
		if !v.IsValid() {
			errFields = append(errFields, i)
		}
	}
	if len(errFields) != 0 {
		return fmt.Errorf("invalid schema definitions for: %v", errFields)
	}
	return nil
}

// Sql.Scanner and Sql.Valuer implementations
func (sdm SchemaDefinitionMap) Value() (driver.Value, error) {
	if len(sdm) == 0 {
		return nil, nil
	}
	return json.Marshal(sdm)
}
func (sdm *SchemaDefinitionMap) Scan(src any) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("failed to assert into []byte. corrupt value")
	}

	err := json.Unmarshal(source, sdm)
	if err != nil {
		return errors.New("failed to unmarshal string into schema definition. corrupt value")
	}
	return nil
}
