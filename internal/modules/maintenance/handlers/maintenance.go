package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/manaschubby/gocms/internal/modules/maintenance/domain"
	"github.com/manaschubby/gocms/internal/modules/maintenance/repository"
	"github.com/manaschubby/gocms/internal/modules/maintenance/services"
	httpTransport "github.com/manaschubby/gocms/internal/transport/http"
)

// ─── Category ────────────────────────────────────────────────────────────────

type CategoryHandlers interface {
	GetCategories(e echo.Context) error
	GetCategory(e echo.Context) error
	CreateCategory(e echo.Context) error
	UpdateCategory(e echo.Context) error
	DeleteCategory(e echo.Context) error
}

type categoryHandlers struct{ s *services.Services }

var _ CategoryHandlers = &categoryHandlers{}

func NewCategoryHandlers(s *services.Services) CategoryHandlers {
	return &categoryHandlers{s: s}
}

func (h *categoryHandlers) GetCategories(e echo.Context) error {
	cs, code, err := h.s.Category.GetAllCategories(e.Request().Context())
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, cs)
}

func (h *categoryHandlers) GetCategory(e echo.Context) error {
	id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid category id")
	}
	c, code, err := h.s.Category.GetCategory(e.Request().Context(), id)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, c)
}

type createCategoryInput struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ManagerEmail string `json:"managerEmail"`
}

func (h *categoryHandlers) CreateCategory(e echo.Context) error {
	var payload createCategoryInput
	if err := e.Bind(&payload); err != nil {
		return validationError(e, err.Error())
	}
	c := &domain.MaintenanceCategory{
		Name:         payload.Name,
		Description:  payload.Description,
		ManagerEmail: payload.ManagerEmail,
	}
	code, err := h.s.Category.CreateCategory(e.Request().Context(), c)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, c)
}

type updateCategoryInput struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ManagerEmail string `json:"managerEmail"`
	Active       bool   `json:"active"`
}

func (h *categoryHandlers) UpdateCategory(e echo.Context) error {
	id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid category id")
	}
	var payload updateCategoryInput
	if err := e.Bind(&payload); err != nil {
		return validationError(e, err.Error())
	}
	c := &domain.MaintenanceCategory{
		Id:           id,
		Name:         payload.Name,
		Description:  payload.Description,
		ManagerEmail: payload.ManagerEmail,
		Active:       payload.Active,
	}
	code, err := h.s.Category.UpdateCategory(e.Request().Context(), c)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, c)
}

func (h *categoryHandlers) DeleteCategory(e echo.Context) error {
	id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid category id")
	}
	code, err := h.s.Category.DeleteCategory(e.Request().Context(), id)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, nil)
}

// ─── Subcategory ─────────────────────────────────────────────────────────────

type SubcategoryHandlers interface {
	GetSubcategories(e echo.Context) error
	GetSubcategory(e echo.Context) error
	CreateSubcategory(e echo.Context) error
	UpdateSubcategory(e echo.Context) error
	DeleteSubcategory(e echo.Context) error
}

type subcategoryHandlers struct{ s *services.Services }

var _ SubcategoryHandlers = &subcategoryHandlers{}

func NewSubcategoryHandlers(s *services.Services) SubcategoryHandlers {
	return &subcategoryHandlers{s: s}
}

func (h *subcategoryHandlers) GetSubcategories(e echo.Context) error {
	categoryId, err := uuid.Parse(e.QueryParam("categoryId"))
	if err != nil {
		return validationError(e, "invalid categoryId query param")
	}
	ss, code, err := h.s.Subcategory.GetSubcategories(e.Request().Context(), categoryId)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, ss)
}

func (h *subcategoryHandlers) GetSubcategory(e echo.Context) error {
	id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid subcategory id")
	}
	s, code, err := h.s.Subcategory.GetSubcategory(e.Request().Context(), id)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, s)
}

