package services

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
)

func ValidateAccount(r repository.AccountRepository, ctx context.Context, aid uuid.UUID) (int, error) {
	account, err := r.GetAccountByUUID(aid, repository.GetAccountOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to retrieve account data: %v", err)
		return http.StatusInternalServerError, errors.New("failed to retrieve account data")
	}

	if account == nil {
		return http.StatusBadRequest, errors.New("account not found")
	}
	return 0, nil
}
