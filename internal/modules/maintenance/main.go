package maintenance

import (
	"github.com/jmoiron/sqlx"
	"github.com/manaschubby/gocms/internal/modules/maintenance/handlers"
	"github.com/manaschubby/gocms/internal/modules/maintenance/repository"
)

type Maintenance struct {
	Repositories *repository.Repositories
	Handlers     *handlers.Handlers
}

func Init(db *sqlx.DB) *Maintenance {
	r := repository.Init(db)
	return &Maintenance{
		Repositories: &r,
		Handlers:     handlers.New(r),
	}
}