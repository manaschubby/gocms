package services

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
)

type ContentTypeService interface {
	CreateNewContentType(ctx context.Context, ct *domain.ContentType) (code int, error error)
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

	// Validate Account
	account, err := s.r.Account.GetAccountByUUID(ct.AccountId, repository.GetAccountOptions{Context: &ctx})
	if err != nil {
		log.Errorf("failed to retrieve account data: %v", err)
		return http.StatusInternalServerError, errors.New("failed to retrieve account data")
	}

	if account == nil {
		return http.StatusBadRequest, errors.New("account not found")
	}

	// Validate Unique Slug
	slug := account.Id.String() + ct.Slug

	encodedSlug, err := url.Parse(slug)
	if err != nil {
		return http.StatusBadRequest, errors.New("invalid characters in slug")
	}

	oldCt, err := s.r.ContentType.GetContentTypeBySlug(encodedSlug.String(), repository.GetContentTypeOptions{Context: &ctx})
	if err != nil {
		log.Errorf("failed to check for current existing content types: %v", err)
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
		log.Errorf("failed to create new content type: %v", err)
		return http.StatusInternalServerError, errors.New("failed to create new content type")
	}

	return http.StatusOK, nil
}
