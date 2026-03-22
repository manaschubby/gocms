package services_test

import (
	"context"
	"database/sql"
	"errors"
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
			name: "Fail - Extremely Long title",
			input: &domain.ContentType{
				AccountId: accountID,
				Name:      EXTREMELY_LONG_STRING,
				SchemaDefinition: map[string]domain.SchemaDefinition{
					"title": {ColumnType: domain.ShortTextColumn, ColumnDefinition: domain.SingleValuedColumn},
				},
			},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "name too long",
		},
		{
			name: "Fail - Extremely Long Slug",
			input: &domain.ContentType{
				AccountId: accountID,
				Slug:      EXTREMELY_LONG_STRING,
				SchemaDefinition: map[string]domain.SchemaDefinition{
					"title": {ColumnType: domain.ShortTextColumn, ColumnDefinition: domain.SingleValuedColumn},
				},
			},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "slug length too long",
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
			expectedCode:  0,
			expectedError: "",
		},
		{
			name: "Fail - Repository Error on Slug Check",
			input: &domain.ContentType{
				AccountId:        accountID,
				Slug:             "error-slug",
				SchemaDefinition: map[string]domain.SchemaDefinition{"t": {ColumnType: domain.ShortTextColumn}},
			},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("GetContentTypeBySlug", mock.Anything, mock.Anything).Return(nil, errors.New("db connection lost"))
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "failed to check for existing content types",
		},
		{
			name: "Fail - Repository Error on Create",
			input: &domain.ContentType{
				AccountId:        accountID,
				Slug:             "new-slug",
				SchemaDefinition: map[string]domain.SchemaDefinition{"t": {ColumnType: domain.ShortTextColumn}},
			},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("GetContentTypeBySlug", mock.Anything, mock.Anything).Return(nil, nil)
				cm.On("CreateNewContentType", mock.Anything, mock.Anything).Return(errors.New("insert failed"))
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "failed to create new content type",
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

func TestContentTypeService_DeleteContentType(t *testing.T) {
	accountID := uuid.New()
	ctID := uuid.New()

	tests := []struct {
		name          string
		input         *domain.ContentType
		setupMock     func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo)
		expectedCode  int
		expectedError string
	}{
		{
			name:  "Success - Delete By ID",
			input: &domain.ContentType{AccountId: accountID, Id: ctID},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("DeleteContentTypeById", ctID, mock.Anything).Return(nil)
			},
			expectedCode: 0,
		},
		{
			name:  "Success - Delete By Slug (No ID provided)",
			input: &domain.ContentType{AccountId: accountID, Slug: "my-slug"},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("DeleteContentTypeBySlug", mock.Anything, mock.Anything).Return(nil)
			},
			expectedCode: 0,
		},
		{
			name:  "Fail - Content Type Not Found (By ID)",
			input: &domain.ContentType{AccountId: accountID, Id: ctID},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("DeleteContentTypeById", ctID, mock.Anything).Return(sql.ErrNoRows)
			},
			expectedCode:  http.StatusNotFound,
			expectedError: "does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAccount := new(mocks.MockAccountRepo)
			mockContent := new(mocks.MockContentTypeRepo)
			svc := services.NewContentTypeService(repository.Repositories{Account: mockAccount, ContentType: mockContent})

			tt.setupMock(mockAccount, mockContent)
			code, err := svc.DeleteContentType(context.Background(), tt.input)

			assert.Equal(t, tt.expectedCode, code)
			if tt.expectedError != "" {
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestContentTypeService_GetAllContentTypes(t *testing.T) {
	accountID := uuid.New()

	tests := []struct {
		name          string
		setupMock     func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo)
		expectedLen   int
		expectedCode  int
		expectedError bool
	}{
		{
			name: "Success - Return Multiple",
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("GetContentTypesByAccountId", accountID, mock.Anything).Return([]*domain.ContentType{{Id: uuid.New()}, {Id: uuid.New()}}, nil)
			},
			expectedLen:  2,
			expectedCode: 0,
		},
		{
			name: "Success - Empty Results (ErrNoRows)",
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("GetContentTypesByAccountId", accountID, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedLen:  0,
			expectedCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAccount := new(mocks.MockAccountRepo)
			mockContent := new(mocks.MockContentTypeRepo)
			svc := services.NewContentTypeService(repository.Repositories{Account: mockAccount, ContentType: mockContent})

			tt.setupMock(mockAccount, mockContent)
			res, code, err := svc.GetAllContentTypes(context.Background(), accountID)

			assert.Equal(t, tt.expectedCode, code)
			assert.Len(t, res, tt.expectedLen)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestContentTypeService_GetContentType_Formatting(t *testing.T) {
	accountID := uuid.New()
	rawSlug := accountID.String() + "my-blog-post"
	expectedCleanSlug := "my-blog-post"

	// Mock response from DB (simulating raw data)
	mockCT := &domain.ContentType{
		AccountId: accountID,
		Slug:      rawSlug, // The DB stores the full prefixed slug
	}

	mockAccount := new(mocks.MockAccountRepo)
	mockContent := new(mocks.MockContentTypeRepo)
	svc := services.NewContentTypeService(repository.Repositories{
		Account:     mockAccount,
		ContentType: mockContent,
	})

	// Setup: Account exists, and Slug search returns the raw object
	mockAccount.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
	mockContent.On("GetContentTypeBySlug", mock.Anything, mock.Anything).Return(mockCT, nil)

	// Act
	result, code, err := svc.GetContentType(context.Background(), &domain.ContentType{
		AccountId: accountID,
		Slug:      expectedCleanSlug,
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, code)
	// Here is the key: we verify the service returned the formatted slug, not the raw one
	assert.Equal(t, expectedCleanSlug, result.Slug, "The slug should be formatted (prefix removed) before returning")
}
