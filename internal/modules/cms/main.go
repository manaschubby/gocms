package cms

import (
	"github.com/jmoiron/sqlx"
	"github.com/manaschubby/gocms/internal/config"
	"github.com/manaschubby/gocms/internal/modules/cms/handlers"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
)

type CMS struct {
	Repositories *repository.Repositories
	Handlers     *handlers.Handlers
}

func Init(cfg *config.Config, db *sqlx.DB) *CMS {
	// Initialize Repos
	r := repository.Init(db)

	return &CMS{
		Repositories: &r,
		Handlers:     handlers.New(r),
	}
}
