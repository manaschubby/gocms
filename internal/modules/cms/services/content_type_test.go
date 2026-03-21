package services_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
	"github.com/manaschubby/gocms/internal/modules/cms/mocks"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
	"github.com/manaschubby/gocms/internal/modules/cms/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContentTypeService_CreateNewContentType(t *testing.T) {
	accountID := uuid.New()

	tests := []struct {
		name          string
		input         *domain.ContentType
		setupMock     func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo)
		expectedCode  int
		expectedError string
	}{
		{
			name: "Fail - Empty Schema",
			input: &domain.ContentType{
				AccountId:        accountID,
				SchemaDefinition: map[string]domain.SchemaDefinition{},
			},
			setupMock:     func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo) {},
			expectedCode:  http.StatusBadRequest,
			expectedError: "at least one schemaDefinition Required",
		},
		{
			name: "Fail - Invalid Schema Column",
			input: &domain.ContentType{
				AccountId: accountID,
				SchemaDefinition: map[string]domain.SchemaDefinition{
					"title": {ColumnType: "invalid-type"},
				},
			},
			setupMock:     func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo) {},
			expectedCode:  http.StatusBadRequest,
			expectedError: "failed to validate schema definition for title",
		},
		{
			name: "Fail - Account Not Found",
			input: &domain.ContentType{
				AccountId: accountID,
				SchemaDefinition: map[string]domain.SchemaDefinition{
					"title": {ColumnType: domain.ShortTextColumn, ColumnDefinition: domain.SingleValuedColumn},
				},
			},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(nil, nil)
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "account not found",
		},
		{
			name: "Fail - Duplicate Slug",
			input: &domain.ContentType{
				AccountId: accountID,
				Slug:      "blog-post",
				SchemaDefinition: map[string]domain.SchemaDefinition{
					"title": {ColumnType: domain.ShortTextColumn, ColumnDefinition: domain.SingleValuedColumn},
				},
			},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("GetContentTypeBySlug", mock.Anything, mock.Anything).Return(&domain.ContentType{}, nil)
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "content type with slug already exists",
		},
		{
			name: "Success - Create Content Type",
			input: &domain.ContentType{
				AccountId: accountID,
				Slug:      "valid-slug",
				SchemaDefinition: map[string]domain.SchemaDefinition{
					"title": {ColumnType: domain.ShortTextColumn, ColumnDefinition: domain.SingleValuedColumn},
				},
			},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("GetContentTypeBySlug", mock.Anything, mock.Anything).Return(nil, nil)
				cm.On("CreateNewContentType", mock.Anything, mock.Anything).Return(nil)
			},
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockAccount := new(mocks.MockAccountRepo)
			mockContent := new(mocks.MockContentTypeRepo)

			// We wrap them in the Repositories struct
			repos := repository.Repositories{
				Account:     mockAccount,
				ContentType: mockContent,
			}

			svc := services.NewContentTypeService(repos)
			tt.setupMock(mockAccount, mockContent)

			// Act
			code, err := svc.CreateNewContentType(context.Background(), tt.input)

			// Assert
			assert.Equal(t, tt.expectedCode, code)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockAccount.AssertExpectations(t)
			mockContent.AssertExpectations(t)
		})
	}
}
