package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
	httpTransport "github.com/manaschubby/gocms/internal/transport/http"
)

type AccountHandlers interface {
	GetAllAccounts(e echo.Context) error
}

type accountHandlers struct {
	cmsRepositories repository.Repositories
}

// Ensure interface compliance
var _ AccountHandlers = &accountHandlers{}

func NewAccountHandlers(r repository.Repositories) *accountHandlers {
	return &accountHandlers{
		cmsRepositories: r,
	}
}

func (h *accountHandlers) GetAllAccounts(e echo.Context) error {
	accounts, err := h.cmsRepositories.Account.GetAccounts(repository.GetAccountsOptions{})
	if err != nil {
		log.Println("failed to fetch accounts from DB: " + err.Error())
		return httpTransport.Err(e, http.StatusInternalServerError, nil)
	}

	return httpTransport.Ok(e, accounts)
}
