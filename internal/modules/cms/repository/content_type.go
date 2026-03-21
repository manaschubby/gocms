package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
)

type contentTypeRepository struct {
	db *sqlx.DB
}

// Ensure type check
var _ ContentTypeRepository = &contentTypeRepository{}

func NewContentTypeRepository(db *sqlx.DB) ContentTypeRepository {
	return &contentTypeRepository{
		db: db,
	}
}

type CreateNewContentTypeOptions struct {
	Tx      *sqlx.Tx
	Context *context.Context
}

func (r *contentTypeRepository) CreateNewContentType(ct *domain.ContentType, options CreateNewContentTypeOptions) error {
	createContentType := `INSERT INTO content_types ("id", "account_id", "name", "slug", "schema_definition", "description", "created_at", "updated_at") VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	ctx, cancel := ensureContext(options.Context)
	defer cancel()

	execer := getExecerContextFromTxOrDB(options.Tx, r.db)

	_, err := execer.ExecContext(ctx, createContentType, ct.Id, ct.AccountId, ct.Name, ct.Slug, ct.SchemaDefinition, ct.Description, ct.CreatedAt, ct.UpdatedAt)
	if err != nil {
		return errors.New("failed to save content-type in db: " + err.Error())
	}
	return nil
}

type GetContentTypeOptions struct {
	Context *context.Context
}

func (r *contentTypeRepository) GetContentTypeBySlug(slug string, options GetContentTypeOptions) (*domain.ContentType, error) {
	getContentTypeBySlug := `SELECT "id", "account_id", "name", "slug", "schema_definition", "description", "created_at", "updated_at" FROM content_types WHERE slug = $1`
	ctx, cancel := ensureContext(options.Context)
	defer cancel()

	var ct domain.ContentType

	err := r.db.GetContext(ctx, &ct, getContentTypeBySlug, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed to fetch content type from db: %w", err)
		}
	}

	return &ct, nil
}
