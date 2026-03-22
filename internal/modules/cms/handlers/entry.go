package handlers

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
	"github.com/manaschubby/gocms/internal/modules/cms/services"
	httpTransport "github.com/manaschubby/gocms/internal/transport/http"
)

type EntryHandlers interface {
	AddEntry(e echo.Context) error
	GetEntry(e echo.Context) error
}

type entryHandlers struct {
	cmsServices     services.Services
	cmsRepositories repository.Repositories
}

var _ EntryHandlers = &entryHandlers{}

func NewEntryHandlers(s services.Services, r repository.Repositories) EntryHandlers {
	return &entryHandlers{
		cmsServices:     s,
		cmsRepositories: r,
	}
}

type AddEntryInput struct {
	Title         string          `json:"title"`
	ContentData   json.RawMessage `json:"contentData"`
	Slug          string          `json:"slug"`
	Status        string          `json:"status"`
	ContentTypeId string          `json:"contentTypeId"`
	AccountId     string          `json:"accountId"`
}

func (h *entryHandlers) AddEntry(e echo.Context) error {
	var payload AddEntryInput

	err := e.Bind(&payload)
	if err != nil {
		return validationError(e, "invalid request payload type")
	}

	if payload.Title == "" || payload.ContentData == nil || payload.Status == "" || payload.Slug == "" || payload.ContentTypeId == "" || payload.AccountId == "" {
		return validationError(e, "title, contentData, slug, accountId, status and contentTypeId is required")
	}

	// Basic Validation (content type validation will be handled in service layer)
	ctId, err := uuid.Parse(payload.ContentTypeId)
	if err != nil {
		return validationError(e, "invalid contentTypeId, should match uuid: "+err.Error())
	}

	aid, err := uuid.Parse(payload.AccountId)
	if err != nil {
		return validationError(e, "invalid accountId, should match uuid: "+err.Error())
	}

	var status domain.EntryStatus
	err = status.Scan(payload.Status)
	if err != nil {
		return validationError(e, "invalid status, must be one of: "+strings.Join([]string{string(domain.StatusArchived), string(domain.StatusDraft), string(domain.StatusPublished)}, ", "))
	}

	var cd map[string]any
	err = json.Unmarshal(payload.ContentData, &cd)
	if err != nil {
		return validationError(e, "invalid contentData, it should be valid JSON")
	}

	entry := domain.Entry{
		Id:            uuid.New(),
		ContentTypeId: ctId,
		Slug:          payload.Slug,
		Title:         payload.Title,
		ContentData:   payload.ContentData,
		Status:        status,
		Version:       0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	code, err := h.cmsServices.Entry.CreateEntry(e.Request().Context(), &entry, aid)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}

	return httpTransport.Ok(e, entry)
}

type GetEntryInput struct {
	Id            string `query:"id"`
	Slug          string `query:"slug"`
	ContentTypeId string `query:"contentTypeId"`
}

func (h *entryHandlers) GetEntry(e echo.Context) error {
	var payload GetEntryInput

	err := e.Bind(&payload)
	if err != nil {
		return validationError(e, "invalid request payload type")
	}

	if payload.Id == "" && payload.Slug == "" && payload.ContentTypeId == "" {
		return validationError(e, "id, slug or contentTypeId is required")
	}

	var entry *domain.Entry
	var code int
	if payload.Id != "" {
		eid, err := uuid.Parse(payload.Id)
		if err != nil {

			return validationError(e, "invalid entry Id, should match UUID: ")
		}
		entry, code, err = h.cmsServices.Entry.GetEntry(e.Request().Context(), &domain.Entry{
			Id: eid,
		})
		if err != nil {
			return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
		}
		return httpTransport.Ok(e, entry)
	}

	ctId, err := uuid.Parse(payload.ContentTypeId)
	if err != nil {
		return validationError(e, "invalid contentTypeId, should match uuid: "+err.Error())
	}

	if payload.Slug != "" {
		entry, code, err = h.cmsServices.Entry.GetEntry(e.Request().Context(), &domain.Entry{
			Slug:          payload.Slug,
			ContentTypeId: ctId,
		})
		if err != nil {
			return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
		}
		return httpTransport.Ok(e, entry)
	}

	entries, code, err := h.cmsServices.Entry.GetAllEntries(e.Request().Context(), &ctId)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}

	return httpTransport.Ok(e, entries)
}
