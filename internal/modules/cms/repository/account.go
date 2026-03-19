package repository

import (
	"context"
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

	var ctx context.Context
	var cancel context.CancelFunc
	if options.Context == nil {
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
	} else {
		ctx = *options.Context
	}

	var accountRows *sqlx.Rows
	var err error
	if options.IsActive != nil {
		accountRows, err = r.db.QueryxContext(ctx, getAccountsByActive, &options.IsActive)
	} else {
		accountRows, err = r.db.QueryxContext(ctx, getAllAccounts)
	}
	if err != nil {
		return nil, errors.New("Failed to Query DB:" + err.Error())
	}
	defer accountRows.Close()

	var accounts []*domain.Account

	for accountRows.Next() {
		account := &domain.Account{}
		var id *string
		accountRows.Scan(&id, &account.Name, &account.Description, &account.IsActive, &account.CreatedAt, &account.UpdatedAt)
		idUuid, err := uuid.Parse(*id)
		if err != nil {
			return nil, errors.New("Failed to parse UUID for account: " + err.Error())
		}
		account.Id = idUuid
		accounts = append(accounts, account)
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

	var ctx context.Context
	var cancel context.CancelFunc
	if options.Context == nil {
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
	} else {
		ctx = *options.Context
	}

	var execer sqlx.ExecerContext
	if options.Tx != nil {
		execer = options.Tx
	} else {
		execer = r.db
	}

	_, err := execer.ExecContext(ctx, createAccount, account.Id.String(), account.Name, account.Description, account.IsActive, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		return errors.New("failed to save account in db: " + err.Error())
	}

	return nil
}
