# AGENTS.md — Contexto del proyecto adventist-stack

Guía para agentes de IA (Claude, Gemini, etc.) que trabajan en este repo. Leer antes de tocar código.

---

## 🔴 En qué estamos AHORA (handoff — leer primero)

**Feature activa:** miembros (`services/members`), rama `feature/1-Members_Management`, Ticket #1.

**Roadmap de la feature:**
- ✅ **Paso 4 — dominio + repo + service de members**: hecho, compila, `go vet` limpio. Usa `entity.Base`, VOs `Email`/`Phone`, enums `Status`/`Gender`, validación con `errors.Join`.
- ✅ **Paso 5 — capa gRPC + Gateway (COMPLETADO)**:
  - ✅ **5.1** `services/members/internal/grpc/server.go` + `mapper.go` (completado).
  - ✅ **5.2** `services/members/cmd/server/main.go` (wiring completado).
  - ✅ **5.3** Gateway:
    - ✅ **5.3.1** Cliente gRPC en `gateway/internal/clients/members_client.go`.
    - ✅ **5.3.2** Handlers HTTP divididos al estilo C# (`members_handler.go`, `members_dto.go`, `members_mapper.go`).
    - ✅ **5.3.3** Registro de rutas en `gateway/internal/router/router.go` bajo `/api/v1/members`.

**Estado actual:** Feature #1 (Members Management) **100% completada y lista para Pull Request a `main`**.

**Cómo se está guiando (teaching mode):** el usuario **escribe members él mismo** para aprender. El agente:
1. Da **moldes** en `products` (servicio de referencia) y `auth`, explica el porqué, y **revisa/corrige** lo que el usuario escribe (a veces manda screenshots de errores del compilador).
2. **No** implementa members por él salvo que lo pida explícito. Sí toca `pkg/`, `products`, `auth` libremente.
3. Explica conceptos en profundidad (VOs, DDD, enums de Go, frontera) — el usuario pregunta mucho el "por qué".
4. **Siempre sugiere el siguiente paso en español** al terminar.

**Decisiones ya tomadas en esta feature** (no re-litigar):
- Enums `Status`/`Gender` en el dominio usan `type Status int` / `type Gender int` con `iota + 1` y métodos `String()` / `IsValid()`.
- Clientes gRPC en el Gateway divididos por archivo en `gateway/internal/clients/` (`auth_client.go`, `members_client.go`, `products_client.go`, `notifications_client.go`).
- `entity.Base` embebido en las 3 entidades (Member, User, Product) con campos de auditoría. Soft-delete/actor: columnas creadas pero lógica **no** activada aún.
- Helpers genéricos extraídos a `pkg/protoconv` (timestamps) y `pkg/ptr` (`Deref`). Los específicos (statusFromProto, toStatus) quedan locales.
- Pendiente futuro (no ahora): `ParseStatus` de string para el gateway (JSON), soft-delete real, wiring de actor desde JWT.

---

## Qué es

Backend de gestión para la Iglesia Adventista. **Monorepo Go de microservicios con gRPC**, más un cliente Flutter (futuro, en `client/`). Dominio: miembros, escuela sabática, boletines, calendario, peticiones de oración.

## Stack

- **Go 1.23.4**, un solo `go.mod` en la raíz (`module github.com/Jeudry/adventist-stack`).
- **gRPC + protobuf** (protoc, NO buf para generar). Protos en `proto/<svc>/v1/`, generado en `gen/<svc>/v1/`. Se genera con `make proto`.
- **Gateway REST** con chi/v5 (traduce JSON ↔ gRPC).
- **PostgreSQL** con pgx/v5 + pgxpool. **sqlc** (queries desde SQL). **golang-migrate** (migraciones embebidas con `embed.FS`, tabla por servicio: `<svc>_schema_migrations`).
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
    repository/               # repo + mappers (db ↔ domain)
    service/                   # orquesta: Normalize → Validate → repo
    grpc/                      # server.go (RPCs finos) + mapper.go (proto ↔ domain)
