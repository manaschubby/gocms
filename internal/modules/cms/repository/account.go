package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
)

type accountRepository struct {
	db *sqlx.DB
}

// Ensure type check
var _ AccountRepository = &accountRepository{}

func NewAccountRepository(db *sqlx.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

type GetAccountsOptions struct {
	IsActive *bool
	Context  *context.Context
}

func (r *accountRepository) GetAccounts(options GetAccountsOptions) ([]*domain.Account, error) {
	getAllAccounts := `SELECT * FROM account`
	getAccountsByActive := getAllAccounts + ` WHERE is_active = $1`

	ctx, cancel := ensureContext(options.Context)
	defer cancel()

	var accounts []*domain.Account
	var err error
	if options.IsActive != nil {
		err = r.db.SelectContext(ctx, &accounts, getAccountsByActive, &options.IsActive)
	} else {
		err = r.db.SelectContext(ctx, &accounts, getAllAccounts)
	}
	if err != nil {
		return nil, errors.New("Failed to Query DB:" + err.Error())
	}

	return accounts, nil
}

type CreateAccountOptions struct {
	Tx      *sqlx.Tx
	Context *context.Context
	// TODO: Implement this (Update the account if it already exists) (then make it public)
	update *bool
}

func (r *accountRepository) CreateAccount(account *domain.Account, options CreateAccountOptions) error {
	createAccount := `INSERT INTO account ("id", "name", "description", "is_active", "created_at", "updated_at") VALUES ($1,$2,$3,$4,$5,$6)`

	ctx, cancel := ensureContext(options.Context)
	defer cancel()

	execer := getExecerContextFromTxOrDB(options.Tx, r.db)

	_, err := execer.ExecContext(ctx, createAccount, account.Id.String(), account.Name, account.Description, account.IsActive, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		return errors.New("failed to save account in db: " + err.Error())
	}

	return nil
}

type GetAccountOptions struct {
	Context *context.Context
}

func (r *accountRepository) GetAccountByUUID(id uuid.UUID, options GetAccountOptions) (*domain.Account, error) {
	getAccount := `SELECT "id", "name", "description", "is_active", "created_at", "updated_at" FROM account WHERE id = $1`

	ctx, cancel := ensureContext(options.Context)
	defer cancel()

	acc := &domain.Account{}

	err := r.db.GetContext(ctx, acc, getAccount, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to fetch account from DB: " + err.Error())
	}

	return acc, nil
}
