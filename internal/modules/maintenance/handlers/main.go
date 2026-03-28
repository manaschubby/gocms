package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	httpTransport "github.com/manaschubby/gocms/internal/transport/http"
	"github.com/manaschubby/gocms/internal/modules/maintenance/repository"
	"github.com/manaschubby/gocms/internal/modules/maintenance/services"
)

type Handlers struct {
	Category    CategoryHandlers
	Subcategory SubcategoryHandlers
	Detail      DetailHandlers
	Worker      WorkerHandlers
	Request     RequestHandlers
	Config      ConfigHandlers
}

func New(r repository.Repositories) *Handlers {
	s := services.New(r)
	return &Handlers{
		Category:    NewCategoryHandlers(s),
		Subcategory: NewSubcategoryHandlers(s),
		Detail:      NewDetailHandlers(s),
		Worker:      NewWorkerHandlers(s),
		Request:     NewRequestHandlers(s),
		Config:      NewConfigHandlers(s),
	}
}

func validationError(e echo.Context, msg string) error {
	return httpTransport.ErrWithMsg(e, http.StatusBadRequest, "failed to process request payload: "+msg, nil)
}