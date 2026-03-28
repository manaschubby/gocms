package services

import (
	"github.com/manaschubby/gocms/internal/modules/maintenance/repository"
)

type Services struct {
	Category    CategoryService
	Subcategory SubcategoryService
	Detail      DetailService
	Worker      WorkerService
	Request     RequestService
	Config      ConfigService
}

func New(r repository.Repositories) *Services {
	return &Services{
		Category:    NewCategoryService(r),
		Subcategory: NewSubcategoryService(r),
		Detail:      NewDetailService(r),
		Worker:      NewWorkerService(r),
		Request:     NewRequestService(r),
		Config:      NewConfigService(r),
	}
}