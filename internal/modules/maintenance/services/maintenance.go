package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/manaschubby/gocms/internal/modules/maintenance/domain"
	"github.com/manaschubby/gocms/internal/modules/maintenance/repository"
)

// ─── Category ────────────────────────────────────────────────────────────────

type CategoryService interface {
	GetAllCategories(ctx context.Context) ([]*domain.MaintenanceCategory, int, error)
	GetCategory(ctx context.Context, id uuid.UUID) (*domain.MaintenanceCategory, int, error)
	CreateCategory(ctx context.Context, c *domain.MaintenanceCategory) (int, error)
	UpdateCategory(ctx context.Context, c *domain.MaintenanceCategory) (int, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) (int, error)
}

type categoryService struct{ r repository.Repositories }

func NewCategoryService(r repository.Repositories) CategoryService {
	return &categoryService{r: r}
}

func (s *categoryService) GetAllCategories(ctx context.Context) ([]*domain.MaintenanceCategory, int, error) {
	cs, err := s.r.Category.GetAllCategories(repository.ReadOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to get categories: %v", err)
		return nil, http.StatusInternalServerError, errors.New("failed to get categories")
	}
	return cs, 0, nil
}

func (s *categoryService) GetCategory(ctx context.Context, id uuid.UUID) (*domain.MaintenanceCategory, int, error) {
	c, err := s.r.Category.GetCategoryById(id, repository.ReadOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to get category: %v", err)
		return nil, http.StatusInternalServerError, errors.New("failed to get category")
	}
	if c == nil {
		return nil, http.StatusNotFound, fmt.Errorf("category %s not found", id)
	}
	return c, 0, nil
}

func (s *categoryService) CreateCategory(ctx context.Context, c *domain.MaintenanceCategory) (int, error) {
	if c.Name == "" || c.ManagerEmail == "" {
		return http.StatusBadRequest, errors.New("name and managerEmail are required")
	}
	c.Id = uuid.New()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	if err := s.r.Category.CreateCategory(c, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to create category: %v", err)
		return http.StatusInternalServerError, errors.New("failed to create category")
	}
	return 0, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, c *domain.MaintenanceCategory) (int, error) {
	if c.Id == uuid.Nil {
		return http.StatusBadRequest, errors.New("id is required")
	}
	if err := s.r.Category.UpdateCategory(c, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to update category: %v", err)
		return http.StatusInternalServerError, errors.New("failed to update category")
	}
	return 0, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uuid.UUID) (int, error) {
	err := s.r.Category.DeleteCategory(id, repository.WriteOptions{Context: &ctx})
	if err == sql.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("category %s not found", id)
	}
	if err != nil {
		log.Printf("failed to delete category: %v", err)
		return http.StatusInternalServerError, errors.New("failed to delete category")
	}
	return 0, nil
}

// ─── Subcategory ─────────────────────────────────────────────────────────────

type SubcategoryService interface {
	GetSubcategories(ctx context.Context, categoryId uuid.UUID) ([]*domain.MaintenanceSubcategory, int, error)
	GetSubcategory(ctx context.Context, id uuid.UUID) (*domain.MaintenanceSubcategory, int, error)
	CreateSubcategory(ctx context.Context, s *domain.MaintenanceSubcategory) (int, error)
	UpdateSubcategory(ctx context.Context, s *domain.MaintenanceSubcategory) (int, error)
	DeleteSubcategory(ctx context.Context, id uuid.UUID) (int, error)
}

type subcategoryService struct{ r repository.Repositories }

func NewSubcategoryService(r repository.Repositories) SubcategoryService {
	return &subcategoryService{r: r}
}

func (s *subcategoryService) GetSubcategories(ctx context.Context, categoryId uuid.UUID) ([]*domain.MaintenanceSubcategory, int, error) {
	ss, err := s.r.Subcategory.GetSubcategoriesByCategoryId(categoryId, repository.ReadOptions{Context: &ctx})
	if err != nil {
		log.Printf("failed to get subcategories: %v", err)
		return nil, http.StatusInternalServerError, errors.New("failed to get subcategories")
	}
	return ss, 0, nil
}

