ifneq (,$(wildcard .env))
  include .env
  export
endif



.PHONY: dev dev-rebuild down
dev:
	@docker compose --project-directory . -f docker/dev/docker-compose.yml up -d
dev-rebuild:
	@docker compose --project-directory . -f docker/dev/docker-compose.yml up --build -d


GO_BIN_DIR := $(shell go env GOBIN)
ifeq ($(GO_BIN_DIR),)
  GO_BIN_DIR := $(shell go env GOPATH)/bin
endif

MIGRATE_BIN ?= $(shell command -v goose 2>/dev/null)
ifeq ($(MIGRATE_BIN),)
  MIGRATE_BIN := $(GO_BIN_DIR)/goose
endif



.PHONY: migrate-create
migrate-create:
	@echo hey $(NAME)

	@if [ -n "$(NAME)" ]; then \
		$(MIGRATE_BIN) create -s $(NAME) sql; \
	else \
		echo "You must provide a NAME for the migration. example: create_init_tables."; \
		exit 1; \
	fi


.PHONY: migrate-status
migrate-status:
	@$(MIGRATE_BIN) status
	@echo "Migrations complete."



.PHONY: migrate
migrate:
	@$(MIGRATE_BIN) up
	@echo "Migrations complete."
