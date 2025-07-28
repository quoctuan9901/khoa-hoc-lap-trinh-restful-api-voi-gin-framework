include .env
export

CONN_STRING = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

MIGRATION_DIRS = internal/db/migrations

ENV_FILE=.env

PROD_COMPOSE=docker-compose.prod.yml
NOAPP_COMPOSE=docker-compose.noapp.yml
DEV_COMPOSE=docker-compose.dev.yml

# Import database
importdb:
	docker exec -i postgres-db psql -U root -d master-golang < ./backupdb-master-golang.sql

# Export database
exportdb:
	docker exec -i postgres-db pg_dump -U root -d master-golang > ./backupdb-master-golang.sql

# Build a binary file
build:
	go build -o bin/myapp ./cmd/api

run-binary:
	./bin/myapp

# Run server
server:
	go run ./cmd/api

# Generate sqlc
sqlc:
	sqlc generate

# Create a new migration (make migrate-create NAME=profiles)
migrate-create:
	migrate create -ext sql -dir $(MIGRATION_DIRS) -seq $(NAME)

# Run all pending migration (make migrate-up)
migrate-up:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" up

# Rollback the last migration
migrate-down:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down 1

# Rollback N migrations
migrate-down-n:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down $(N)

# Force migration version (use with caution example: make migrate-force VERSION=1) 
migrate-force:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" force $(VERSION)

# Drop everything (include schema migration)
migrate-drop:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" drop

# Apply specific migration version (make migrate-goto VERSION=1)
migrate-goto:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" goto $(VERSION)

# Build container for noapp
noapp:
	docker compose -f $(NOAPP_COMPOSE) down
	docker compose -f $(NOAPP_COMPOSE) --env-file $(ENV_FILE) up -d --build

# Stop all container noapp
stop-noapp:
	docker compose -f $(NOAPP_COMPOSE) down

# Build container for dev
dev:
	docker compose -f $(DEV_COMPOSE) down
	docker compose -f $(DEV_COMPOSE) --env-file $(ENV_FILE) up --build

# Build container for production
prod:
	docker compose -f $(PROD_COMPOSE) down
	docker compose -f $(PROD_COMPOSE) --env-file $(ENV_FILE) up -d --build

# Stop all container production
stop-prod:
	docker compose -f $(PROD_COMPOSE) down

# View logs container production
logs-prod:
	docker compose -f $(PROD_COMPOSE) logs -f --tail=100

# Go to container golang-api
bash:
	docker exec -it golang-api /bin/sh

.PHONY: importdb exportdb server migrate-create migrate-up migrate-down migrate-force migrate-drop migrate-goto migrate-down-n sqlc build run-binary prod stop-prod logs-prod bash noapp stop-noapp dev