package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
)

type EntryService interface {
	CreateEntry(ctx context.Context, e *domain.Entry, accountId uuid.UUID) (code int, err error)
	GetEntry(ctx context.Context, e *domain.Entry) (entry *domain.Entry, code int, err error)
	GetAllEntries(ctx context.Context, e *domain.Entry) (entries []*domain.Entry, code int, err error)
	UpdateEntry(ctx context.Context, e *domain.Entry) (entry *domain.Entry, code int, err error)
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
	ct, code, err := ValidateContentType(s.r.ContentType, ctx, e.ContentTypeId)
	if err != nil {
		return code, err
	}

	err = ct.Validate()
	if err != nil {
		return http.StatusBadRequest, errors.New("contentType schema is invalid")
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
	cd, err := ValidateContentData(e.ContentData, ct.SchemaDefinition)
	if err != nil {
		return http.StatusBadRequest, err
	}

	e.ContentData = cd
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

	_, code, err = ValidateContentType(s.r.ContentType, ctx, e.ContentTypeId)
	if err != nil {
		return entry, code, err
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

func (s *entryService) GetAllEntries(ctx context.Context, e *domain.Entry) (entries []*domain.Entry, code int, err error) {
	_, code, err = ValidateContentType(s.r.ContentType, ctx, e.ContentTypeId)
	if err != nil {
		return entries, code, err
	}

	entries, err = s.r.Entry.GetEntriesByFilter(e, repository.GetEntryOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to fetch entries from database: %v", err)
		return entries, http.StatusInternalServerError, errors.New("failed to fetch entries from database")
	}
	if entries == nil {
		return []*domain.Entry{}, 0, nil
	}

	return entries, code, err
}

func (s *entryService) UpdateEntry(ctx context.Context, e *domain.Entry) (entry *domain.Entry, code int, err error) {
	entry, err = s.r.Entry.GetEntryById(e.Id, repository.GetEntryOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed tp fetch entry from Db: %v", err)
		return entry, http.StatusInternalServerError, errors.New("failed to fetch entry from db")
	}
	if entry == nil {
		return entry, http.StatusBadRequest, errors.New("entry does not exist")
	}

	ct, code, err := ValidateContentType(s.r.ContentType, ctx, entry.ContentTypeId)
	if err != nil {
		return entry, code, err
	}

	if e.ContentData != nil {
		contentData, err := ValidateContentData(e.ContentData, ct.SchemaDefinition)
		if err != nil {
			return entry, http.StatusBadRequest, err
		}

		entry.ContentData = contentData
	}

	if e.Status != "" {
		err := entry.Status.Scan(e.Status)
		if err != nil {
			return entry, http.StatusBadRequest, fmt.Errorf("invalid status found: %w", err)
		}
	}

	if e.Title != "" {
		entry.Title = e.Title
	}

	if entry.IsDifferentTo(*e) {

		entry.Version = entry.Version + 1
		entry.UpdatedAt = time.Now()

		err := s.r.Entry.UpdateEntry(entry, repository.UpdateEntryOptions{Context: &ctx})
		if err != nil {
			log.Printf("failed to update entry in DB: %v", err)
			return entry, http.StatusInternalServerError, errors.New("failed to update entry in DB")
		}
	}

	return entry, code, err
}
