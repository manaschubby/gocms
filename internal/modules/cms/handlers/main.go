package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
	"github.com/manaschubby/gocms/internal/modules/cms/services"
	httpTransport "github.com/manaschubby/gocms/internal/transport/http"
)

type Handlers struct {
	Account     AccountHandlers
	ContentType ContentTypeHandlers
	Entry       EntryHandlers
}

func New(r repository.Repositories) *Handlers {
	cmsService := services.New(r)
	return &Handlers{
		Account:     NewAccountHandlers(r),
		ContentType: NewContentTypeHandlers(r, *cmsService),
		Entry:       NewEntryHandlers(*cmsService, r),
	}
}

func validationError(e echo.Context, msg string) error {
	return httpTransport.ErrWithMsg(e, http.StatusBadRequest, "failed to process request payload: "+msg, nil)
}
