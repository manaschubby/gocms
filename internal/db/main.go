package db

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/manaschubby/gocms/internal/config"

	_ "github.com/lib/pq"
)

func Connect(cfg config.Config) (*sqlx.DB, error) {
	datasourceURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB, cfg.PostgresSSLMode)
	db, err := sqlx.Connect("postgres", datasourceURL)
	if err != nil {
		return nil, errors.New("Failed to connect to the Database using provided Credentials: " + err.Error())
	}
	return db, nil
}
