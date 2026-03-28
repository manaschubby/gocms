package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/manaschubby/gocms/internal/modules/maintenance/domain"
)

type CategoryRepository interface {
	CreateCategory(c *domain.MaintenanceCategory, options WriteOptions) error
	GetCategoryById(id uuid.UUID, options ReadOptions) (*domain.MaintenanceCategory, error)
	GetAllCategories(options ReadOptions) ([]*domain.MaintenanceCategory, error)
	UpdateCategory(c *domain.MaintenanceCategory, options WriteOptions) error
	DeleteCategory(id uuid.UUID, options WriteOptions) error
}

type SubcategoryRepository interface {
	CreateSubcategory(s *domain.MaintenanceSubcategory, options WriteOptions) error
	GetSubcategoryById(id uuid.UUID, options ReadOptions) (*domain.MaintenanceSubcategory, error)
	GetSubcategoriesByCategoryId(categoryId uuid.UUID, options ReadOptions) ([]*domain.MaintenanceSubcategory, error)
	UpdateSubcategory(s *domain.MaintenanceSubcategory, options WriteOptions) error
	DeleteSubcategory(id uuid.UUID, options WriteOptions) error
}

type DetailRepository interface {
	CreateDetail(d *domain.MaintenanceDetail, options WriteOptions) error
	GetDetailById(id uuid.UUID, options ReadOptions) (*domain.MaintenanceDetail, error)
	GetDetailsBySubcategoryId(subcategoryId uuid.UUID, options ReadOptions) ([]*domain.MaintenanceDetail, error)
	UpdateDetail(d *domain.MaintenanceDetail, options WriteOptions) error
	DeleteDetail(id uuid.UUID, options WriteOptions) error
}

type WorkerRepository interface {
	CreateWorker(w *domain.MaintenanceWorker, options WriteOptions) error
	GetWorkerById(id uuid.UUID, options ReadOptions) (*domain.MaintenanceWorker, error)
	GetWorkersBySubcategoryId(subcategoryId uuid.UUID, options ReadOptions) ([]*domain.MaintenanceWorker, error)
	UpdateWorker(w *domain.MaintenanceWorker, options WriteOptions) error
	DeleteWorker(id uuid.UUID, options WriteOptions) error
}

type RequestRepository interface {
	CreateRequest(r *domain.MaintenanceRequest, options WriteOptions) error
	GetRequestById(id uuid.UUID, options ReadOptions) (*domain.MaintenanceRequest, error)
	GetRequests(filters RequestFilters, options ReadOptions) ([]*domain.MaintenanceRequest, error)
	UpdateRequest(r *domain.MaintenanceRequest, options WriteOptions) error
	AddStatusLog(log *domain.MaintenanceStatusLog, options WriteOptions) error
	AddEscalationLog(log *domain.MaintenanceEscalationLog, options WriteOptions) error
	GetPendingEscalations(options ReadOptions) ([]*domain.MaintenanceRequest, error)
}

type ConfigRepository interface {
	GetConfig(options ReadOptions) (*domain.MaintenanceConfig, error)
	UpsertConfig(c *domain.MaintenanceConfig, options WriteOptions) error
}

type Repositories struct {
	Category    CategoryRepository
	Subcategory SubcategoryRepository
	Detail      DetailRepository
	Worker      WorkerRepository
	Request     RequestRepository
	Config      ConfigRepository
}

func Init(db *sqlx.DB) Repositories {
	return Repositories{
		Category:    NewCategoryRepository(db),
		Subcategory: NewSubcategoryRepository(db),
		Detail:      NewDetailRepository(db),
		Worker:      NewWorkerRepository(db),
		Request:     NewRequestRepository(db),
		Config:      NewConfigRepository(db),
	}
}

// Shared options
type WriteOptions struct {
	Tx      *sqlx.Tx
	Context *context.Context
}

type ReadOptions struct {
	Context *context.Context
}

type RequestFilters struct {
	RequesterEmail  *string
	CategoryId      *uuid.UUID
	SubcategoryId   *uuid.UUID
	Status          *domain.MaintenanceStatus
	EscalationLevel *domain.EscalationLevel
}

func getExecer(tx *sqlx.Tx, db *sqlx.DB) sqlx.ExecerContext {
	if tx != nil {
		return tx
	}
	return db
}

func ensureContext(pctx *context.Context) (context.Context, context.CancelFunc) {
	if pctx == nil {
		return context.WithCancel(context.Background())
	}
	return *pctx, func() {}
}