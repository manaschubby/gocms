package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
)

func TestAccountDTO(t *testing.T) {
	single := domain.Account{Id: uuid.New(), Name: "Checking"}
	slice := []domain.Account{{Id: uuid.New(), Name: "Savings"}}

	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name:    "Single Object",
			input:   single,
			wantErr: false,
		},
		{
			name:    "Pointer to Object",
			input:   &single,
			wantErr: false,
		},
		{
			name:    "Slice of Objects",
			input:   slice,
			wantErr: false,
		},
		{
			name:    "Slice of Pointers",
			input:   []*domain.Account{&single},
			wantErr: false,
		},
		{
			name:    "Invalid Type (string)",
			input:   "not an account",
			wantErr: true,
		},
		{
			name:    "Nil Input",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.AccountDTO(tt.input)

			// Check error expectations
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountDTO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we don't expect an error, verify the JSON is valid
			if !tt.wantErr {
				if len(got) == 0 {
					t.Error("AccountDTO() returned empty bytes for valid input")
				}

				// Optional: Unmarshal back to check integrity
				var check any
				if err := json.Unmarshal(got, &check); err != nil {
					t.Errorf("AccountDTO() produced invalid JSON: %v", err)
				}
			}
		})
	}
}
