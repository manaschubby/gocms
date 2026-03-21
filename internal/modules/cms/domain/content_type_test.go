package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColumnType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		ct       domain.ColumnType
		expected bool
	}{
		{"valid short text", domain.ShortTextColumn, true},
		{"valid reference", domain.ReferenceColumn, true},
		{"invalid type", domain.ColumnType("unknown"), false},
		{"empty type", domain.ColumnType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.ct.IsValid())
		})
	}
}

func TestSchemaDefinition_ValidateAny(t *testing.T) {
	tests := []struct {
		name    string
		schema  domain.SchemaDefinition
		input   any
		wantErr bool
	}{
		// --- ShortText Tests ---
		{
			name:    "ShortText valid",
			schema:  domain.SchemaDefinition{ColumnType: domain.ShortTextColumn, ColumnDefinition: domain.SingleValuedColumn},
			input:   "hello world",
			wantErr: false,
		},
		{
			name:    "ShortText too long",
			schema:  domain.SchemaDefinition{ColumnType: domain.ShortTextColumn, ColumnDefinition: domain.SingleValuedColumn},
			input:   string(make([]byte, domain.SHORT_TEXT_MAX_LENGTH+1)),
			wantErr: true,
		},
		{
			name:    "ShortText invalid type",
			schema:  domain.SchemaDefinition{ColumnType: domain.ShortTextColumn, ColumnDefinition: domain.SingleValuedColumn},
			input:   123,
			wantErr: true,
		},

		// --- Number Tests (json.Number) ---
		{
			name:    "Number valid",
			schema:  domain.SchemaDefinition{ColumnType: domain.NumberColumn, ColumnDefinition: domain.SingleValuedColumn},
			input:   json.Number("42"),
			wantErr: false,
		},
		{
			name:    "Number invalid raw int",
			schema:  domain.SchemaDefinition{ColumnType: domain.NumberColumn, ColumnDefinition: domain.SingleValuedColumn},
			input:   42, // ValidateAny expects json.Number specifically
			wantErr: true,
		},

		// --- List Valued Tests ---
		{
			name:    "List of Booleans valid",
			schema:  domain.SchemaDefinition{ColumnType: domain.BooleanColumn, ColumnDefinition: domain.ListValuedColumn},
			input:   []any{true, false, true},
			wantErr: false,
		},
		{
			name:    "List invalid non-array input",
			schema:  domain.SchemaDefinition{ColumnType: domain.BooleanColumn, ColumnDefinition: domain.ListValuedColumn},
			input:   true,
			wantErr: true,
		},
		{
			name:    "List empty array is valid",
			schema:  domain.SchemaDefinition{ColumnType: domain.BooleanColumn, ColumnDefinition: domain.ListValuedColumn},
			input:   []any{},
			wantErr: false,
		},

		// --- Reference Tests ---
		{
			name:    "Reference valid UUID string",
			schema:  domain.SchemaDefinition{ColumnType: domain.ReferenceColumn, ColumnDefinition: domain.SingleValuedColumn},
			input:   uuid.New().String(),
			wantErr: false,
		},
		{
			name:    "Reference valid UUID object",
			schema:  domain.SchemaDefinition{ColumnType: domain.ReferenceColumn, ColumnDefinition: domain.SingleValuedColumn},
			input:   uuid.New(),
			wantErr: false,
		},
		{
			name:    "Reference invalid string",
			schema:  domain.SchemaDefinition{ColumnType: domain.ReferenceColumn, ColumnDefinition: domain.SingleValuedColumn},
			input:   "not-a-uuid",
			wantErr: true,
		},

		// --- Time Tests ---
		{
			name:    "DateTime valid",
			schema:  domain.SchemaDefinition{ColumnType: domain.DateTimeColumn, ColumnDefinition: domain.SingleValuedColumn},
			input:   time.Now(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.schema.ValidateAny(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSchemaDefinition_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		sd       domain.SchemaDefinition
		expected bool
	}{
		{
			name: "Empty Default Value",
			sd: domain.SchemaDefinition{
				ColumnType:       domain.ShortTextColumn,
				ColumnDefinition: domain.SingleValuedColumn,
			},
			expected: true,
		},
		{
			name: "Valid ShortText",
			sd: domain.SchemaDefinition{
				ColumnType:       domain.ShortTextColumn,
				ColumnDefinition: domain.SingleValuedColumn,
				DefaultValue:     "Hello",
			},
			expected: true,
		},
		{
			name: "Invalid ColumnType",
			sd: domain.SchemaDefinition{
				ColumnType:       "wrong",
				ColumnDefinition: domain.SingleValuedColumn,
			},
			expected: false,
		},
		{
			name: "Invalid ColumnDefinition",
			sd: domain.SchemaDefinition{
				ColumnType:       domain.ShortTextColumn,
				ColumnDefinition: domain.ColumnDefinition(99),
			},
			expected: false,
		},
		{
			name: "ShortText Too Long",
			sd: domain.SchemaDefinition{
				ColumnType:       domain.ShortTextColumn,
				ColumnDefinition: domain.SingleValuedColumn,
				DefaultValue:     string(make([]byte, domain.SHORT_TEXT_MAX_LENGTH+1)),
			},
			expected: false,
		},
		{
			name: "Valid ListValued",
			sd: domain.SchemaDefinition{
				ColumnType:       domain.BooleanColumn,
				ColumnDefinition: domain.ListValuedColumn,
				DefaultValue:     []any{true, false},
			},
			expected: true,
		},
		{
			name: "Invalid ListValued - Not an array",
			sd: domain.SchemaDefinition{
				ColumnType:       domain.BooleanColumn,
				ColumnDefinition: domain.ListValuedColumn,
				DefaultValue:     true, // Should be []any
			},
			expected: false,
		},
		{
			name: "Valid Reference UUID Object",
			sd: domain.SchemaDefinition{
				ColumnType:       domain.ReferenceColumn,
				ColumnDefinition: domain.SingleValuedColumn,
				DefaultValue:     uuid.New(),
			},
			expected: true,
		},
		{
			name: "Valid Reference UUID String",
			sd: domain.SchemaDefinition{
				ColumnType:       domain.ReferenceColumn,
				ColumnDefinition: domain.SingleValuedColumn,
				DefaultValue:     uuid.New().String(),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.sd.IsValid())
		})
	}
}

func TestContentType_Validate(t *testing.T) {
	t.Run("Valid ContentType", func(t *testing.T) {
		ct := &domain.ContentType{
			SchemaDefinition: map[string]domain.SchemaDefinition{
				"title": {
					ColumnType:       domain.ShortTextColumn,
					ColumnDefinition: domain.SingleValuedColumn,
					DefaultValue:     "New Post",
				},
			},
		}
		assert.NoError(t, ct.Validate())
	})

	t.Run("Invalid Schema Field", func(t *testing.T) {
		ct := &domain.ContentType{
			SchemaDefinition: map[string]domain.SchemaDefinition{
				"broken_field": {
					ColumnType: "non-existent",
				},
			},
		}
		err := ct.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "broken_field")
	})

	t.Run("Empty Default Value", func(t *testing.T) {
		ct := &domain.ContentType{
			SchemaDefinition: map[string]domain.SchemaDefinition{
				"title": {
					ColumnType:       domain.ShortTextColumn,
					ColumnDefinition: domain.SingleValuedColumn,
				},
			},
		}
		assert.NoError(t, ct.Validate())
	})
}

func TestSchemaDefinitionMap_Value(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.SchemaDefinitionMap
		want    any // driver.Value is an alias for any
		wantErr bool
	}{
		{
			name:    "Success - Valid Map",
			input:   domain.SchemaDefinitionMap{"title": {ColumnType: domain.ShortTextColumn, ColumnDefinition: 0, Required: false}},
			want:    []byte(`{"title": {"columnDefinition":0, "columnType":"short-text", "required":false}}`), // Adjust keys based on your actual JSON tags
			wantErr: false,
		},
		{
			name:    "Success - Empty Map returns nil",
			input:   domain.SchemaDefinitionMap{},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Success - Nil Map returns nil",
			input:   nil,
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Value()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.want == nil {
					assert.Nil(t, got)
				} else {
					// We compare as JSON strings to avoid byte slice pointer issues
					assert.JSONEq(t, string(tt.want.([]byte)), string(got.([]byte)))
				}
			}
		})
	}
}

func TestSchemaDefinitionMap_Scan(t *testing.T) {
	validJSON := []byte(`{"title":{"columnType":"short-text"}}`)

	tests := []struct {
		name          string
		src           any
		expectedState domain.SchemaDefinitionMap
		expectedErr   string
	}{
		{
			name: "Success - Valid Byte Slice",
			src:  validJSON,
			expectedState: domain.SchemaDefinitionMap{
				"title": {ColumnType: domain.ShortTextColumn},
			},
		},
		{
			name:        "Fail - Invalid Type (String instead of []byte)",
			src:         `{"title":"value"}`,
			expectedErr: "failed to assert into []byte",
		},
		{
			name:        "Fail - Corrupt JSON",
			src:         []byte(`{"title": { "broken" ...`),
			expectedErr: "failed to unmarshal string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sdm domain.SchemaDefinitionMap
			err := sdm.Scan(tt.src)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedState, sdm)
			}
		})
	}
}
