package handlers

import (
	"net/http"
	"time"

	"log"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
	"github.com/manaschubby/gocms/internal/modules/cms/services"
	httpTransport "github.com/manaschubby/gocms/internal/transport/http"
)

type ContentTypeHandlers interface {
	CreateContentType(e echo.Context) error
	DeleteContentType(e echo.Context) error
	GetContentType(e echo.Context) error
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
		log.Printf("failed to bind input from request: %v", err)
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

type DeleteContentTypeInput struct {
	Id        *string `json:"id,omitempty"`
	Slug      *string `json:"slug,omitempty"`
	AccountId *string `json:"accountId"`
}

func (h *contentTypeHandlers) DeleteContentType(e echo.Context) error {
	var payload DeleteContentTypeInput

	err := e.Bind(&payload)
	if err != nil {
		return validationError(e, "failed to parse request payload")
	}

	if payload.Id == nil && payload.Slug == nil {
		return validationError(e, "atleast slug or Id is required")
	}

	aid, err := uuid.Parse(*payload.AccountId)
	if err != nil {
		return validationError(e, "invalid account id: "+err.Error())
	}

	var id uuid.UUID
	if payload.Id != nil {
		id, err = uuid.Parse(*payload.Id)
		if err != nil {
			return validationError(e, "invalid content_type id: "+err.Error())
		}
	}

	code, err := h.cmsServices.ContentType.DeleteContentType(e.Request().Context(), &domain.ContentType{
		Id:        id,
		Slug:      *payload.Slug,
		AccountId: aid,
	})
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}

	return httpTransport.Ok(e, nil)
}

type GetContentTypeInput struct {
	Id        *string `query:"id"`
	Slug      *string `query:"slug"`
	AccountId *string `query:"accountId"`
}

func (h *contentTypeHandlers) GetContentType(e echo.Context) error {
	var payload GetContentTypeInput

	err := e.Bind(&payload)
	if err != nil {
		return validationError(e, "failed to parse request payload")
	}

	aid, err := uuid.Parse(*payload.AccountId)
	if err != nil {
		return validationError(e, "invalid account id: "+err.Error())
	}

	if payload.Id == nil && payload.Slug == nil {
		contentTypes, code, err := h.cmsServices.ContentType.GetAllContentTypes(e.Request().Context(), aid)
		if err != nil {
			return httpTransport.ErrWithMsg(e, code, err.Error(), contentTypes)
		}
		return httpTransport.Ok(e, contentTypes)
	}

	var ctId uuid.UUID
	if payload.Id != nil {
		ctId, err = uuid.Parse(*payload.Id)
		if err != nil {
			return validationError(e, "invalid content_type id: "+err.Error())
		}
	}

	var slug string
	if payload.Slug != nil {
		slug = *payload.Slug
	}

	ct := &domain.ContentType{
		AccountId: aid,
		Id:        ctId,
		Slug:      slug,
	}
	contentTypes, code, err := h.cmsServices.ContentType.GetContentType(e.Request().Context(), ct)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), contentTypes)
	}
	return httpTransport.Ok(e, contentTypes)

}