func (s *subcategoryService) GetSubcategory(ctx context.Context, id uuid.UUID) (*domain.MaintenanceSubcategory, int, error) {
	sub, err := s.r.Subcategory.GetSubcategoryById(id, repository.ReadOptions{Context: &ctx})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to get subcategory")
	}
	if sub == nil {
		return nil, http.StatusNotFound, fmt.Errorf("subcategory %s not found", id)
	}
	return sub, 0, nil
}

func (s *subcategoryService) CreateSubcategory(ctx context.Context, sub *domain.MaintenanceSubcategory) (int, error) {
	if sub.Name == "" || sub.SupervisorEmail == "" || sub.CategoryId == uuid.Nil {
		return http.StatusBadRequest, errors.New("name, supervisorEmail and categoryId are required")
	}
	sub.Id = uuid.New()
	sub.CreatedAt = time.Now()
	sub.UpdatedAt = time.Now()
	if err := s.r.Subcategory.CreateSubcategory(sub, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to create subcategory: %v", err)
		return http.StatusInternalServerError, errors.New("failed to create subcategory")
	}
	return 0, nil
}

func (s *subcategoryService) UpdateSubcategory(ctx context.Context, sub *domain.MaintenanceSubcategory) (int, error) {
	if sub.Id == uuid.Nil {
		return http.StatusBadRequest, errors.New("id is required")
	}
	if err := s.r.Subcategory.UpdateSubcategory(sub, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to update subcategory: %v", err)
		return http.StatusInternalServerError, errors.New("failed to update subcategory")
	}
	return 0, nil
}

func (s *subcategoryService) DeleteSubcategory(ctx context.Context, id uuid.UUID) (int, error) {
	err := s.r.Subcategory.DeleteSubcategory(id, repository.WriteOptions{Context: &ctx})
	if err == sql.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("subcategory %s not found", id)
	}
	if err != nil {
		log.Printf("failed to delete subcategory: %v", err)
		return http.StatusInternalServerError, errors.New("failed to delete subcategory")
	}
	return 0, nil
}

// ─── Detail ──────────────────────────────────────────────────────────────────

type DetailService interface {
	GetDetails(ctx context.Context, subcategoryId uuid.UUID) ([]*domain.MaintenanceDetail, int, error)
	GetDetail(ctx context.Context, id uuid.UUID) (*domain.MaintenanceDetail, int, error)
	CreateDetail(ctx context.Context, d *domain.MaintenanceDetail) (int, error)
	UpdateDetail(ctx context.Context, d *domain.MaintenanceDetail) (int, error)
	DeleteDetail(ctx context.Context, id uuid.UUID) (int, error)
}

type detailService struct{ r repository.Repositories }

func NewDetailService(r repository.Repositories) DetailService {
	return &detailService{r: r}
}

func (s *detailService) GetDetails(ctx context.Context, subcategoryId uuid.UUID) ([]*domain.MaintenanceDetail, int, error) {
	ds, err := s.r.Detail.GetDetailsBySubcategoryId(subcategoryId, repository.ReadOptions{Context: &ctx})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to get details")
	}
	return ds, 0, nil
}

func (s *detailService) GetDetail(ctx context.Context, id uuid.UUID) (*domain.MaintenanceDetail, int, error) {
	d, err := s.r.Detail.GetDetailById(id, repository.ReadOptions{Context: &ctx})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to get detail")
	}
	if d == nil {
		return nil, http.StatusNotFound, fmt.Errorf("detail %s not found", id)
	}
	return d, 0, nil
}

func (s *detailService) CreateDetail(ctx context.Context, d *domain.MaintenanceDetail) (int, error) {
	if d.Name == "" || d.SubcategoryId == uuid.Nil {
		return http.StatusBadRequest, errors.New("name and subcategoryId are required")
	}
	d.Id = uuid.New()
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	if err := s.r.Detail.CreateDetail(d, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to create detail: %v", err)
		return http.StatusInternalServerError, errors.New("failed to create detail")
	}
	return 0, nil
}

