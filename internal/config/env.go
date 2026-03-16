package config

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/manaschubby/gocms/internal/store"
)

func (s *Config) LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	// Database Specific Config
	// Load essential configurations
	db, err := getEnvOrError(store.POSTGRES_DB_VAR)
	if err != nil {
		return err
	}
	host, err := getEnvOrError(store.POSTGRES_HOST_VAR)
	if err != nil {
		return err
	}

	user, err := getEnvOrError(store.POSTGRES_USER_VAR)
	if err != nil {
		return err
	}

	password, err := getEnvOrError(store.POSTGRES_PASSWORD_VAR)
	if err != nil {
		return err
	}
	s.PostgresDB = db
	s.PostgresHost = host
	s.PostgresUser = user
	s.PostgresPassword = password

	// Load optional strings
	s.PostgresSSLMode = getEnvOr(store.POSTGRES_SSLMODE_VAR, "disable")
	s.PostgresPort = getEnvOr(store.POSTGRES_PORT_VAR, "5432")

	return nil
}

func getEnvOr(envVar, defaultValue string) string {
	env := os.Getenv(envVar)
	if strings.TrimSpace(env) == "" {
		return defaultValue
	} else {
		return env
	}
}
func getEnvOrError(envVar string) (string, error) {
	env := os.Getenv(envVar)
	if strings.TrimSpace(env) != "" {
		return env, nil
	} else {
		return "", errors.New("Required ENV var:" + envVar)
	}
}
