package handlers

import (
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
)

type Handlers struct {
	Account AccountHandler
}

func New(r repository.Repositories) *Handlers {
	return &Handlers{
		Account: NewAccountHandler(r),
	}
}
