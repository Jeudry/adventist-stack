# GEMINI.md — Contexto del proyecto adventist-stack

Guía para agentes de IA que trabajan en este repo. Leer antes de tocar código.

## Qué es

Backend de gestión para la Iglesia Adventista. **Monorepo Go de microservicios con gRPC**, más un cliente Flutter (futuro, en `client/`). Dominio: miembros, escuela sabática, boletines, calendario de eventos, peticiones de oración.

## Stack

- **Go 1.23.4**, un solo `go.mod` en la raíz (`module github.com/Jeudry/adventist-stack`).
- **gRPC + protobuf** (protoc, NO buf para generar). Protos en `proto/<svc>/v1/`, generado en `gen/<svc>/v1/`. Se genera con `make proto`.
- **Gateway REST** con chi/v5 (traduce JSON ↔ gRPC).
- **PostgreSQL** con pgx/v5 + pgxpool. **sqlc** para generar queries desde SQL. **golang-migrate** (migraciones embebidas con `embed.FS`, tabla por servicio: `<svc>_schema_migrations`).
- **JWT** (HS256) + bcrypt en auth. **Redis** (pub/sub) y **SMTP/MailHog** en notifications.

## Arquitectura

Cada servicio en `services/<svc>/` sigue Clean/Hexagonal + **DDD-lite**:

```
services/<svc>/
  cmd/server/main.go          # wiring: config → migrate → connect → repo → service → grpc
  embed.go                    # //go:embed migrations/*.sql → MigrationsFS
  migrations/                 # 000001_*.up.sql / .down.sql
  query.sql                   # queries para sqlc
  sqlc.yaml
  internal/
    domain/                   # entidad rica: Normalize/Validate, factories, enums, VOs
    db/                        # generado por sqlc (NO editar)
    repository/                # repo + mappers (db ↔ domain)
    service/                   # orquesta: Normalize → Validate → repo
    grpc/                      # server.go (RPCs finos) + mapper.go (proto ↔ domain)
```

### Paquetes compartidos (`pkg/`)
- `pkg/entity` → `Base` (ID + CreatedAt/UpdatedAt/DeletedAt/CreatedBy/UpdatedBy/DeletedBy). Se **embebe** en cada entidad de dominio (`type Member struct { entity.Base; ... }`).
- `pkg/vo` → Value Objects: `Email`, `Phone`. Inmutables, campo privado, `NewX`/`NewOptionalX`/`Ptr()`/`String()`/`IsZero()`.
- `pkg/pagination` → `ListRequest` → `ToQuery()` → `Query`; `Page[T]` + `NewPage`.
- `pkg/protoconv` → `TimeFromProto`/`TimeToProto` (`*timestamppb.Timestamp` ↔ `*time.Time`).
- `pkg/ptr` → `Deref[T](*T) T` (genérico).
- `pkg/strutil` → `TrimPtr`. `pkg/config`, `pkg/database`, `pkg/jwt`, `pkg/logger`, `pkg/mailer`, `pkg/redis`, `pkg/middleware`.

## Convenciones (IMPORTANTE)

### Dominio / DDD-lite
- **Entidad rica**: métodos `Normalize()` (trim/lowercase, defaults) y `Validate() error`. Validación con **constantes** (`NameMaxLen`) + **una función por propiedad** (`validateFirstName`) + **`errors.Join`** para acumular errores.
- **Value Objects** solo donde la validación es real y reusable (email, phone). Se construyen **en la frontera** con `vo.NewX` (por eso esos mappers devuelven `error`). Campos planos (`FirstName`, `Address *string`) se asignan directo.
- **Enums = typed string + constantes + `IsValid()`**, NO struct-VO (perdería constantes/switch). Ej: `type Status string`, `type Gender string`.
- **Regla de frontera**: ningún tipo externo (proto/JSON/DB) entra crudo al dominio. Se convierte en el borde: `vo.NewOptionalEmail`, `protoconv.TimeFromProto`, `statusFromProto`.
- **Métodos de comportamiento** en la entidad (ubiquitous language): ej. `User.SetPassword/Authenticate/IsAdmin`.
- La política de negocio vive en el **dominio**, no en el service. El service solo orquesta.

### Repositorio
- Usa sqlc. Mappers `toDomain` (db → domain) y `toCreateParams`/`toUpdateParams`/`toListParams`.
- `toDomain` **devuelve `error`** cuando rehidrata VOs (reconstruir un VO puede fallar).
- Errores de "no encontrado": `errors.Is(err, pgx.ErrNoRows)` → error de dominio (`ErrXNotFound`).

### gRPC (capa de transporte)
- `server.go`: struct embebe `Unimplemented<Svc>ServiceServer`, RPCs **finos** (input → domain → svc → domain → proto, con `toStatus`).
- `mapper.go`: helpers de mapeo. `<x>FromCreate/Update` (proto → domain), `<x>ToProto`, `statusFromProto`/`statusToProto`, `toStatus` (errores dominio → `codes`).
- Genérico y sin acople → va a `pkg/` (`protoconv`, `ptr`). Referencia tipos del servicio (enums, sentinels) → local.

### Idioma
- **Código, comentarios y mensajes de error: en INGLÉS.** Sin comentarios verbosos de IA.
- **Docs (`docs/`) y comunicación con el usuario: en ESPAÑOL.**

### Git (obligatorio)
- Cuenta: **Jeudry**. Remoto: GitHub (`gh` CLI). Rama base: `main`. Rama actual de trabajo: `feature/1-Members_Management`.
- **Un commit atómico + push inmediato por cada cambio pedido.** Conventional commits (`feat`, `fix`, `refactor`, `chore`). NUNCA `git add .` — stage explícito.
- Terminar los mensajes de commit con: `Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>` (ajustar al modelo real).
- Verificar antes de commitear: `go build ./...`, `go test ./...`, `go vet`.

## Estado actual

- **auth** ✅ — end-to-end (domain con `Email` VO + comportamiento, sqlc, gRPC, gateway). Tests de dominio verdes.
- **products** 📐 — servicio de **REFERENCIA/molde** (no se despliega). Tiene el patrón completo: sqlc, domain, service, **gRPC server (`server.go` + `mapper.go`) + `cmd/server/main.go`**. Copiar de acá para members.
- **members** 🚧 — `domain` + `repository` + `service` listos y compilan (usa `entity.Base`, VOs `Email`/`Phone`, enums `Status`/`Gender`). **En progreso: capa gRPC** (`internal/grpc/server.go` + `mapper.go`) — la escribe el usuario siguiendo el molde de products. Pendiente: `cmd/server/main.go` y wiring en el gateway.
- **notifications** — Redis pub/sub + mailer.

### Diferencias members vs el molde products
1. members tiene **VOs** (`Email`/`Phone`) → `memberFromCreate` devuelve `(domain.Member, error)` y llama `vo.NewOptionalEmail/Phone`. products no tiene VOs (sin error).
2. `toStatus` de members suma `vo.ErrInvalidEmail`/`vo.ErrInvalidPhone`.
3. members tiene 2 fechas opcionales (`BirthDate`, `BaptismDate`) y `Gender`.

## Dinámica de trabajo

El usuario está **aprendiendo** y escribe members él mismo (teaching mode). El agente provee moldes en products, explica y revisa — no hace el trabajo de members por él salvo que lo pida explícito. **Siempre sugerir el siguiente paso en español** al terminar.

## Comandos

```bash
make proto          # genera Go desde los .proto
make tidy           # go mod tidy
go build ./...      # compila todo
go test ./...       # tests
sqlc generate       # (dentro de services/<svc>/) regenera queries tras cambiar migraciones/query.sql
```
