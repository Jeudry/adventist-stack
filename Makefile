.DEFAULT_GOAL := help
MODULE := github.com/Jeudry/adventist-stack

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
