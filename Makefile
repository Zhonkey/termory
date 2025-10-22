ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
CMD ?= ./main.go

# ============ DEVELOPMENT COMMANDS ============
build:
	docker compose build

up: build
	docker compose up -d

sh: up
	docker compose exec app sh -c "/bin/bash"

npm: up
	docker compose exec vue sh -c "/bin/bash"

db: up
	docker compose exec db psql -U trainer -d trainer

down:
	docker compose down --remove-orphans

debug:
	docker compose exec app sh -c 'go build -gcflags "all=-N -l" -o main_debug.bin .'
	docker compose exec app sh -c 'dlv exec ./main_debug.bin --headless --listen=:2345 --api-version=2 --accept-multiclient'

# ============ PRODUCTION COMMANDS ============

# Сборка production образа
build-prod:
	docker compose -f docker-compose.prod.yml build

# Запуск production контейнеров
up-prod: build-prod
	docker compose -f docker-compose.prod.yml up -d

# Остановка production контейнеров
down-prod:
	docker compose -f docker-compose.prod.yml down --remove-orphans

# Просмотр логов production
logs-prod:
	docker compose -f docker-compose.prod.yml logs -f app

# Просмотр логов БД
logs-db-prod:
	docker compose -f docker-compose.prod.yml logs -f db

# Запуск миграций в production
migrate-prod:
	docker compose -f docker-compose.prod.yml exec app /app/migrate up

# Откат миграций в production
migrate-down-prod:
	docker compose -f docker-compose.prod.yml exec app /app/migrate down

# Статус миграций в production
migrate-status-prod:
	docker compose -f docker-compose.prod.yml exec app /app/migrate status

# Полный production деплой (сборка + запуск + миграции)
deploy-prod: up-prod migrate-prod
	@echo "✅ Production deployment completed successfully"

# Локальная сборка бинарников (без Docker)
build-local:
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/api ./cmd/api/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/migrate ./cmd/migrate/main.go
	@echo "✅ Binaries built successfully in bin/"

# Запуск локально собранного API (требует БД)
run-local: build-local
	@echo "Starting API server..."
	./bin/api

# Запуск миграций локально
migrate-up-local: build-local
	@echo "Running migrations..."
	./bin/migrate up

# Запуск миграций локально
migrate-down-local: build-local
	@echo "Running migrations..."
	./bin/migrate down

# Помощь
help:
	@echo "=== Development Commands ==="
	@echo "make build          - Build development Docker image"
	@echo "make up             - Start development containers"
	@echo "make down           - Stop development containers"
	@echo "make sh             - Open shell in app container"
	@echo "make npm            - Open shell in Vue container"
	@echo "make db             - Connect to PostgreSQL"
	@echo "make debug          - Start debugger"
	@echo ""
	@echo "=== Production Commands ==="
	@echo "make build-prod     - Build production Docker image"
	@echo "make up-prod        - Start production containers"
	@echo "make down-prod      - Stop production containers"
	@echo "make deploy-prod    - Full production deployment (build + run + migrate)"
	@echo "make migrate-prod   - Run migrations in production"
	@echo "make migrate-down-prod - Rollback migrations in production"
	@echo "make migrate-status-prod - Check migration status"
	@echo "make logs-prod      - View production API logs"
	@echo "make logs-db-prod   - View production DB logs"
	@echo ""
	@echo "=== Local Build Commands ==="
	@echo "make build-local    - Build binaries locally (no Docker)"
	@echo "make run-local      - Run API locally"
	@echo "make migrate-up-local  - Run migrations locally"
	@echo "make migrate-down-local  - Down migrations locally"
	@echo ""
	@echo "make help           - Show this help message"

.PHONY: build up sh npm db down debug build-prod up-prod down-prod logs-prod logs-db-prod migrate-prod migrate-down-prod migrate-status-prod deploy-prod build-local run-local migrate-local help