func (s *detailService) UpdateDetail(ctx context.Context, d *domain.MaintenanceDetail) (int, error) {
	if d.Id == uuid.Nil {
		return http.StatusBadRequest, errors.New("id is required")
	}
	if err := s.r.Detail.UpdateDetail(d, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to update detail: %v", err)
		return http.StatusInternalServerError, errors.New("failed to update detail")
	}
	return 0, nil
}

func (s *detailService) DeleteDetail(ctx context.Context, id uuid.UUID) (int, error) {
	err := s.r.Detail.DeleteDetail(id, repository.WriteOptions{Context: &ctx})
	if err == sql.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("detail %s not found", id)
	}
	if err != nil {
		log.Printf("failed to delete detail: %v", err)
		return http.StatusInternalServerError, errors.New("failed to delete detail")
	}
	return 0, nil
}

// ─── Worker ──────────────────────────────────────────────────────────────────

type WorkerService interface {
	GetWorkers(ctx context.Context, subcategoryId uuid.UUID) ([]*domain.MaintenanceWorker, int, error)
	GetWorker(ctx context.Context, id uuid.UUID) (*domain.MaintenanceWorker, int, error)
	CreateWorker(ctx context.Context, w *domain.MaintenanceWorker) (int, error)
	UpdateWorker(ctx context.Context, w *domain.MaintenanceWorker) (int, error)
	DeleteWorker(ctx context.Context, id uuid.UUID) (int, error)
}

type workerService struct{ r repository.Repositories }

func NewWorkerService(r repository.Repositories) WorkerService {
	return &workerService{r: r}
}

func (s *workerService) GetWorkers(ctx context.Context, subcategoryId uuid.UUID) ([]*domain.MaintenanceWorker, int, error) {
	ws, err := s.r.Worker.GetWorkersBySubcategoryId(subcategoryId, repository.ReadOptions{Context: &ctx})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to get workers")
	}
	return ws, 0, nil
}

func (s *workerService) GetWorker(ctx context.Context, id uuid.UUID) (*domain.MaintenanceWorker, int, error) {
	w, err := s.r.Worker.GetWorkerById(id, repository.ReadOptions{Context: &ctx})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to get worker")
	}
	if w == nil {
		return nil, http.StatusNotFound, fmt.Errorf("worker %s not found", id)
	}
	return w, 0, nil
}

func (s *workerService) CreateWorker(ctx context.Context, w *domain.MaintenanceWorker) (int, error) {
	if w.Name == "" || w.Phone == "" || w.SubcategoryId == uuid.Nil {
		return http.StatusBadRequest, errors.New("name, phone and subcategoryId are required")
	}
	w.Id = uuid.New()
	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()
	if err := s.r.Worker.CreateWorker(w, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to create worker: %v", err)
		return http.StatusInternalServerError, errors.New("failed to create worker")
	}
	return 0, nil
}

func (s *workerService) UpdateWorker(ctx context.Context, w *domain.MaintenanceWorker) (int, error) {
	if w.Id == uuid.Nil {
		return http.StatusBadRequest, errors.New("id is required")
	}
	if err := s.r.Worker.UpdateWorker(w, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to update worker: %v", err)
		return http.StatusInternalServerError, errors.New("failed to update worker")
	}
	return 0, nil
}

func (s *workerService) DeleteWorker(ctx context.Context, id uuid.UUID) (int, error) {
	err := s.r.Worker.DeleteWorker(id, repository.WriteOptions{Context: &ctx})
	if err == sql.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("worker %s not found", id)
	}
	if err != nil {
		log.Printf("failed to delete worker: %v", err)
		return http.StatusInternalServerError, errors.New("failed to delete worker")
	}
	return 0, nil
}

// ─── Request ─────────────────────────────────────────────────────────────────

