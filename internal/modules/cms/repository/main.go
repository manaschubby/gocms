package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
)

type AccountRepository interface {
	GetAccounts(options GetAccountsOptions) ([]*domain.Account, error)
	CreateAccount(account *domain.Account, options CreateAccountOptions) error
}

type Repositories struct {
	Account AccountRepository
}

func Init(db *sqlx.DB) Repositories {
	accountRepo := NewAccountRepository(db)
	r := Repositories{
		Account: accountRepo,
	}
	return r
}
