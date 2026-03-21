package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"log"

	"github.com/google/uuid"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
)

type ContentTypeService interface {
	GetAllContentTypes(ctx context.Context, accountId uuid.UUID) (contentTypes []*domain.ContentType, code int, error error)
	GetContentType(ctx context.Context, ct *domain.ContentType) (contentType *domain.ContentType, code int, error error)
	CreateNewContentType(ctx context.Context, ct *domain.ContentType) (code int, error error)
	DeleteContentType(ctx context.Context, ct *domain.ContentType) (code int, error error)
}

type contentTypeService struct {
	r repository.Repositories
}

func NewContentTypeService(r repository.Repositories) ContentTypeService {
	return &contentTypeService{
		r: r,
	}
}

func (s *contentTypeService) CreateNewContentType(ctx context.Context, ct *domain.ContentType) (code int, error error) {
	// Validate Schema Definitions
	length := len(ct.SchemaDefinition)
	if length == 0 {
		return http.StatusBadRequest, errors.New("at least one schemaDefinition Required")
	}

	for k, v := range ct.SchemaDefinition {
		if !v.IsValid() {
			return http.StatusBadRequest, errors.New("failed to validate schema definition for " + k)
		}
	}

	errCode, err := s.validateAccount(ctx, ct.AccountId)
	if err != nil {
		return errCode, err
	}

	// Validate Unique Slug
	slug := domain.GetContentTypeSlugFor(ct.AccountId, ct.Slug)

	encodedSlug, err := url.Parse(slug)
	if err != nil {
		return http.StatusBadRequest, errors.New("invalid characters in slug")
	}

	oldCt, err := s.r.ContentType.GetContentTypeBySlug(encodedSlug.String(), repository.GetContentTypeOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to check for current existing content types: %v", err)
		return http.StatusInternalServerError, errors.New("failed to check for existing content types")
	}
	if oldCt != nil {
		return http.StatusBadRequest, errors.New("content type with slug already exists")
	}

	// Create Content type
	ct.Id = uuid.New()

	ct.Slug = encodedSlug.String()
	err = s.r.ContentType.CreateNewContentType(ct, repository.CreateNewContentTypeOptions{
		Context: &ctx,
	})
	if err != nil {
		log.Printf("failed to create new content type: %v", err)
		return http.StatusInternalServerError, errors.New("failed to create new content type")
	}

	return 0, nil
}

func (s *contentTypeService) DeleteContentType(ctx context.Context, ct *domain.ContentType) (code int, error error) {
	errCode, err := s.validateAccount(ctx, ct.AccountId)
	if err != nil {
		return errCode, err
	}
	var emptyUUID uuid.UUID
	// Delete by ID
	if ct.Id.String() != emptyUUID.String() && ct.Id.String() != "" {
		err := s.r.ContentType.DeleteContentTypeById(ct.Id, repository.DeleteContentTypeOptions{Context: &ctx})
		if err != nil {
			if err == sql.ErrNoRows {
				return http.StatusNotFound, fmt.Errorf("content_type with id %s does not exist", ct.Id.String())
			}
			log.Printf("failed to delete content type from DB: %v", err)
			return http.StatusInternalServerError, errors.New("failed to delete content type from DB")
		}
		return 0, nil
	}

	// Else delete by slug
	slug := domain.GetContentTypeSlugFor(ct.AccountId, ct.Slug)
	err = s.r.ContentType.DeleteContentTypeBySlug(slug, repository.DeleteContentTypeOptions{Context: &ctx})
	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("content_type with slug %s does not exist", ct.Slug)
		}
		log.Printf("failed to delete content type from DB: %v", err)
		return http.StatusInternalServerError, errors.New("failed to delete content type from DB")
	}

	return 0, nil
}

func (s *contentTypeService) GetAllContentTypes(ctx context.Context, accountId uuid.UUID) (contentTypes []*domain.ContentType, code int, error error) {
	errCode, err := s.validateAccount(ctx, accountId)
	if err != nil {
		return nil, errCode, err
	}
	defer func() {
		for _, contentType := range contentTypes {
			if contentType != nil {
				contentType.Format(accountId)
			}
		}
	}()

	contentTypes, err = s.r.ContentType.GetContentTypesByAccountId(accountId, repository.GetContentTypeOptions{Context: &ctx})
	if err != nil {
		if err == sql.ErrNoRows {
			return []*domain.ContentType{}, 0, nil
		}
		log.Printf("failed to fetch content types from db: %v", err)
		return []*domain.ContentType{}, http.StatusInternalServerError, err
	}
	return contentTypes, 0, nil
}

func (s *contentTypeService) GetContentType(ctx context.Context, ct *domain.ContentType) (contentType *domain.ContentType, code int, error error) {
	errCode, err := s.validateAccount(ctx, ct.AccountId)
	if err != nil {
		return contentType, errCode, err
	}

	defer func() {
		if contentType != nil {
			contentType.Format(ct.AccountId)
		}
	}()

	var emptyUUID uuid.UUID
	// Get by ID
	if ct.Id.String() != emptyUUID.String() && ct.Id.String() != "" {
		contentType, err := s.r.ContentType.GetContentTypeById(ct.Id, repository.GetContentTypeOptions{Context: &ctx})
		if err != nil {
			if err == sql.ErrNoRows {
				return contentType, http.StatusNotFound, fmt.Errorf("content_type with id %s does not exist", ct.Id.String())
			}
			log.Printf("failed to get content type from DB: %v", err)
			return contentType, http.StatusInternalServerError, errors.New("failed to get content type from DB")
		}
		return contentType, 0, nil
	}

	// Else get by slug
	slug := domain.GetContentTypeSlugFor(ct.AccountId, ct.Slug)
	contentType, err = s.r.ContentType.GetContentTypeBySlug(slug, repository.GetContentTypeOptions{Context: &ctx})
	if err != nil {
		if err == sql.ErrNoRows {
			return contentType, http.StatusNotFound, fmt.Errorf("content_type with slug %s does not exist", ct.Slug)
		}
		log.Printf("failed to get content type from DB: %v", err)
		return contentType, http.StatusInternalServerError, errors.New("failed to get content type from DB")
	}

	return contentType, 0, nil
}

func (s *contentTypeService) validateAccount(ctx context.Context, aid uuid.UUID) (int, error) {
	account, err := s.r.Account.GetAccountByUUID(aid, repository.GetAccountOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to retrieve account data: %v", err)
		return http.StatusInternalServerError, errors.New("failed to retrieve account data")
	}

	if account == nil {
		return http.StatusBadRequest, errors.New("account not found")
	}
	return 0, nil
}
