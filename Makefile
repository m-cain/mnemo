.PHONY: all db backend frontend clean migrate-create migrate-up migrate-down-one migrate-down-to migrate-status migrate-version

all: db backend frontend

backend:
	@echo "Starting backend with Air..."
	cd backend && DATABASE_URL="postgres://mnemo_user:mnemo_password@localhost:5432/mnemo_db?sslmode=disable" air

frontend:
	@echo "Starting frontend development server..."
	cd frontend && npm run dev

# DB_DRIVER and DB_DSN might be set from dbconf.yml or env vars for goose
# For simplicity, we can specify the environment from dbconf.yml
DB_ENV ?= development
MIGRATIONS_DIR := backend/migrations

# Set environment variables for goose based on the development environment in dbconf.yml
GOOSE_DRIVER := postgres
GOOSE_DBSTRING := postgres://mnemo_user:mnemo_password@localhost:5432/mnemo_db?sslmode=disable

migrate-create:
	@read -p "Enter migration name: " name; \
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(MIGRATIONS_DIR) create $$name sql

migrate-up:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(MIGRATIONS_DIR) up

migrate-down-one: # Goose 'down' rolls back one by default
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(MIGRATIONS_DIR) down

migrate-down-to:
	@read -p "Enter version to migrate down to: " version; \
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(MIGRATIONS_DIR) down-to $$version

migrate-status:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(MIGRATIONS_DIR) status

migrate-version:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(MIGRATIONS_DIR) version