type RequestService interface {
	GetRequests(ctx context.Context, filters repository.RequestFilters) ([]*domain.MaintenanceRequest, int, error)
	GetRequest(ctx context.Context, id uuid.UUID) (*domain.MaintenanceRequest, int, error)
	CreateRequest(ctx context.Context, r *domain.MaintenanceRequest) (int, error)
	UpdateRequest(ctx context.Context, r *domain.MaintenanceRequest) (int, error)
	AddStatusLog(ctx context.Context, l *domain.MaintenanceStatusLog) (int, error)
	AddEscalationLog(ctx context.Context, l *domain.MaintenanceEscalationLog) (int, error)
	GetPendingEscalations(ctx context.Context) ([]*domain.MaintenanceRequest, int, error)
	AssignWorker(ctx context.Context, requestId uuid.UUID, workerId uuid.UUID, userEmail string) (int, error)
	Resolve(ctx context.Context, requestId uuid.UUID, notes string, userEmail string) (int, error)
	Reject(ctx context.Context, requestId uuid.UUID, comments string, userEmail string) (int, error)
}

type requestService struct{ r repository.Repositories }

func NewRequestService(r repository.Repositories) RequestService {
	return &requestService{r: r}
}

func (s *requestService) GetRequests(ctx context.Context, filters repository.RequestFilters) ([]*domain.MaintenanceRequest, int, error) {
	reqs, err := s.r.Request.GetRequests(filters, repository.ReadOptions{Context: &ctx})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to get requests")
	}
	return reqs, 0, nil
}

func (s *requestService) GetRequest(ctx context.Context, id uuid.UUID) (*domain.MaintenanceRequest, int, error) {
	req, err := s.r.Request.GetRequestById(id, repository.ReadOptions{Context: &ctx})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to get request")
	}
	if req == nil {
		return nil, http.StatusNotFound, fmt.Errorf("request %s not found", id)
	}
	return req, 0, nil
}

func (s *requestService) CreateRequest(ctx context.Context, req *domain.MaintenanceRequest) (int, error) {
	if req.RequesterEmail == "" || req.Location == "" || req.Description == "" ||
		req.CategoryId == uuid.Nil || req.SubcategoryId == uuid.Nil {
		return http.StatusBadRequest, errors.New("requesterEmail, location, description, categoryId and subcategoryId are required")
	}
	req.Id = uuid.New()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	if err := s.r.Request.CreateRequest(req, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to create request: %v", err)
		return http.StatusInternalServerError, errors.New("failed to create request")
	}
	return 0, nil
}

func (s *requestService) UpdateRequest(ctx context.Context, req *domain.MaintenanceRequest) (int, error) {
	if req.Id == uuid.Nil {
		return http.StatusBadRequest, errors.New("id is required")
	}
	if err := s.r.Request.UpdateRequest(req, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to update request: %v", err)
		return http.StatusInternalServerError, errors.New("failed to update request")
	}
	return 0, nil
}

func (s *requestService) AddStatusLog(ctx context.Context, l *domain.MaintenanceStatusLog) (int, error) {
	if l.RequestId == uuid.Nil || l.UserEmail == "" || l.Action == "" {
		return http.StatusBadRequest, errors.New("requestId, userEmail and action are required")
	}
	l.Id = uuid.New()
	l.Timestamp = time.Now()
	if err := s.r.Request.AddStatusLog(l, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to add status log: %v", err)
		return http.StatusInternalServerError, errors.New("failed to add status log")
	}
	return 0, nil
}

func (s *requestService) AddEscalationLog(ctx context.Context, l *domain.MaintenanceEscalationLog) (int, error) {
	if l.RequestId == uuid.Nil {
		return http.StatusBadRequest, errors.New("requestId is required")
	}
	l.Id = uuid.New()
	l.Timestamp = time.Now()
	if err := s.r.Request.AddEscalationLog(l, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to add escalation log: %v", err)
		return http.StatusInternalServerError, errors.New("failed to add escalation log")
	}
	return 0, nil
}

func (s *requestService) GetPendingEscalations(ctx context.Context) ([]*domain.MaintenanceRequest, int, error) {
	reqs, err := s.r.Request.GetPendingEscalations(repository.ReadOptions{Context: &ctx})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to get pending escalations")
	}
	return reqs, 0, nil
}

