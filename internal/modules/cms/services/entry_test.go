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
	"github.com/manaschubby/gocms/internal/modules/cms/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEntryService_CreateEntry(t *testing.T) {
	accountID := uuid.New()
	ctID := uuid.New()
	ctx := context.Background()

	// Setup a valid Content Type with a schema
	validCT := &domain.ContentType{
		Id:        ctID,
		AccountId: accountID,
		SchemaDefinition: map[string]domain.SchemaDefinition{
			"title": {
				ColumnType:       domain.ShortTextColumn,
				ColumnDefinition: domain.SingleValuedColumn,
				Required:         true,
			},
			"status": {
				ColumnType:       domain.ShortTextColumn,
				ColumnDefinition: domain.SingleValuedColumn,
				DefaultValue:     "draft",
			},
		},
	}

	tests := []struct {
		name          string
		inputEntry    *domain.Entry
		setupMock     func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo, em *mocks.MockEntryRepo)
		expectedCode  int
		expectedError string
	}{
		{
			name:       "Fail - Content Type Not Found",
			inputEntry: &domain.Entry{ContentTypeId: ctID, Slug: "test"},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo, em *mocks.MockEntryRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("GetContentTypeById", ctID, mock.Anything).Return(nil, nil)
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "contentType does not exist",
		},
		{
			name:       "Fail - Duplicate Slug",
			inputEntry: &domain.Entry{ContentTypeId: ctID, Slug: "existing-slug"},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo, em *mocks.MockEntryRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("GetContentTypeById", ctID, mock.Anything).Return(validCT, nil)
				em.On("GetEntryByContentTypeAndSlug", ctID, "existing-slug", mock.Anything).Return(&domain.Entry{}, nil)
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "already exists",
		},
		{
			name: "Fail - Schema Validation (Missing Required)",
			inputEntry: &domain.Entry{
				ContentTypeId: ctID,
				Slug:          "new-post",
				ContentData:   json.RawMessage(`{"status": "published"}`), // title missing
			},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo, em *mocks.MockEntryRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("GetContentTypeById", ctID, mock.Anything).Return(validCT, nil)
				em.On("GetEntryByContentTypeAndSlug", ctID, "new-post", mock.Anything).Return(nil, nil)
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "required field",
		},
		{
			name: "Success - With Default Values",
			inputEntry: &domain.Entry{
				ContentTypeId: ctID,
				Slug:          "valid-post",
				ContentData:   json.RawMessage(`{"title": "Hello World"}`), // status missing, should get 'draft'
			},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo, em *mocks.MockEntryRepo) {
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("GetContentTypeById", ctID, mock.Anything).Return(validCT, nil)
				em.On("GetEntryByContentTypeAndSlug", ctID, "valid-post", mock.Anything).Return(nil, nil)
				em.On("AddEntry", mock.MatchedBy(func(e *domain.Entry) bool {
					var data map[string]any
					json.Unmarshal(e.ContentData, &data)
					return data["status"] == "draft" // Verify default was applied
				}), mock.Anything).Return(nil)
			},
			expectedCode: 0,
		},
		{
			name: "Success - With Numerical Value",
			inputEntry: &domain.Entry{
				ContentTypeId: ctID,
				Slug:          "valid-post",
				ContentData:   json.RawMessage(`{"title": "Hello World", "age": 25}`),
			},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo, em *mocks.MockEntryRepo) {
				schema := map[string]domain.SchemaDefinition{
					"title": {ColumnType: domain.ShortTextColumn, Required: true},
					"age":   {ColumnType: domain.NumberColumn, Required: true},
				}
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("GetContentTypeById", ctID, mock.Anything).Return(&domain.ContentType{Id: ctID, SchemaDefinition: schema}, nil)
				em.On("GetEntryByContentTypeAndSlug", ctID, "valid-post", mock.Anything).Return(nil, nil)
				em.On("AddEntry", mock.MatchedBy(func(e *domain.Entry) bool {
					var data map[string]any
					json.Unmarshal(e.ContentData, &data)
					return data["age"] == float64(25)
				}), mock.Anything).Return(nil)
			},
			expectedCode: 0,
		},
		{
			name: "Fail - Schema Validation (Type Mismatch: Expected Number, got String)",
			inputEntry: &domain.Entry{
				ContentTypeId: ctID,
				Slug:          "bad-type-post",
				// 'age' is defined as NumberColumn in the setupMock below
				ContentData: json.RawMessage(`{"title": "Valid Title", "age": "twenty-five"}`),
			},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo, em *mocks.MockEntryRepo) {
				schema := map[string]domain.SchemaDefinition{
					"title": {ColumnType: domain.ShortTextColumn, Required: true},
					"age":   {ColumnType: domain.NumberColumn, Required: true},
				}
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("GetContentTypeById", ctID, mock.Anything).Return(&domain.ContentType{Id: ctID, SchemaDefinition: schema}, nil)
				em.On("GetEntryByContentTypeAndSlug", ctID, "bad-type-post", mock.Anything).Return(nil, nil)
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "schema validation failed for following columns: {age: needs to be a number}",
		},
		{
			name: "Fail - Schema Validation (Constraint: ShortText too long)",
			inputEntry: &domain.Entry{
				ContentTypeId: ctID,
				Slug:          "long-text-post",
				ContentData:   json.RawMessage(`{"title": "` + EXTREMELY_LONG_STRING + `"}`),
			},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo, em *mocks.MockEntryRepo) {
				schema := map[string]domain.SchemaDefinition{
					"title": {ColumnType: domain.ShortTextColumn, Required: true},
				}
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("GetContentTypeById", ctID, mock.Anything).Return(&domain.ContentType{Id: ctID, SchemaDefinition: schema}, nil)
				em.On("GetEntryByContentTypeAndSlug", ctID, "long-text-post", mock.Anything).Return(nil, nil)
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "should be less than", // Matches domain.SHORT_TEXT_MAX_LENGTH error
		},
		{
			name: "Fail - Schema Validation (List Constraint: Expected Array, got Single Value)",
			inputEntry: &domain.Entry{
				ContentTypeId: ctID,
				Slug:          "not-an-array",
				ContentData:   json.RawMessage(`{"tags": "golang"}`),
			},
			setupMock: func(am *mocks.MockAccountRepo, cm *mocks.MockContentTypeRepo, em *mocks.MockEntryRepo) {
				schema := map[string]domain.SchemaDefinition{
					"tags": {
						ColumnType:       domain.ShortTextColumn,
						ColumnDefinition: domain.ListValuedColumn, // Explicitly a list
					},
				}
				am.On("GetAccountByUUID", accountID, mock.Anything).Return(&domain.Account{Id: accountID}, nil)
				cm.On("GetContentTypeById", ctID, mock.Anything).Return(&domain.ContentType{Id: ctID, SchemaDefinition: schema}, nil)
				em.On("GetEntryByContentTypeAndSlug", ctID, "not-an-array", mock.Anything).Return(nil, nil)
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "needs to be an array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ma, mc, me := new(mocks.MockAccountRepo), new(mocks.MockContentTypeRepo), new(mocks.MockEntryRepo)
			repos := repository.Repositories{Account: ma, ContentType: mc, Entry: me}
			svc := services.NewEntryService(repos)

			tt.setupMock(ma, mc, me)
			code, err := svc.CreateEntry(ctx, tt.inputEntry, accountID)

			assert.Equal(t, tt.expectedCode, code)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestEntryService_GetEntry(t *testing.T) {
	ctID := uuid.New()
	entryID := uuid.New()
	ctx := context.Background()

	tests := []struct {
		name         string
		input        *domain.Entry
		setupMock    func(mc *mocks.MockContentTypeRepo, me *mocks.MockEntryRepo)
		expectedCode int
	}{
		{
			name:  "Success - Get By ID",
			input: &domain.Entry{Id: entryID},
			setupMock: func(mc *mocks.MockContentTypeRepo, me *mocks.MockEntryRepo) {
				me.On("GetEntryById", entryID, mock.Anything).Return(&domain.Entry{Id: entryID}, nil)
			},
			expectedCode: 0,
		},
		{
			name:  "Success - Get By Slug",
			input: &domain.Entry{ContentTypeId: ctID, Slug: "my-slug"},
			setupMock: func(mc *mocks.MockContentTypeRepo, me *mocks.MockEntryRepo) {
				mc.On("GetContentTypeById", ctID, mock.Anything).Return(&domain.ContentType{Id: ctID}, nil)
				me.On("GetEntryByContentTypeAndSlug", ctID, "my-slug", mock.Anything).Return(&domain.Entry{Slug: "my-slug"}, nil)
			},
			expectedCode: 0,
		},
		{
			name:  "Fail - Entry Not Found By ID",
			input: &domain.Entry{Id: entryID},
			setupMock: func(mc *mocks.MockContentTypeRepo, me *mocks.MockEntryRepo) {
				me.On("GetEntryById", entryID, mock.Anything).Return(nil, nil)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "Fail - DB Error GetEntryById",
			input: &domain.Entry{Id: entryID},
			setupMock: func(mc *mocks.MockContentTypeRepo, me *mocks.MockEntryRepo) {
				me.On("GetEntryById", entryID, mock.Anything).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:  "Fail - DB Error GetContentTypeById",
			input: &domain.Entry{ContentTypeId: ctID, Slug: "my-slug"},
			setupMock: func(mc *mocks.MockContentTypeRepo, me *mocks.MockEntryRepo) {
				mc.On("GetContentTypeById", ctID, mock.Anything).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:  "Fail - ContentType Does Not Exist",
			input: &domain.Entry{ContentTypeId: ctID, Slug: "my-slug"},
			setupMock: func(mc *mocks.MockContentTypeRepo, me *mocks.MockEntryRepo) {
				mc.On("GetContentTypeById", ctID, mock.Anything).Return(nil, nil)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "Fail - DB Error GetEntryByContentTypeAndSlug",
			input: &domain.Entry{ContentTypeId: ctID, Slug: "my-slug"},
			setupMock: func(mc *mocks.MockContentTypeRepo, me *mocks.MockEntryRepo) {
				mc.On("GetContentTypeById", ctID, mock.Anything).Return(&domain.ContentType{Id: ctID}, nil)
				me.On("GetEntryByContentTypeAndSlug", ctID, "my-slug", mock.Anything).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:  "Fail - Entry Not Found By Slug",
			input: &domain.Entry{ContentTypeId: ctID, Slug: "my-slug"},
			setupMock: func(mc *mocks.MockContentTypeRepo, me *mocks.MockEntryRepo) {
				mc.On("GetContentTypeById", ctID, mock.Anything).Return(&domain.ContentType{Id: ctID}, nil)
				me.On("GetEntryByContentTypeAndSlug", ctID, "my-slug", mock.Anything).Return(nil, nil)
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc, me := new(mocks.MockContentTypeRepo), new(mocks.MockEntryRepo)
			svc := services.NewEntryService(repository.Repositories{ContentType: mc, Entry: me})
			tt.setupMock(mc, me)

			_, code, _ := svc.GetEntry(ctx, tt.input)
			assert.Equal(t, tt.expectedCode, code)
		})
	}
}

func TestEntryService_GetAllEntries(t *testing.T) {
	ctID := uuid.New()
	ctx := context.Background()

	tests := []struct {
		name         string
		input        *domain.Entry
		setupMock    func(mc *mocks.MockContentTypeRepo, me *mocks.MockEntryRepo)
		expectedCode int
		expectedLen  int
		expectedErr  bool
	}{
		{
			name:  "Success - Found Entries",
			input: &domain.Entry{ContentTypeId: ctID},
			setupMock: func(mc *mocks.MockContentTypeRepo, me *mocks.MockEntryRepo) {
				mc.On("GetContentTypeById", ctID, mock.Anything).Return(&domain.ContentType{Id: ctID}, nil)
				me.On("GetEntriesByFilter", mock.Anything, mock.Anything).Return([]*domain.Entry{{Id: uuid.New()}, {Id: uuid.New()}}, nil)
			},
			expectedCode: 0,
			expectedLen:  2,
			expectedErr:  false,
		},
		{
			name:  "Success - No Entries Found (Empty Slice)",
			input: &domain.Entry{ContentTypeId: ctID},
			setupMock: func(mc *mocks.MockContentTypeRepo, me *mocks.MockEntryRepo) {
				mc.On("GetContentTypeById", ctID, mock.Anything).Return(&domain.ContentType{Id: ctID}, nil)
				me.On("GetEntriesByFilter", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedCode: 0,
			expectedLen:  0,
			expectedErr:  false,
		},
		{
			name:  "Fail - ContentType Not Found",
			input: &domain.Entry{ContentTypeId: ctID},
			setupMock: func(mc *mocks.MockContentTypeRepo, me *mocks.MockEntryRepo) {
				mc.On("GetContentTypeById", ctID, mock.Anything).Return(nil, nil)
			},
			expectedCode: http.StatusBadRequest,
			expectedLen:  0,
			expectedErr:  true,
		},
		{
			name:  "Fail - DB Error on GetEntriesByFilter",
			input: &domain.Entry{ContentTypeId: ctID},
			setupMock: func(mc *mocks.MockContentTypeRepo, me *mocks.MockEntryRepo) {
				mc.On("GetContentTypeById", ctID, mock.Anything).Return(&domain.ContentType{Id: ctID}, nil)
				me.On("GetEntriesByFilter", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedLen:  0,
			expectedErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc, me := new(mocks.MockContentTypeRepo), new(mocks.MockEntryRepo)
			svc := services.NewEntryService(repository.Repositories{ContentType: mc, Entry: me})
			tt.setupMock(mc, me)

			res, code, err := svc.GetAllEntries(ctx, tt.input)

			assert.Equal(t, tt.expectedCode, code)
			assert.Len(t, res, tt.expectedLen)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