type createSubcategoryInput struct {
	CategoryId      string `json:"categoryId"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	SupervisorEmail string `json:"supervisorEmail"`
}

func (h *subcategoryHandlers) CreateSubcategory(e echo.Context) error {
	var payload createSubcategoryInput
	if err := e.Bind(&payload); err != nil {
		return validationError(e, err.Error())
	}
	categoryId, err := uuid.Parse(payload.CategoryId)
	if err != nil {
		return validationError(e, "invalid categoryId")
	}
	s := &domain.MaintenanceSubcategory{
		CategoryId:      categoryId,
		Name:            payload.Name,
		Description:     payload.Description,
		SupervisorEmail: payload.SupervisorEmail,
	}
	code, err := h.s.Subcategory.CreateSubcategory(e.Request().Context(), s)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, s)
}

type updateSubcategoryInput struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	SupervisorEmail string `json:"supervisorEmail"`
	Active          bool   `json:"active"`
}

func (h *subcategoryHandlers) UpdateSubcategory(e echo.Context) error {
	id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid subcategory id")
	}
	var payload updateSubcategoryInput
	if err := e.Bind(&payload); err != nil {
		return validationError(e, err.Error())
	}
	s := &domain.MaintenanceSubcategory{
		Id:              id,
		Name:            payload.Name,
		Description:     payload.Description,
		SupervisorEmail: payload.SupervisorEmail,
		Active:          payload.Active,
	}
	code, err := h.s.Subcategory.UpdateSubcategory(e.Request().Context(), s)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, s)
}

func (h *subcategoryHandlers) DeleteSubcategory(e echo.Context) error {
	id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid subcategory id")
	}
	code, err := h.s.Subcategory.DeleteSubcategory(e.Request().Context(), id)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, nil)
}

// ─── Detail ──────────────────────────────────────────────────────────────────

type DetailHandlers interface {
	GetDetails(e echo.Context) error
	GetDetail(e echo.Context) error
	CreateDetail(e echo.Context) error
	UpdateDetail(e echo.Context) error
	DeleteDetail(e echo.Context) error
}

type detailHandlers struct{ s *services.Services }

var _ DetailHandlers = &detailHandlers{}

func NewDetailHandlers(s *services.Services) DetailHandlers {
	return &detailHandlers{s: s}
}

func (h *detailHandlers) GetDetails(e echo.Context) error {
	subcategoryId, err := uuid.Parse(e.QueryParam("subcategoryId"))
	if err != nil {
		return validationError(e, "invalid subcategoryId query param")
	}
	ds, code, err := h.s.Detail.GetDetails(e.Request().Context(), subcategoryId)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, ds)
}

func (h *detailHandlers) GetDetail(e echo.Context) error {
	id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid detail id")
	}
	d, code, err := h.s.Detail.GetDetail(e.Request().Context(), id)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, d)
}

type createDetailInput struct {
	SubcategoryId string `json:"subcategoryId"`
	Name          string `json:"name"`
	Description   string `json:"description"`
}

func (h *detailHandlers) CreateDetail(e echo.Context) error {
	var payload createDetailInput
	if err := e.Bind(&payload); err != nil {
		return validationError(e, err.Error())
	}
	subcategoryId, err := uuid.Parse(payload.SubcategoryId)
	if err != nil {
		return validationError(e, "invalid subcategoryId")
	}
	d := &domain.MaintenanceDetail{
		SubcategoryId: subcategoryId,
		Name:          payload.Name,
		Description:   payload.Description,
	}
	code, err := h.s.Detail.CreateDetail(e.Request().Context(), d)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, d)
}

type updateDetailInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

func (h *detailHandlers) UpdateDetail(e echo.Context) error {
	id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid detail id")
	}
	var payload updateDetailInput
	if err := e.Bind(&payload); err != nil {
		return validationError(e, err.Error())
	}
	d := &domain.MaintenanceDetail{
		Id:          id,
		Name:        payload.Name,
		Description: payload.Description,
		Active:      payload.Active,
	}
	code, err := h.s.Detail.UpdateDetail(e.Request().Context(), d)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, d)
}

func (h *detailHandlers) DeleteDetail(e echo.Context) error {
	id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid detail id")
	}
	code, err := h.s.Detail.DeleteDetail(e.Request().Context(), id)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, nil)
}

// ─── Worker ──────────────────────────────────────────────────────────────────

type WorkerHandlers interface {
	GetWorkers(e echo.Context) error
	GetWorker(e echo.Context) error
	CreateWorker(e echo.Context) error
	UpdateWorker(e echo.Context) error
	DeleteWorker(e echo.Context) error
}

type workerHandlers struct{ s *services.Services }

var _ WorkerHandlers = &workerHandlers{}

func NewWorkerHandlers(s *services.Services) WorkerHandlers {
	return &workerHandlers{s: s}
}

func (h *workerHandlers) GetWorkers(e echo.Context) error {
	subcategoryId, err := uuid.Parse(e.QueryParam("subcategoryId"))
	if err != nil {
		return validationError(e, "invalid subcategoryId query param")
	}
	ws, code, err := h.s.Worker.GetWorkers(e.Request().Context(), subcategoryId)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, ws)
}

func (h *workerHandlers) GetWorker(e echo.Context) error {
	id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid worker id")
	}
	w, code, err := h.s.Worker.GetWorker(e.Request().Context(), id)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, w)
}

type createWorkerInput struct {
	SubcategoryId string `json:"subcategoryId"`
	Name          string `json:"name"`
	UserEmail     string `json:"userEmail"`
	Phone         string `json:"phone"`
}

func (h *workerHandlers) CreateWorker(e echo.Context) error {
	var payload createWorkerInput
	if err := e.Bind(&payload); err != nil {
		return validationError(e, err.Error())
	}
	subcategoryId, err := uuid.Parse(payload.SubcategoryId)
	if err != nil {
		return validationError(e, "invalid subcategoryId")
	}
	w := &domain.MaintenanceWorker{
		SubcategoryId: subcategoryId,
		Name:          payload.Name,
		UserEmail:     payload.UserEmail,
		Phone:         payload.Phone,
	}
	code, err := h.s.Worker.CreateWorker(e.Request().Context(), w)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, w)
}

type updateWorkerInput struct {
	Name      string `json:"name"`
	UserEmail string `json:"userEmail"`
	Phone     string `json:"phone"`
	Active    bool   `json:"active"`
}

func (h *workerHandlers) UpdateWorker(e echo.Context) error {
	id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid worker id")
	}
	var payload updateWorkerInput
	if err := e.Bind(&payload); err != nil {
		return validationError(e, err.Error())
	}
	w := &domain.MaintenanceWorker{
		Id:        id,
		Name:      payload.Name,
		UserEmail: payload.UserEmail,
		Phone:     payload.Phone,
		Active:    payload.Active,
	}
	code, err := h.s.Worker.UpdateWorker(e.Request().Context(), w)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, w)
}

func (h *workerHandlers) DeleteWorker(e echo.Context) error {
	id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid worker id")
	}
	code, err := h.s.Worker.DeleteWorker(e.Request().Context(), id)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, nil)
}

// ─── Request ─────────────────────────────────────────────────────────────────

type RequestHandlers interface {
	GetRequests(e echo.Context) error
	GetRequest(e echo.Context) error
	CreateRequest(e echo.Context) error
	AssignWorker(e echo.Context) error
	Resolve(e echo.Context) error
	Reject(e echo.Context) error
}

type requestHandlers struct{ s *services.Services }

var _ RequestHandlers = &requestHandlers{}

func NewRequestHandlers(s *services.Services) RequestHandlers {
	return &requestHandlers{s: s}
}

func (h *requestHandlers) GetRequests(e echo.Context) error {
	filters := repository.RequestFilters{}

	if v := e.QueryParam("requesterEmail"); v != "" {
		filters.RequesterEmail = &v
	}
	if v := e.QueryParam("categoryId"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			return validationError(e, "invalid categoryId")
		}
		filters.CategoryId = &id
	}
	if v := e.QueryParam("subcategoryId"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			return validationError(e, "invalid subcategoryId")
		}
		filters.SubcategoryId = &id
	}
	if v := e.QueryParam("status"); v != "" {
		s := domain.MaintenanceStatus(v)
		filters.Status = &s
	}

	reqs, code, err := h.s.Request.GetRequests(e.Request().Context(), filters)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, reqs)
}

func (h *requestHandlers) GetRequest(e echo.Context) error {
	id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid request id")
	}
	req, code, err := h.s.Request.GetRequest(e.Request().Context(), id)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, req)
}

type createRequestInput struct {
	RequesterEmail string  `json:"requesterEmail"`
	RequesterName  string  `json:"requesterName"`
	Location       string  `json:"location"`
	CategoryId     string  `json:"categoryId"`
	SubcategoryId  string  `json:"subcategoryId"`
	DetailId       *string `json:"detailId"`
	Description    string  `json:"description"`
}

func (h *requestHandlers) CreateRequest(e echo.Context) error {
	var payload createRequestInput
	if err := e.Bind(&payload); err != nil {
		return validationError(e, err.Error())
	}

	categoryId, err := uuid.Parse(payload.CategoryId)
	if err != nil {
		return validationError(e, "invalid categoryId")
	}
	subcategoryId, err := uuid.Parse(payload.SubcategoryId)
	if err != nil {
		return validationError(e, "invalid subcategoryId")
	}

	req := &domain.MaintenanceRequest{
		RequesterEmail: payload.RequesterEmail,
		RequesterName:  payload.RequesterName,
		Location:       payload.Location,
		CategoryId:     categoryId,
		SubcategoryId:  subcategoryId,
		Description:    payload.Description,
	}

	if payload.DetailId != nil {
		did, err := uuid.Parse(*payload.DetailId)
		if err != nil {
			return validationError(e, "invalid detailId")
		}
		req.DetailId = &did
	}

	code, err := h.s.Request.CreateRequest(e.Request().Context(), req)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, req)
}

type assignWorkerInput struct {
	WorkerId  string `json:"workerId"`
	UserEmail string `json:"userEmail"`
}

func (h *requestHandlers) AssignWorker(e echo.Context) error {
	requestId, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid request id")
	}
	var payload assignWorkerInput
	if err := e.Bind(&payload); err != nil {
		return validationError(e, err.Error())
	}
	workerId, err := uuid.Parse(payload.WorkerId)
	if err != nil {
		return validationError(e, "invalid workerId")
	}
	code, err := h.s.Request.AssignWorker(e.Request().Context(), requestId, workerId, payload.UserEmail)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, nil)
}

type resolveInput struct {
	ResolutionNotes string `json:"resolutionNotes"`
	UserEmail       string `json:"userEmail"`
}

func (h *requestHandlers) Resolve(e echo.Context) error {
	requestId, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid request id")
	}
	var payload resolveInput
	if err := e.Bind(&payload); err != nil {
		return validationError(e, err.Error())
	}
	code, err := h.s.Request.Resolve(e.Request().Context(), requestId, payload.ResolutionNotes, payload.UserEmail)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, nil)
}

type rejectInput struct {
	Comments  string `json:"comments"`
	UserEmail string `json:"userEmail"`
}

func (h *requestHandlers) Reject(e echo.Context) error {
	requestId, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return validationError(e, "invalid request id")
	}
	var payload rejectInput
	if err := e.Bind(&payload); err != nil {
		return validationError(e, err.Error())
	}
	code, err := h.s.Request.Reject(e.Request().Context(), requestId, payload.Comments, payload.UserEmail)
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, nil)
}

// ─── Config ──────────────────────────────────────────────────────────────────

type ConfigHandlers interface {
	GetConfig(e echo.Context) error
	UpdateConfig(e echo.Context) error
}

type configHandlers struct{ s *services.Services }

var _ ConfigHandlers = &configHandlers{}

func NewConfigHandlers(s *services.Services) ConfigHandlers {
	return &configHandlers{s: s}
}

func (h *configHandlers) GetConfig(e echo.Context) error {
	c, code, err := h.s.Config.GetConfig(e.Request().Context())
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, c)
}

type updateConfigInput struct {
	DeanEmail string `json:"deanEmail"`
}

func (h *configHandlers) UpdateConfig(e echo.Context) error {
	var payload updateConfigInput
	if err := e.Bind(&payload); err != nil {
		return validationError(e, err.Error())
	}
	code, err := h.s.Config.UpdateConfig(e.Request().Context(), &domain.MaintenanceConfig{
		DeanEmail: payload.DeanEmail,
	})
	if err != nil {
		return httpTransport.ErrWithMsg(e, code, err.Error(), nil)
	}
	return httpTransport.Ok(e, nil)
}