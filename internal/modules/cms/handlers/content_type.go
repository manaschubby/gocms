package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
	"github.com/manaschubby/gocms/internal/modules/cms/services"
	httpTransport "github.com/manaschubby/gocms/internal/transport/http"
)

type ContentTypeHandlers interface {
	CreateContentType(e echo.Context) error
}

type contentTypeHandlers struct {
	cmsRepos    repository.Repositories
	cmsServices services.Services
}

var _ ContentTypeHandlers = &contentTypeHandlers{}

func NewContentTypeHandlers(r repository.Repositories, s services.Services) ContentTypeHandlers {
	return &contentTypeHandlers{
		cmsRepos:    r,
		cmsServices: s,
	}
}

type CreateContentTypeInput struct {
	AccountId        string                             `json:"accountId"`
	Name             string                             `json:"name"`
	Slug             string                             `json:"slug"`
	Description      string                             `json:"description,omitempty"`
	SchemaDefinition map[string]domain.SchemaDefinition `json:"schemaDefinition"`
}

func validationError(e echo.Context, msg string) error {
	return httpTransport.ErrWithMsg(e, http.StatusBadRequest, "failed to process request payload: "+msg, nil)
}

func (h *contentTypeHandlers) CreateContentType(e echo.Context) error {
	// UnMarshal Request
	var payload *CreateContentTypeInput
	err := e.Bind(&payload)
	if err != nil || payload == nil {
		log.Errorf("failed to bind input from request: %v", err)
		return validationError(e, err.Error())
	}

	// Validate Request
	aid, err := uuid.Parse(payload.AccountId)
	if err != nil {
		return validationError(e, "invalid account id: "+err.Error())
	}

	if payload.Name == "" || payload.Slug == "" || payload.Description == "" || payload.SchemaDefinition == nil {
		return validationError(e, "name, description, slug and schemaDefinition Required ")
	}

	// Validate Schema Definitions
	length := len(payload.SchemaDefinition)
	if length == 0 {
		return validationError(e, "at least one schemaDefinition Required ")
	}

	for k, v := range payload.SchemaDefinition {
		if !v.IsValid() {
			return validationError(e, "failed to validate schema definition for "+k)
		}
	}

	ctx := e.Request().Context()
	ct := &domain.ContentType{
		Name:             payload.Name,
		AccountId:        aid,
		Slug:             payload.Slug,
		Description:      payload.Description,
		SchemaDefinition: payload.SchemaDefinition,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	code, err := h.cmsServices.ContentType.CreateNewContentType(ctx, ct)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}

	return httpTransport.Ok(e, ct)
}
