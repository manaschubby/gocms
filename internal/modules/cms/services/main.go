package services

import "github.com/manaschubby/gocms/internal/modules/cms/repository"

type Services struct {
	ContentType ContentTypeService
}

func New(r repository.Repositories) *Services {
	return &Services{
		ContentType: NewContentTypeService(r),
	}
}
