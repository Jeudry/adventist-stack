.DEFAULT_GOAL := help
MODULE := github.com/Jeudry/adventist-stack
# Lee DATABASE_URL desde .env para los comandos de migración.
DB_URL := $(shell grep -E '^DATABASE_URL=' .env 2>/dev/null | cut -d= -f2-)
# DSN con tabla de migraciones POR SERVICIO, para que sus versiones no choquen
# entre servicios que comparten la misma base (usa '=' para expandir $(svc) al usarse).
MIGRATE_DB = "$(DB_URL)&x-migrations-table=$(svc)_schema_migrations"

.PHONY: help
help: ## Muestra esta ayuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

.PHONY: proto
proto: ## Genera el código Go desde los .proto
	protoc -I proto \
		--go_out=. --go_opt=module=$(MODULE) \
		--go-grpc_out=. --go-grpc_opt=module=$(MODULE) \
		$$(find proto -name '*.proto')

.PHONY: tidy
tidy: ## Ordena dependencias (go mod tidy)
	go mod tidy

.PHONY: migration
migration: ## Crea una migración vacía: make migration svc=<servicio> name=<nombre>
	@if [ -z "$(svc)" ] || [ -z "$(name)" ]; then \
		echo "Uso:  make migration svc=<servicio> name=<nombre>"; \
		echo "Ej:   make migration svc=members name=create_members_table"; \
		exit 1; \
	fi
	@mkdir -p services/$(svc)/migrations
	@migrate create -ext sql -dir services/$(svc)/migrations -seq $(name)
	@echo "✓ Migración creada en services/$(svc)/migrations (llená los .up.sql y .down.sql)"

.PHONY: migrate-up
migrate-up: ## Aplica migraciones pendientes: make migrate-up svc=<servicio>
	@[ -n "$(svc)" ] || { echo "Uso: make migrate-up svc=<servicio>"; exit 1; }
	migrate -path services/$(svc)/migrations -database $(MIGRATE_DB) up

.PHONY: migrate-down
migrate-down: ## Retrocede N migraciones (default 1): make migrate-down svc=<servicio> [n=1]
	@[ -n "$(svc)" ] || { echo "Uso: make migrate-down svc=<servicio> [n=1]"; exit 1; }
	migrate -path services/$(svc)/migrations -database $(MIGRATE_DB) down $(or $(n),1)

.PHONY: migrate-version
migrate-version: ## Muestra la versión aplicada: make migrate-version svc=<servicio>
	@[ -n "$(svc)" ] || { echo "Uso: make migrate-version svc=<servicio>"; exit 1; }
	migrate -path services/$(svc)/migrations -database $(MIGRATE_DB) version

.PHONY: migrate-force
migrate-force: ## Fuerza una versión (arregla estado 'dirty'): make migrate-force svc=<servicio> version=<N>
	@[ -n "$(svc)" ] && [ -n "$(version)" ] || { echo "Uso: make migrate-force svc=<servicio> version=<N>"; exit 1; }
	migrate -path services/$(svc)/migrations -database $(MIGRATE_DB) force $(version)

.PHONY: build
build: ## Compila los tres binarios en ./bin
	go build -o bin/gateway ./gateway/cmd/server
	go build -o bin/auth ./services/auth/cmd/server
	go build -o bin/notifications ./services/notifications/cmd/server

.PHONY: run-gateway
run-gateway: ## Corre el gateway
	go run ./gateway/cmd/server

.PHONY: run-auth
run-auth: ## Corre el servicio auth
	go run ./services/auth/cmd/server

.PHONY: run-notifications
run-notifications: ## Corre el servicio notifications
	go run ./services/notifications/cmd/server

.PHONY: test
test: ## Corre los tests
	go test ./...

.PHONY: infra-up
infra-up: ## Levanta postgres, redis y mailhog (docker-compose)
	docker compose -f deploy/docker-compose.yml up -d

.PHONY: infra-down
infra-down: ## Detiene la infraestructura (sin borrar datos)
	docker compose -f deploy/docker-compose.yml down

.PHONY: client
client: ## Corre la app Flutter (client/)
	cd client && flutter run
