package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
)

type entryRepositry struct {
	db *sqlx.DB
}

func NewEntryRepository(db *sqlx.DB) EntryRepository {
	return &entryRepositry{
		db: db,
	}
}

type AddEntryOptions struct {
	Context *context.Context
	Tx      *sqlx.Tx
}

func (r *entryRepositry) AddEntry(e *domain.Entry, options AddEntryOptions) error {
	addEntry := `INSERT INTO entries ("id", "content_type_id", "slug", "title", "content_data", "status", "version", "created_at", "updated_at") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	ctx, cancel := ensureContext(options.Context)
	defer cancel()

	execer := getExecerContextFromTxOrDB(options.Tx, r.db)

	result, err := execer.ExecContext(ctx, addEntry, e.Id, e.ContentTypeId, e.Slug, e.Title, e.ContentData, e.Status, e.Version, e.CreatedAt, e.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save entry to db: %w", err)
	}
	added, _ := result.RowsAffected()
	if added != 1 {
		return fmt.Errorf("no entry was saved to db")
	}

	return nil
}

type GetEntryOptions struct {
	Context *context.Context
}

func (r *entryRepositry) GetEntryByContentTypeAndSlug(ctId uuid.UUID, slug string, options GetEntryOptions) (e *domain.Entry, err error) {
	getEntryByContentTypeAndSlug := `SELECT "id", "content_type_id", "slug", "title", "content_data", "status", "version", "created_at", "updated_at" FROM entries WHERE content_type_id = $1 AND slug = $2`

	ctx, cancel := ensureContext(options.Context)
	defer cancel()

	e = &domain.Entry{}
	err = r.db.GetContext(ctx, e, getEntryByContentTypeAndSlug, ctId, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get entry from db: " + err.Error())
	}

	return e, err
}

func (r *entryRepositry) GetEntryById(eid uuid.UUID, options GetEntryOptions) (e *domain.Entry, err error) {
	getEntryById := `SELECT "id", "content_type_id", "slug", "title", "content_data", "status", "version", "created_at", "updated_at" FROM entries WHERE id = $1`

	ctx, cancel := ensureContext(options.Context)
	defer cancel()

	e = &domain.Entry{}
	err = r.db.GetContext(ctx, e, getEntryById, eid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get entry from db: " + err.Error())
	}

	return e, err
}

func (r *entryRepositry) GetEntriesByContentType(ctId uuid.UUID, options GetEntryOptions) (e []*domain.Entry, err error) {
	getEntriesByContentTypeId := `SELECT "id", "content_type_id", "slug", "title", "content_data", "status", "version", "created_at", "updated_at" FROM entries WHERE content_type_id = $1`

	ctx, cancel := ensureContext(options.Context)
	defer cancel()

	e = make([]*domain.Entry, 0)
	err = r.db.SelectContext(ctx, &e, getEntriesByContentTypeId, ctId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get entry from db: " + err.Error())
	}

	return e, err
}