func (s *requestService) AssignWorker(ctx context.Context, requestId uuid.UUID, workerId uuid.UUID, userEmail string) (int, error) {
	req, err := s.r.Request.GetRequestById(requestId, repository.ReadOptions{Context: &ctx})
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to get request")
	}
	if req == nil {
		return http.StatusNotFound, fmt.Errorf("request %s not found", requestId)
	}
	now := time.Now()
	prev := req.Status
	req.AssignedWorkerId = &workerId
	req.AssignedAt = &now
	req.Status = domain.StatusAssigned
	req.UpdatedAt = now
	if err := s.r.Request.UpdateRequest(req, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to assign worker: %v", err)
		return http.StatusInternalServerError, errors.New("failed to assign worker")
	}
	newStatus := domain.StatusAssigned
	logEntry := &domain.MaintenanceStatusLog{
		RequestId:      requestId,
		UserEmail:      userEmail,
		Action:         "assign_worker",
		PreviousStatus: &prev,
		NewStatus:      &newStatus,
	}
	if code, err := s.AddStatusLog(ctx, logEntry); err != nil {
		return code, err
	}
	return 0, nil
}

func (s *requestService) Resolve(ctx context.Context, requestId uuid.UUID, notes string, userEmail string) (int, error) {
	req, err := s.r.Request.GetRequestById(requestId, repository.ReadOptions{Context: &ctx})
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to get request")
	}
	if req == nil {
		return http.StatusNotFound, fmt.Errorf("request %s not found", requestId)
	}
	now := time.Now()
	prev := req.Status
	req.ResolvedAt = &now
	req.ResolutionNotes = notes
	req.Status = domain.StatusResolved
	req.UpdatedAt = now
	if err := s.r.Request.UpdateRequest(req, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to resolve request: %v", err)
		return http.StatusInternalServerError, errors.New("failed to resolve request")
	}
	newStatus := domain.StatusResolved
	logEntry := &domain.MaintenanceStatusLog{
		RequestId:      requestId,
		UserEmail:      userEmail,
		Action:         "resolve",
		PreviousStatus: &prev,
		NewStatus:      &newStatus,
		Comments:       notes,
	}
	if code, err := s.AddStatusLog(ctx, logEntry); err != nil {
		return code, err
	}
	return 0, nil
}

func (s *requestService) Reject(ctx context.Context, requestId uuid.UUID, comments string, userEmail string) (int, error) {
	req, err := s.r.Request.GetRequestById(requestId, repository.ReadOptions{Context: &ctx})
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to get request")
	}
	if req == nil {
		return http.StatusNotFound, fmt.Errorf("request %s not found", requestId)
	}
	prev := req.Status
	req.Status = domain.StatusRejected
	req.UpdatedAt = time.Now()
	if err := s.r.Request.UpdateRequest(req, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to reject request: %v", err)
		return http.StatusInternalServerError, errors.New("failed to reject request")
	}
	newStatus := domain.StatusRejected
	logEntry := &domain.MaintenanceStatusLog{
		RequestId:      requestId,
		UserEmail:      userEmail,
		Action:         "reject",
		PreviousStatus: &prev,
		NewStatus:      &newStatus,
		Comments:       comments,
	}
	if code, err := s.AddStatusLog(ctx, logEntry); err != nil {
		return code, err
	}
	return 0, nil
}

// ─── Config ──────────────────────────────────────────────────────────────────

type ConfigService interface {
	GetConfig(ctx context.Context) (*domain.MaintenanceConfig, int, error)
	UpdateConfig(ctx context.Context, c *domain.MaintenanceConfig) (int, error)
}

type configService struct{ r repository.Repositories }

func NewConfigService(r repository.Repositories) ConfigService {
	return &configService{r: r}
}

func (s *configService) GetConfig(ctx context.Context) (*domain.MaintenanceConfig, int, error) {
	c, err := s.r.Config.GetConfig(repository.ReadOptions{Context: &ctx})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to get config")
	}
	return c, 0, nil
}

func (s *configService) UpdateConfig(ctx context.Context, c *domain.MaintenanceConfig) (int, error) {
	if err := s.r.Config.UpsertConfig(c, repository.WriteOptions{Context: &ctx}); err != nil {
		log.Printf("failed to update config: %v", err)
		return http.StatusInternalServerError, errors.New("failed to update config")
	}
	return 0, nil
}
