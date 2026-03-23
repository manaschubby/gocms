package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
)

type AccountRepository interface {
	GetAccounts(options GetAccountsOptions) ([]*domain.Account, error)
	CreateAccount(account *domain.Account, options CreateAccountOptions) error
	GetAccountByUUID(id uuid.UUID, options GetAccountOptions) (*domain.Account, error)
}

type ContentTypeRepository interface {
	CreateNewContentType(ct *domain.ContentType, options CreateNewContentTypeOptions) error
	GetContentTypeBySlug(slug string, options GetContentTypeOptions) (*domain.ContentType, error)
	GetContentTypeById(id uuid.UUID, options GetContentTypeOptions) (*domain.ContentType, error)
	GetContentTypesByAccountId(account uuid.UUID, options GetContentTypeOptions) ([]*domain.ContentType, error)
	DeleteContentTypeBySlug(slug string, options DeleteContentTypeOptions) error
	DeleteContentTypeById(id uuid.UUID, options DeleteContentTypeOptions) error
}

type EntryRepository interface {
	AddEntry(e *domain.Entry, options AddEntryOptions) error
	GetEntryByContentTypeAndSlug(ctId uuid.UUID, slug string, options GetEntryOptions) (e *domain.Entry, err error)
	GetEntryById(eid uuid.UUID, options GetEntryOptions) (e *domain.Entry, err error)
	GetEntriesByContentType(ctId uuid.UUID, options GetEntryOptions) (e []*domain.Entry, err error)
	GetEntriesByFilter(entry *domain.Entry, options GetEntryOptions) (e []*domain.Entry, err error)
	UpdateEntry(entry *domain.Entry, options UpdateEntryOptions) (err error)
}

type Repositories struct {
	Account     AccountRepository
	ContentType ContentTypeRepository
	Entry       EntryRepository
}

func Init(db *sqlx.DB) Repositories {
	r := Repositories{
		Account:     NewAccountRepository(db),
		ContentType: NewContentTypeRepository(db),
		Entry:       NewEntryRepository(db),
	}
	return r
}

func getExecerContextFromTxOrDB(tx *sqlx.Tx, db *sqlx.DB) sqlx.ExecerContext {
	var execer sqlx.ExecerContext
	if tx != nil {
		execer = tx
	} else {
		execer = db
	}
	return execer
}

func ensureContext(pctx *context.Context) (context.Context, context.CancelFunc) {
	var ctx context.Context
	var cancel context.CancelFunc
	if pctx == nil {
		ctx, cancel = context.WithCancel(context.Background())
		return ctx, cancel
	}
	ctx = *pctx
	return ctx, func() {}
}
