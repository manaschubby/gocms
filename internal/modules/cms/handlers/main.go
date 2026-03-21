package handlers

import (
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
	"github.com/manaschubby/gocms/internal/modules/cms/services"
)

type Handlers struct {
	Account     AccountHandlers
	ContentType ContentTypeHandlers
}

func New(r repository.Repositories) *Handlers {
	cmsService := services.New(r)
	return &Handlers{
		Account:     NewAccountHandlers(r),
		ContentType: NewContentTypeHandlers(r, *cmsService),
	}
}
