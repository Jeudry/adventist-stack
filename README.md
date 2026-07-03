# Adventist Stack

Plataforma de gestión de iglesia adventista. Backend de **microservicios en Go
con gRPC** y un **API Gateway REST**, con app **Flutter multiplataforma** en
`client/`.

## Arquitectura

```
┌────────────┐        REST/JSON        ┌──────────────┐
│  Flutter   │  ───────────────────▶   │   Gateway    │  (chi, JWT, Swagger,
│  (client/) │                         │  :8090       │   rate limit, CORS)
└────────────┘                         └──────┬───────┘
                                              │ gRPC
                          ┌───────────────────┼────────────────────┐
                          ▼                                        ▼
                  ┌──────────────┐                        ┌──────────────────┐
                  │   auth       │  :50051                │  notifications   │  :50052
                  │  (Postgres)  │                        │  (Redis, SMTP)   │
                  └──────────────┘                        └──────────────────┘
```

| Carpeta                 | Qué es                                                        |
| ----------------------- | ------------------------------------------------------------- |
| `gateway/`              | API Gateway REST (chi) → traduce a gRPC. Swagger, JWT, ratelimit |
| `services/auth/`        | Microservicio de autenticación (gRPC + PostgreSQL)            |
| `services/notifications/` | Correos con plantillas + notificaciones por Redis pub/sub   |
| `proto/`                | Contratos gRPC (`.proto`), fuente de verdad de las APIs        |
| `gen/`                  | Código Go generado desde los `.proto` (`make proto`)          |
| `pkg/`                  | Librerías compartidas: config, database, redis, jwt, mailer, logger, middleware |
| `client/`               | App Flutter (iOS, Android, web, macOS, Windows, Linux)        |
| `deploy/`               | `docker-compose.yml` (Postgres, Redis, MailHog)               |

Cada servicio sigue la misma estructura por capas:
`domain` (core) · `models` (DB) · `repository` · `service` (lógica) · `grpc` (transporte).

## Requisitos

- Go 1.23+
- Flutter 3.x
- Docker (para Postgres/Redis/MailHog)
- `protoc` + `protoc-gen-go` + `protoc-gen-go-grpc` (solo si regeneras protos)

## Puesta en marcha

```bash
# 1. Configura el entorno (secretos NUNCA se commitean)
cp .env.example .env
#    genera un JWT_SECRET fuerte:
openssl rand -base64 48        # pégalo en JWT_SECRET dentro de .env

# 2. Levanta la infraestructura (Postgres, Redis, MailHog)
make infra-up

# 3. Resuelve dependencias
make tidy

# 4. Corre cada servicio (en terminales separadas)
make run-auth
make run-notifications
make run-gateway

# 5. Corre el cliente Flutter
make client
```

- API Gateway: <http://localhost:8090>
- Swagger UI: <http://localhost:8090/swagger>
- Bandeja de correos (MailHog): <http://localhost:8025>

## Probar la API

```bash
# Registro
curl -X POST http://localhost:8090/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"pastor@iglesia.org","name":"Pastor","password":"secreto123"}'

# Login
curl -X POST http://localhost:8090/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"pastor@iglesia.org","password":"secreto123"}'

# Ruta protegida (usa el access_token del paso anterior)
curl http://localhost:8090/api/v1/me -H 'Authorization: Bearer <ACCESS_TOKEN>'
```

## Comandos útiles

```bash
make help          # lista todos los comandos
make proto         # regenera el código gRPC desde proto/
make build         # compila los 3 binarios en ./bin
make test          # corre los tests
```

## Seguridad de secretos

- `.env` está en `.gitignore`; solo se versiona `.env.example` con valores vacíos.
- La config se carga en structs tipadas y **valida al arranque** (`required`):
  la app falla rápido si falta una credencial, en vez de romperse en runtime.
- API keys y credenciales viven siempre en el entorno, nunca en el código.
