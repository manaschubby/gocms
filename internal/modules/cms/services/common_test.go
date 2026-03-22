package services_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
	"github.com/manaschubby/gocms/internal/modules/cms/mocks"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
	"github.com/manaschubby/gocms/internal/modules/cms/services" // or wherever validateAccount is defined
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	EXTREMELY_LONG_STRING = "EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING_EXTREMELY_LONG_STRING"
)

func TestValidateAccount(t *testing.T) {
	aid := uuid.New()
	ctx := context.Background()

	tests := []struct {
		name          string
		setupMock     func(m *mocks.MockAccountRepo)
		expectedCode  int
		expectedError string
	}{
		{
			name: "Success - Account Found",
			setupMock: func(m *mocks.MockAccountRepo) {
				m.On("GetAccountByUUID", aid, mock.MatchedBy(func(opt repository.GetAccountOptions) bool {
					return opt.Context != nil
				})).Return(&domain.Account{Id: aid}, nil)
			},
			expectedCode:  0,
			expectedError: "",
		},
		{
			name: "Fail - Repository Error",
			setupMock: func(m *mocks.MockAccountRepo) {
				m.On("GetAccountByUUID", aid, mock.Anything).
					Return(nil, errors.New("connection refused"))
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "failed to retrieve account data",
		},
		{
			name: "Fail - Account Not Found (Nil)",
			setupMock: func(m *mocks.MockAccountRepo) {
				// Mock returning (nil, nil) which your function handles as "not found"
				m.On("GetAccountByUUID", aid, mock.Anything).Return(nil, nil)
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "account not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(mocks.MockAccountRepo)
			tt.setupMock(mockRepo)

			// Act
			code, err := services.ValidateAccount(mockRepo, ctx, aid)

			// Assert
			assert.Equal(t, tt.expectedCode, code)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
