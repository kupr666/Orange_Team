include .env
export

export PROJECT_ROOT=${shell pwd}

app-run:
	@go run ./cmd/api

env-up:
	@docker compose up -d app-postgres

env-down:
	@docker compose down app-postgres

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Missing required parameter seq. Example: make migrate-create seq=init"; \
		exit 1; \
	fi;
	docker compose run --rm app-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Missing required parameter action. Example: make migrate-action action=up action_args=1"; \
		exit 1; \
	fi; 
	@docker compose run  --rm app-postgres-migrate \
		-path /migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@app-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
		"$(action)" $(action_args)

loki-up:
	@docker compose up -d loki

loki-logs:
	@docker compose logs -f loki

loki-stop:
	@docker compose stop loki

alloy-up:
	@docker compose up -d loki alloy

alloy-logs:
	@docker compose logs -f alloy

alloy-stop:
	@docker compose stop alloy loki

.PHONY: dev-frontend
dev-frontend:
	@echo "Переходим в frontend..."
	cd frontend && \
	( [ -f .env ] || cp .env.example .env ) && \
	echo "Устанавливаем зависимости..." && \
	npm ci && \
	echo "Запускаем dev-сервер..." && \
	npm run dev
