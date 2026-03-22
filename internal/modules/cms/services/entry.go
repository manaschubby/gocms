package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/manaschubby/gocms/internal/modules/cms/common"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
)

type EntryService interface {
	CreateEntry(ctx context.Context, e *domain.Entry, accountId uuid.UUID) (code int, err error)
	GetEntry(ctx context.Context, e *domain.Entry) (entry *domain.Entry, code int, err error)
	GetAllEntries(ctx context.Context, ctId *uuid.UUID) (entries []*domain.Entry, code int, err error)
}

type entryService struct {
	r repository.Repositories
}

func NewEntryService(r repository.Repositories) EntryService {
	return &entryService{
		r: r,
	}
}

func (s *entryService) CreateEntry(ctx context.Context, e *domain.Entry, accountId uuid.UUID) (code int, err error) {
	// Validate accountId
	code, err = ValidateAccount(s.r.Account, ctx, accountId)
	if err != nil {
		return code, err
	}

	// Get Content Type
	ct, err := s.r.ContentType.GetContentTypeById(e.ContentTypeId, repository.GetContentTypeOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to fetch content type: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("failed to fetch content type")
	}
	if ct == nil {
		return http.StatusBadRequest, errors.New("content_type does not exist")
	}

	err = ct.Validate()
	if err != nil {
		return http.StatusBadRequest, errors.New("content_type schema is invalid")
	}

	// Check if pre-exists
	oldEntry, err := s.r.Entry.GetEntryByContentTypeAndSlug(e.ContentTypeId, e.Slug, repository.GetEntryOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to fetch entry from db: %v", err)
		return http.StatusInternalServerError, errors.New("failed to fetch entry from db")
	}
	if oldEntry != nil {
		return http.StatusBadRequest, errors.New("entry with this slug already exists in this content_type")
	}

	// Validate content data
	schema := ct.SchemaDefinition
	var contentData map[string]any
	err = common.JsonNumberDecode(e.ContentData, &contentData)
	if err != nil { // should not happen (callers should check for validation before passing)
		return http.StatusBadRequest, fmt.Errorf("failed to parse content data: %w", err)
	}

	errorColumns := make([]string, 0)
	for k, v := range schema {
		value := contentData[k]
		if v.DefaultValue != nil && value == nil {
			value = v.DefaultValue
		}
		if v.Required && value == nil {
			errorColumns = append(errorColumns, k+": required field and no default value exists")
			continue
		}

		err = v.ValidateAny(value)
		if err != nil {
			errorColumns = append(errorColumns, k+": "+err.Error())
		}
		contentData[k] = value
	}

	if len(errorColumns) != 0 {
		return http.StatusBadRequest, errors.New("schema validation failed for following columns: {" + strings.Join(errorColumns, ", ") + "}")
	}

	cd, err := json.Marshal(contentData)
	if err != nil {
		log.Printf("failed to prepare contentData for db insert: %v", err)
		return http.StatusInternalServerError, errors.New("failed to prepare contentData for db insert")
	}
	e.ContentData = json.RawMessage(cd)
	err = s.r.Entry.AddEntry(e, repository.AddEntryOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to insert entry into db: %v", err)
		return http.StatusInternalServerError, errors.New("failed to insert entry into db")
	}

	return 0, nil
}

func (s *entryService) GetEntry(ctx context.Context, e *domain.Entry) (entry *domain.Entry, code int, err error) {
	if e.Id != uuid.Nil {
		entry, err := s.r.Entry.GetEntryById(e.Id, repository.GetEntryOptions{Context: &ctx})
		if err != nil {
			log.Printf("failed to fetch entry from database: %v", err)
			return entry, http.StatusInternalServerError, errors.New("failed to fetch entry from database")
		}
		if entry == nil {
			return entry, http.StatusBadRequest, errors.New("entry does not exist")
		}

		return entry, code, err
	}

	contentType, err := s.r.ContentType.GetContentTypeById(e.ContentTypeId, repository.GetContentTypeOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to fetch contentType from database: %v", err)
		return entry, http.StatusInternalServerError, errors.New("failed to fetch contentType from database")
	}

	if contentType == nil {
		return entry, http.StatusBadRequest, errors.New("contentType does not exist")
	}

	entry, err = s.r.Entry.GetEntryByContentTypeAndSlug(e.ContentTypeId, e.Slug, repository.GetEntryOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to fetch entry from database: %v", err)
		return entry, http.StatusInternalServerError, errors.New("failed to fetch entry from database")
	}
	if entry == nil {
		return entry, http.StatusBadRequest, errors.New("entry does not exist")
	}
	return entry, code, err
}

func (s *entryService) GetAllEntries(ctx context.Context, ctId *uuid.UUID) (entries []*domain.Entry, code int, err error) {

	contentType, err := s.r.ContentType.GetContentTypeById(*ctId, repository.GetContentTypeOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to fetch contentType from database: %v", err)
		return entries, http.StatusInternalServerError, errors.New("failed to fetch contentType from database")
	}

	if contentType == nil {
		return entries, http.StatusBadRequest, errors.New("contentType does not exist")
	}

	entries, err = s.r.Entry.GetEntriesByContentType(*ctId, repository.GetEntryOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to fetch entries from database: %v", err)
		return entries, http.StatusInternalServerError, errors.New("failed to fetch contentType from database")
	}
	if entries == nil {
		return []*domain.Entry{}, 0, nil
	}

	return entries, code, err
}