```

### Paquetes compartidos (`pkg/`)
- `pkg/entity` → `Base` (ID + CreatedAt/UpdatedAt/DeletedAt/CreatedBy/UpdatedBy/DeletedBy). Se **embebe** en cada entidad (`type Member struct { entity.Base; ... }`).
- `pkg/vo` → Value Objects: `Email`, `Phone`. Inmutables, campo privado, `NewX`/`NewOptionalX`/`Ptr()`/`String()`/`IsZero()`.
- `pkg/pagination` → `ListRequest` → `ToQuery()` → `Query`; `Page[T]` + `NewPage`.
- `pkg/protoconv` → `TimeFromProto`/`TimeToProto` (`*timestamppb.Timestamp` ↔ `*time.Time`).
- `pkg/ptr` → `Deref[T](*T) T` (genérico).
- `pkg/strutil` → `TrimPtr`. `pkg/config`, `pkg/database`, `pkg/jwt`, `pkg/logger`, `pkg/mailer`, `pkg/redis`, `pkg/middleware`.

## Convenciones (IMPORTANTE)

### Dominio / DDD-lite
- **Entidad rica**: `Normalize()` (trim/lowercase, defaults) y `Validate() error`. Validación con **constantes** (`NameMaxLen`) + **una función por propiedad** (`validateFirstName`) + **`errors.Join`**.
- **Value Objects** solo donde la validación es real y reusable (email, phone). Se construyen **en la frontera** con `vo.NewX` (por eso esos mappers devuelven `error`). Campos planos (`FirstName`, `Address *string`) van directo.
- **Enums = typed string + constantes + `IsValid()`**, NO struct-VO.
- **Regla de frontera**: ningún tipo externo (proto/JSON/DB) entra crudo al dominio. Se convierte en el borde: `vo.NewOptionalEmail`, `protoconv.TimeFromProto`, `statusFromProto`.
- **Métodos de comportamiento** en la entidad (ubiquitous language): ej. `User.SetPassword/Authenticate/IsAdmin`.
- La política de negocio vive en el **dominio**, no en el service. El service solo orquesta.

### Repositorio
- Usa sqlc. Mappers `toDomain` (db → domain) y `toCreateParams`/`toUpdateParams`/`toListParams`.
- `toDomain` **devuelve `error`** cuando rehidrata VOs.
- No encontrado: `errors.Is(err, pgx.ErrNoRows)` → error de dominio (`ErrXNotFound`).

### gRPC (transporte)
- `server.go`: struct embebe `Unimplemented<Svc>ServiceServer`, RPCs **finos** (input → domain → svc → domain → proto, con `toStatus`).
- `mapper.go`: `<x>FromCreate/Update`, `<x>ToProto`, `statusFromProto`/`statusToProto`, `toStatus` (dominio → `codes`).
- Genérico y sin acople → `pkg/`. Referencia tipos del servicio → local.

### Idioma
- **Código, comentarios y errores: en INGLÉS.** Sin comentarios verbosos de IA.
- **Docs y comunicación con el usuario: en ESPAÑOL.**

### Git (obligatorio)
- Cuenta: **Jeudry**. Remoto: GitHub (`gh`). Base: `main`. Rama actual: `feature/1-Members_Management`.
- **Un commit atómico + push inmediato por cambio pedido.** Conventional commits. NUNCA `git add .` — stage explícito.
- Terminar mensajes con `Co-Authored-By: <modelo>`. Verificar con `go build ./...` / `go test ./...` / `go vet` antes de commitear.

## Estado por servicio

- **auth** ✅ — end-to-end (domain con `Email` VO + comportamiento, sqlc, gRPC, gateway). Tests verdes.
- **members** ✅ — end-to-end (domain, sqlc, gRPC server/wiring, gateway client/handlers/routes). 100% completado.
- **notifications** — Redis pub/sub + mailer.

## Comandos

```bash
make proto          # genera Go desde los .proto
make tidy           # go mod tidy
go build ./...      # compila todo
go test ./...       # tests
sqlc generate       # (dentro de services/<svc>/) regenera tras cambiar migraciones/query.sql
```
