package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/manaschubby/gocms/internal/modules/maintenance/domain"
	"github.com/lib/pq"
)

// ─── Category ────────────────────────────────────────────────────────────────

type categoryRepository struct{ db *sqlx.DB }

var _ CategoryRepository = &categoryRepository{}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(c *domain.MaintenanceCategory, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	_, err := getExecer(o.Tx, r.db).ExecContext(ctx,
		`INSERT INTO maintenance_categories (id, name, description, manager_email, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		c.Id, c.Name, c.Description, c.ManagerEmail, c.Active, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *categoryRepository) GetCategoryById(id uuid.UUID, o ReadOptions) (*domain.MaintenanceCategory, error) {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	var c domain.MaintenanceCategory
	err := r.db.GetContext(ctx, &c, `SELECT * FROM maintenance_categories WHERE id=$1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func (r *categoryRepository) GetAllCategories(o ReadOptions) ([]*domain.MaintenanceCategory, error) {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	var cs []*domain.MaintenanceCategory
	err := r.db.SelectContext(ctx, &cs, `SELECT * FROM maintenance_categories ORDER BY name`)
	return cs, err
}

func (r *categoryRepository) UpdateCategory(c *domain.MaintenanceCategory, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	_, err := getExecer(o.Tx, r.db).ExecContext(ctx,
		`UPDATE maintenance_categories SET name=$1, description=$2, manager_email=$3, active=$4, updated_at=$5 WHERE id=$6`,
		c.Name, c.Description, c.ManagerEmail, c.Active, time.Now(), c.Id,
	)
	return err
}

func (r *categoryRepository) DeleteCategory(id uuid.UUID, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	res, err := getExecer(o.Tx, r.db).ExecContext(ctx, `DELETE FROM maintenance_categories WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ─── Subcategory ─────────────────────────────────────────────────────────────

type subcategoryRepository struct{ db *sqlx.DB }

var _ SubcategoryRepository = &subcategoryRepository{}

func NewSubcategoryRepository(db *sqlx.DB) SubcategoryRepository {
	return &subcategoryRepository{db: db}
}

func (r *subcategoryRepository) CreateSubcategory(s *domain.MaintenanceSubcategory, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	_, err := getExecer(o.Tx, r.db).ExecContext(ctx,
		`INSERT INTO maintenance_subcategories (id, category_id, name, description, supervisor_email, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		s.Id, s.CategoryId, s.Name, s.Description, s.SupervisorEmail, s.Active, s.CreatedAt, s.UpdatedAt,
	)
	return err
}

func (r *subcategoryRepository) GetSubcategoryById(id uuid.UUID, o ReadOptions) (*domain.MaintenanceSubcategory, error) {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	var s domain.MaintenanceSubcategory
	err := r.db.GetContext(ctx, &s, `SELECT * FROM maintenance_subcategories WHERE id=$1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (r *subcategoryRepository) GetSubcategoriesByCategoryId(categoryId uuid.UUID, o ReadOptions) ([]*domain.MaintenanceSubcategory, error) {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	var ss []*domain.MaintenanceSubcategory
	err := r.db.SelectContext(ctx, &ss, `SELECT * FROM maintenance_subcategories WHERE category_id=$1 ORDER BY name`, categoryId)
	return ss, err
}

func (r *subcategoryRepository) UpdateSubcategory(s *domain.MaintenanceSubcategory, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	_, err := getExecer(o.Tx, r.db).ExecContext(ctx,
		`UPDATE maintenance_subcategories SET name=$1, description=$2, supervisor_email=$3, active=$4, updated_at=$5 WHERE id=$6`,
		s.Name, s.Description, s.SupervisorEmail, s.Active, time.Now(), s.Id,
	)
	return err
}

func (r *subcategoryRepository) DeleteSubcategory(id uuid.UUID, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	res, err := getExecer(o.Tx, r.db).ExecContext(ctx, `DELETE FROM maintenance_subcategories WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ─── Detail ──────────────────────────────────────────────────────────────────

type detailRepository struct{ db *sqlx.DB }

var _ DetailRepository = &detailRepository{}

func NewDetailRepository(db *sqlx.DB) DetailRepository {
	return &detailRepository{db: db}
}

func (r *detailRepository) CreateDetail(d *domain.MaintenanceDetail, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	_, err := getExecer(o.Tx, r.db).ExecContext(ctx,
		`INSERT INTO maintenance_details (id, subcategory_id, name, description, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		d.Id, d.SubcategoryId, d.Name, d.Description, d.Active, d.CreatedAt, d.UpdatedAt,
	)
	return err
}

func (r *detailRepository) GetDetailById(id uuid.UUID, o ReadOptions) (*domain.MaintenanceDetail, error) {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	var d domain.MaintenanceDetail
	err := r.db.GetContext(ctx, &d, `SELECT * FROM maintenance_details WHERE id=$1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &d, err
}

func (r *detailRepository) GetDetailsBySubcategoryId(subcategoryId uuid.UUID, o ReadOptions) ([]*domain.MaintenanceDetail, error) {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	var ds []*domain.MaintenanceDetail
	err := r.db.SelectContext(ctx, &ds, `SELECT * FROM maintenance_details WHERE subcategory_id=$1 ORDER BY name`, subcategoryId)
	return ds, err
}

func (r *detailRepository) UpdateDetail(d *domain.MaintenanceDetail, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	_, err := getExecer(o.Tx, r.db).ExecContext(ctx,
		`UPDATE maintenance_details SET name=$1, description=$2, active=$3, updated_at=$4 WHERE id=$5`,
		d.Name, d.Description, d.Active, time.Now(), d.Id,
	)
	return err
}

func (r *detailRepository) DeleteDetail(id uuid.UUID, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	res, err := getExecer(o.Tx, r.db).ExecContext(ctx, `DELETE FROM maintenance_details WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ─── Worker ──────────────────────────────────────────────────────────────────

type workerRepository struct{ db *sqlx.DB }

var _ WorkerRepository = &workerRepository{}

func NewWorkerRepository(db *sqlx.DB) WorkerRepository {
	return &workerRepository{db: db}
}

func (r *workerRepository) CreateWorker(w *domain.MaintenanceWorker, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	_, err := getExecer(o.Tx, r.db).ExecContext(ctx,
		`INSERT INTO maintenance_workers (id, name, user_email, phone, subcategory_id, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		w.Id, w.Name, w.UserEmail, w.Phone, w.SubcategoryId, w.Active, w.CreatedAt, w.UpdatedAt,
	)
	return err
}

func (r *workerRepository) GetWorkerById(id uuid.UUID, o ReadOptions) (*domain.MaintenanceWorker, error) {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	var w domain.MaintenanceWorker
	err := r.db.GetContext(ctx, &w, `SELECT * FROM maintenance_workers WHERE id=$1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &w, err
}

func (r *workerRepository) GetWorkersBySubcategoryId(subcategoryId uuid.UUID, o ReadOptions) ([]*domain.MaintenanceWorker, error) {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	var ws []*domain.MaintenanceWorker
	err := r.db.SelectContext(ctx, &ws, `SELECT * FROM maintenance_workers WHERE subcategory_id=$1 AND active=true ORDER BY name`, subcategoryId)
	return ws, err
}

func (r *workerRepository) UpdateWorker(w *domain.MaintenanceWorker, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	_, err := getExecer(o.Tx, r.db).ExecContext(ctx,
		`UPDATE maintenance_workers SET name=$1, user_email=$2, phone=$3, active=$4, updated_at=$5 WHERE id=$6`,
		w.Name, w.UserEmail, w.Phone, w.Active, time.Now(), w.Id,
	)
	return err
}

func (r *workerRepository) DeleteWorker(id uuid.UUID, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	res, err := getExecer(o.Tx, r.db).ExecContext(ctx, `DELETE FROM maintenance_workers WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ─── Request ─────────────────────────────────────────────────────────────────

type requestRepository struct{ db *sqlx.DB }

var _ RequestRepository = &requestRepository{}

func NewRequestRepository(db *sqlx.DB) RequestRepository {
	return &requestRepository{db: db}
}

func (r *requestRepository) CreateRequest(req *domain.MaintenanceRequest, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	_, err := getExecer(o.Tx, r.db).ExecContext(ctx,
		`INSERT INTO maintenance_requests
		 (id, requester_email, requester_name, location, category_id, subcategory_id, detail_id, description, status, escalation_level, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		req.Id, req.RequesterEmail, req.RequesterName, req.Location,
		req.CategoryId, req.SubcategoryId, req.DetailId, req.Description,
		req.Status, req.EscalationLevel, req.CreatedAt, req.UpdatedAt,
	)
	return err
}

func (r *requestRepository) GetRequestById(id uuid.UUID, o ReadOptions) (*domain.MaintenanceRequest, error) {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	var req domain.MaintenanceRequest
	err := r.db.GetContext(ctx, &req, `SELECT * FROM maintenance_requests WHERE id=$1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &req, err
}

func (r *requestRepository) GetRequests(filters RequestFilters, o ReadOptions) ([]*domain.MaintenanceRequest, error) {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()

	query := `SELECT * FROM maintenance_requests WHERE 1=1`
	args := []any{}
	i := 1

	if filters.RequesterEmail != nil {
		query += fmt.Sprintf(" AND requester_email=$%d", i)
		args = append(args, *filters.RequesterEmail)
		i++
	}
	if filters.CategoryId != nil {
		query += fmt.Sprintf(" AND category_id=$%d", i)
		args = append(args, *filters.CategoryId)
		i++
	}
	if filters.SubcategoryId != nil {
		query += fmt.Sprintf(" AND subcategory_id=$%d", i)
		args = append(args, *filters.SubcategoryId)
		i++
	}
	if filters.Status != nil {
		query += fmt.Sprintf(" AND status=$%d", i)
		args = append(args, *filters.Status)
		i++
	}
	if filters.EscalationLevel != nil {
		query += fmt.Sprintf(" AND escalation_level=$%d", i)
		args = append(args, *filters.EscalationLevel)
		i++
	}

	query += " ORDER BY created_at DESC"

	var reqs []*domain.MaintenanceRequest
	err := r.db.SelectContext(ctx, &reqs, query, args...)
	return reqs, err
}

func (r *requestRepository) UpdateRequest(req *domain.MaintenanceRequest, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	_, err := getExecer(o.Tx, r.db).ExecContext(ctx,
		`UPDATE maintenance_requests SET
		 status=$1, escalation_level=$2, last_escalated_at=$3,
		 assigned_worker_id=$4, assigned_at=$5,
		 resolved_at=$6, resolution_notes=$7, updated_at=$8
		 WHERE id=$9`,
		req.Status, req.EscalationLevel, req.LastEscalatedAt,
		req.AssignedWorkerId, req.AssignedAt,
		req.ResolvedAt, req.ResolutionNotes, time.Now(),
		req.Id,
	)
	return err
}

func (r *requestRepository) AddStatusLog(log *domain.MaintenanceStatusLog, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	_, err := getExecer(o.Tx, r.db).ExecContext(ctx,
		`INSERT INTO maintenance_status_log (id, request_id, user_email, action, previous_status, new_status, comments, timestamp)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		log.Id, log.RequestId, log.UserEmail, log.Action,
		log.PreviousStatus, log.NewStatus, log.Comments, log.Timestamp,
	)
	return err
}

func (r *requestRepository) AddEscalationLog(log *domain.MaintenanceEscalationLog, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	_, err := getExecer(o.Tx, r.db).ExecContext(ctx,
		`INSERT INTO maintenance_escalation_log (id, request_id, escalation_level, notified_emails, timestamp)
		 VALUES ($1,$2,$3,$4,$5)`,
		log.Id, log.RequestId, log.EscalationLevel, pq.Array(log.NotifiedEmails), log.Timestamp,
	)
	return err
}

func (r *requestRepository) GetPendingEscalations(o ReadOptions) ([]*domain.MaintenanceRequest, error) {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	var reqs []*domain.MaintenanceRequest
	err := r.db.SelectContext(ctx, &reqs,
		`SELECT * FROM maintenance_requests
		 WHERE status NOT IN ('Resolved','Rejected')
		 AND escalation_level != 'Level3'
		 ORDER BY created_at ASC`,
	)
	return reqs, err
}

// ─── Config ──────────────────────────────────────────────────────────────────

type configRepository struct{ db *sqlx.DB }

var _ ConfigRepository = &configRepository{}

func NewConfigRepository(db *sqlx.DB) ConfigRepository {
	return &configRepository{db: db}
}

func (r *configRepository) GetConfig(o ReadOptions) (*domain.MaintenanceConfig, error) {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	var c domain.MaintenanceConfig
	err := r.db.GetContext(ctx, &c, `SELECT dean_email, updated_at FROM maintenance_config WHERE id=1`)
	if err == sql.ErrNoRows {
		return &domain.MaintenanceConfig{}, nil
	}
	return &c, err
}

func (r *configRepository) UpsertConfig(c *domain.MaintenanceConfig, o WriteOptions) error {
	ctx, cancel := ensureContext(o.Context)
	defer cancel()
	_, err := getExecer(o.Tx, r.db).ExecContext(ctx,
		`INSERT INTO maintenance_config (id, dean_email, updated_at) VALUES (1,$1,$2)
		 ON CONFLICT (id) DO UPDATE SET dean_email=$1, updated_at=$2`,
		c.DeanEmail, time.Now(),
	)
	return err
}