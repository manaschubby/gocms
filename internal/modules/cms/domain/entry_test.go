package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntry_IsDifferentTo(t *testing.T) {
	baseEntry := domain.Entry{
		Id:          uuid.New(),
		Title:       "Original Title",
		Status:      domain.StatusDraft,
		ContentData: []byte(`{"key": "value"}`),
		Slug:        "original-slug",
		Version:     1,
	}

	tests := []struct {
		name     string
		compare  domain.Entry
		expected bool
	}{
		{
			name:     "Identical entries",
			compare:  baseEntry,
			expected: false,
		},
		{
			name: "Different Title",
			compare: func() domain.Entry {
				e := baseEntry
				e.Title = "New Title"
				return e
			}(),
			expected: true,
		},
		{
			name: "Different Status",
			compare: func() domain.Entry {
				e := baseEntry
				e.Status = domain.StatusPublished
				return e
			}(),
			expected: true,
		},
		{
			name: "Different ContentData",
			compare: func() domain.Entry {
				e := baseEntry
				e.ContentData = []byte(`{"key": "different"}`)
				return e
			}(),
			expected: true,
		},
		{
			name: "Different metadata (Id/Slug/Version) - Should be False",
			compare: func() domain.Entry {
				e := baseEntry
				e.Id = uuid.New()
				e.Slug = "something-else"
				e.Version = 2
				return e
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := baseEntry.IsDifferentTo(tt.compare)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestEntryStatus_Value(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.EntryStatus
		expected string
	}{
		{"Draft", domain.StatusDraft, "draft"},
		{"Published", domain.StatusPublished, "published"},
		{"Archived", domain.StatusArchived, "archived"},
		{"Empty/Invalid", domain.EntryStatus("unknown"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Value()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestEntryStatus_Scan(t *testing.T) {
	tests := []struct {
		name          string
		src           any
		expectedState domain.EntryStatus
		expectedErr   string
	}{
		{
			name:          "Success - String 'draft'",
			src:           "draft",
			expectedState: domain.StatusDraft,
		},
		{
			name:          "Success - Byte slice 'published'",
			src:           []byte("published"),
			expectedState: domain.StatusPublished,
		},
		{
			name:          "Success - Nil handling",
			src:           nil,
			expectedState: domain.EntryStatus(""), // Should remain unchanged/zero-value
		},
		{
			name:        "Fail - Invalid Enum Value",
			src:         "deleted",
			expectedErr: "invalid entry status enum value",
		},
		{
			name:        "Fail - Incompatible Type (int)",
			src:         123,
			expectedErr: "incompatible type for EntryStatus",
		},
		{
			name:        "Fail - Incompatible Type (bool)",
			src:         true,
			expectedErr: "incompatible type for EntryStatus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var es domain.EntryStatus
			err := es.Scan(tt.src)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedState, es)
			}
		})
	}
}
