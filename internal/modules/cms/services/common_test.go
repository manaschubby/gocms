package services_test

import (
	"context"
	"encoding/json"
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

func TestValidateContentType(t *testing.T) {
	ctId := uuid.New()
	ctx := context.Background()

	tests := []struct {
		name          string
		setupMock     func(m *mocks.MockContentTypeRepo)
		expectedCode  int
		expectedError string
	}{
		{
			name: "Success - ContentType Found",
			setupMock: func(m *mocks.MockContentTypeRepo) {
				m.On("GetContentTypeById", ctId, mock.MatchedBy(func(opt repository.GetContentTypeOptions) bool {
					return opt.Context != nil
				})).Return(&domain.ContentType{Id: ctId}, nil)
			},
			expectedCode:  0,
			expectedError: "",
		},
		{
			name: "Fail - Repository Error",
			setupMock: func(m *mocks.MockContentTypeRepo) {
				m.On("GetContentTypeById", ctId, mock.Anything).
					Return(nil, errors.New("connection reset"))
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "failed to fetch contentType from database",
		},
		{
			name: "Fail - ContentType Not Found (Nil)",
			setupMock: func(m *mocks.MockContentTypeRepo) {
				m.On("GetContentTypeById", ctId, mock.Anything).Return(nil, nil)
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "contentType does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(mocks.MockContentTypeRepo)
			tt.setupMock(mockRepo)

			// Act
			res, code, err := services.ValidateContentType(mockRepo, ctx, ctId)

			// Assert
			assert.Equal(t, tt.expectedCode, code)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, ctId, res.Id)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestValidateContentData(t *testing.T) {
	// Setup a baseline schema for testing validation logic
	schema := domain.SchemaDefinitionMap{
		"title": domain.SchemaDefinition{
			ColumnType: domain.ShortTextColumn,
			Required:   true,
		},
		"status": domain.SchemaDefinition{
			ColumnType:   domain.ShortTextColumn,
			DefaultValue: "draft",
		},
	}

	tests := []struct {
		name          string
		inputJSON     json.RawMessage
		schema        domain.SchemaDefinitionMap
		expectedError string
	}{
		{
			name:      "Success - Valid Data All Fields",
			inputJSON: json.RawMessage(`{"title": "Hello World", "status": "published"}`),
			schema:    schema,
		},
		{
			name:      "Success - Missing field but has Default Value",
			inputJSON: json.RawMessage(`{"title": "Hello World"}`), // 'status' missing, will fall back to default
			schema:    schema,
		},
		{
			name:          "Fail - Missing Required Field",
			inputJSON:     json.RawMessage(`{"status": "published"}`), // 'title' is missing
			schema:        schema,
			expectedError: "schema validation failed for following columns: {title: required field and no default value exists}",
		},
		{
			name:      "Success - Optional Field Omitted Completely",
			inputJSON: json.RawMessage(`{"title": "Hello"}`),
			schema: domain.SchemaDefinitionMap{
				"title":          domain.SchemaDefinition{Required: true, ColumnType: domain.ShortTextColumn},
				"optional_field": domain.SchemaDefinition{Required: false, ColumnType: domain.ShortTextColumn}, // No default, not required
			},
			expectedError: "",
		},
		{
			name:          "Success - Explicit Null triggers Default Value",
			inputJSON:     json.RawMessage(`{"title": "Hello", "status": null}`),
			schema:        schema,
			expectedError: "",
		},
		{
			name:      "Success - Explicit Null for Optional Field is Skipped",
			inputJSON: json.RawMessage(`{"title": "Hello", "optional_field": null}`),
			schema: domain.SchemaDefinitionMap{
				"title":          domain.SchemaDefinition{Required: true, ColumnType: domain.ShortTextColumn},
				"optional_field": domain.SchemaDefinition{Required: false, ColumnType: domain.ShortTextColumn},
			},
			expectedError: "",
		},
		{
			name:      "Fail - Explicit Null for Required Field",
			inputJSON: json.RawMessage(`{"title": null}`),
			schema: domain.SchemaDefinitionMap{
				"title": domain.SchemaDefinition{Required: true, ColumnType: domain.ShortTextColumn},
			},
			expectedError: "schema validation failed for following columns: {title: required field and no default value exists}",
		},
		{
			name:          "Fail - Extra Columns in Content Data",
			inputJSON:     json.RawMessage(`{"title": "Hello", "status": "draft", "extra_unwanted_field": "oops"}`),
			schema:        schema,
			expectedError: "extra columns found in contentData",
		},
		{
			name:          "Fail - Extra Columns with Default Triggered",
			inputJSON:     json.RawMessage(`{"title": "Hello", "extra_field": "oops"}`), // 'status' triggers default, but extra field still fails
			schema:        schema,
			expectedError: "extra columns found in contentData",
		},
		{
			name:          "Fail - Invalid JSON Parsing",
			inputJSON:     json.RawMessage(`{ bad json format `),
			schema:        schema,
			expectedError: "failed to parse content data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			res, err := services.ValidateContentData(tt.inputJSON, tt.schema)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)

				// For cases where default was applied, we decode back and verify
				if tt.name == "Success - Missing field but has Default Value" {
					var parsed map[string]any
					err = json.Unmarshal(res, &parsed)
					assert.NoError(t, err)
					assert.Equal(t, "draft", parsed["status"]) // verifies the default injection worked
				}
			}
		})
	}
}
