# AGENTS.md — Contexto del proyecto adventist-stack

Guía para agentes de IA (Claude, Gemini, etc.) que trabajan en este repo. Leer antes de tocar código.

---

## 🔴 En qué estamos AHORA (handoff — leer primero)

**Feature activa:** escuela sabática (`services/sabbath_school`), rama `feature/3-Sabbath_School_Management`, Ticket #3.

**Roadmap de la feature:**
- ✅ **Paso 1 — Transporte HTTP**: `services/sabbath_school/internal/http/handler.go` con `Create` y `GetByID` ya expuestos. Ya no hay contrato `.proto`.
- 🔴 **Paso 2 — Migración + DB**: `services/sabbath_school/migrations` y `internal/db/` (queries pgx escritas a mano).
- 🔴 **Paso 3 — Dominio (`services/sabbath_school/internal/domain`)**: entidades y enums.
- 🔴 **Paso 4 — Repositorio (`services/sabbath_school/internal/repository`)**: repo + mappers.
- 🔴 **Paso 5 — Servicio (`services/sabbath_school/internal/service`)**: lógica de negocio.
- ✅ **Paso 6 — Capa HTTP + Wiring (`services/sabbath_school/internal/http` y `cmd/server/main.go`)**.
- ✅ **Paso 7 — Ruta en el gateway**: una línea en `gateway/internal/router/router.go` (`r.Mount("/sabbath-schools", ...)`) más la URL del servicio en `.env`. Ya no hay client, models ni mappers que escribir.

**Estado actual:** Iniciando **Feature #3 (Sabbath School Management)**.

**Cómo se está guiando (teaching mode):** el usuario **escribe members él mismo** para aprender. El agente:
1. Da **moldes** en `products` (servicio de referencia) y `auth`, explica el porqué, y **revisa/corrige** lo que el usuario escribe (a veces manda screenshots de errores del compilador).
2. **No** implementa members por él salvo que lo pida explícito. Sí toca `pkg/`, `products`, `auth` libremente.
3. Explica conceptos en profundidad (VOs, DDD, enums de Go, frontera) — el usuario pregunta mucho el "por qué".
4. **Siempre sugiere el siguiente paso en español** al terminar.

**Decisiones ya tomadas en esta feature** (no re-litigar):
- Enums del dominio (`Status`, `Gender`, `Role`) son **typed string**: `type Status string` con constantes literales (`StatusActive Status = "active"`), `String()`, `IsValid()` y `ParseX(string)`. El valor cero es `""` (unspecified), y las columnas se guardan como `VARCHAR` con `CHECK IN (...)`. Antes eran `int` con `iota + 1`; se migraron porque el número posicional no significaba nada en la DB ni en los logs.
- El Gateway **no deserializa nada**: `gateway/internal/proxy` reenvía por prefijo con `httputil.ReverseProxy`. No hay clientes, handlers ni mappers por servicio.
- `entity.Base` embebido en las 3 entidades (Member, User, Product) con campos de auditoría. Soft-delete/actor: columnas creadas pero lógica **no** activada aún.
- Helpers genéricos en `pkg/httpx` (JSON, identidad, paginación, fechas, `Serve`) y `pkg/ptr` (`Deref`). El mapeo de errores de dominio a códigos HTTP (`writeDomainError`) queda **local** a cada servicio.
- Pendiente futuro (no ahora): soft-delete real, wiring de actor desde JWT.
- Ya no hay protobuf en el proyecto: cliente, gateway y servicios hablan JSON. gRPC volvería solo si un servicio necesitara llamar a otro, que hoy no pasa.

---

## Qué es

Backend de gestión para la Iglesia Adventista. **Monorepo Go de microservicios sobre HTTP/JSON**, más un cliente Flutter (futuro, en `client/`). Dominio: miembros, escuela sabática, boletines, calendario, peticiones de oración.

## Stack

- **Go 1.23.4**, un solo `go.mod` en la raíz (`module github.com/Jeudry/adventist-stack`).
- **HTTP/JSON con chi/v5** en todos los servicios. Sin generación de código: cada servicio escribe sus handlers y sus DTOs.
- **Gateway** con chi/v5: enruta y **no abre el cuerpo de la petición**. Valida el JWT y afirma la identidad en `X-User-Id` / `X-User-Role`.
- **PostgreSQL** con pgx/v5, sin generación de código. El SQL vive en constantes dentro del repositorio y se ejecuta con **`pgx.StrictNamedArgs`**: parámetros por nombre (`@name`), nunca posicionales. **golang-migrate** (migraciones embebidas con `embed.FS`, tabla por servicio: `<svc>_schema_migrations`).
- **JWT** (HS256) + bcrypt en auth. **Redis** (pub/sub) y **SMTP/MailHog** en notifications.

## Arquitectura

Cada servicio en `services/<svc>/` sigue Clean/Hexagonal + **DDD-lite**:

```
services/<svc>/
  cmd/server/main.go          # wiring: config → migrate → connect → repo → service → http
  embed.go                    # //go:embed migrations/*.sql → MigrationsFS
  migrations/                 # 000001_*.up.sql / .down.sql
  internal/
    domain/                   # entidad rica: Normalize/Validate, factories, enums, VOs
    repository/               # SQL + Scan directo a domain (tiene el *pgxpool.Pool)
    service/                   # orquesta: Normalize → Validate → repo
    http/                      # handler.go: rutas chi, DTOs, VMs y writeDomainError
```

### Paquetes compartidos (`pkg/`)
- `pkg/entity` → `Base` (ID + CreatedAt/UpdatedAt/DeletedAt/CreatedBy/UpdatedBy/DeletedBy). Se **embebe** en cada entidad (`type Member struct { entity.Base; ... }`).
- `pkg/vo` → Value Objects: `Email`, `Phone`. Inmutables, campo privado, `NewX`/`NewOptionalX`/`Ptr()`/`String()`/`IsZero()`.
- `pkg/pagination` → `ListRequest` → `ToQuery()` → `Query`; `Page[T]` + `NewPage`.
- `pkg/httpx` → `NewAPI` (monta huma sobre chi), `In[T]`/`Out[T]`, `WriteJSON`/`WriteError`/`DecodeJSON`, `UserID`/`Role` (headers del gateway), `ListRequest`, `BaseVM`, `PageResponse`, `ParseDate`/`FormatDate`, y `Serve` (arranque + apagado ordenado).
- `pkg/ptr` → `Deref[T](*T) T` (genérico).
- `pkg/strutil` → `TrimPtr`. `pkg/config`, `pkg/database`, `pkg/jwt`, `pkg/logger`, `pkg/mailer`, `pkg/redis`, `pkg/middleware`.

## Convenciones (IMPORTANTE)

### Dominio / DDD-lite
- **Entidad rica**: `Normalize()` (trim/lowercase, defaults) y `Validate() error`. Validación con **constantes** (`NameMaxLen`) + **una función por propiedad** (`validateFirstName`) + **`errors.Join`**.
- **Value Objects** solo donde la validación es real y reusable (email, phone). Se construyen **en la frontera** con `vo.NewX` (por eso esos mappers devuelven `error`). Campos planos (`FirstName`, `Address *string`) van directo.
- **Enums = typed string + constantes + `IsValid()`**, NO struct-VO.
- **Regla de frontera**: ningún tipo externo (JSON/DB) entra crudo al dominio. Se convierte en el borde: `vo.NewOptionalEmail`, `httpx.ParseDate`, `domain.ParseStatus`.
- **Métodos de comportamiento** en la entidad (ubiquitous language): ej. `User.SetPassword/Authenticate/IsAdmin`.
- La política de negocio vive en el **dominio**, no en el service. El service solo orquesta.

### Repositorio
- El SQL vive en constantes al tope del `<entidad>_repository.go`, **siempre con parámetros nombrados** y ejecutado con `pgx.StrictNamedArgs`. Nunca posicionales: dos columnas del mismo tipo intercambiadas por error corren igual, y la base escribe la equivocada. Con args nombrados eso falla antes de tocar la base — *"argument X found in sql query but not present in StrictNamedArgs"*.
- La lista de columnas se declara **una sola vez** por entidad (`<entidad>Columns`) y la comparten el `SELECT`, el `INSERT ... RETURNING` y el `UPDATE ... RETURNING`. El orden de esa constante es el contrato que el `Scan` del mapper tiene que respetar.
- `<entidad>_mapper.go` tiene el `scan<Entidad>(row)` y el `writableArgs`. `writableArgs` lleva solo las columnas que el cliente escribe; el update le saca `created_by` y le pone `updated_by`, para que ninguna de las dos sentencias escriba la columna de auditoría de la otra.
- `pgx.ErrNoRows` → error de dominio (`ErrXNotFound`) es trabajo del repositorio.
- Una query que solo usan los tests (`select<Entidad>IncludingDeleted`) va como constante junto a las demás: el `_test.go` está en el mismo paquete y la usa con el mismo `scan`.
- `toDomain` **devuelve `error`** cuando rehidrata VOs (email, phone).

### HTTP (transporte)

**Los cinco servicios corren sobre huma.** chi sigue debajo como router; huma va encima.
- `handler.go`: `Register(api huma.API)` declara cada operación con `huma.Operation` — eso es lo que antes era el bloque de esa ruta en `openapi.yaml`, y **el spec se genera de ahí**.
- Los handlers reciben y devuelven **tipos**, no `http.ResponseWriter`. `In[T]`/`Out[T]` cubren los que solo mueven un cuerpo; si el endpoint lee ruta o query, necesita su struct con tags `path:` / `query:` — ahí es donde se declaran, y de ahí sale la validación y el spec.
- **La validación de entrada va en tags** (`minimum`, `maximum`, `format:"uuid"`): huma rechaza antes de llamar al handler. La de negocio sigue en el dominio.
- Los errores son **RFC 7807** (`{title,status,detail,errors}`) vía `domainError`. Validación fallida es **422**, no 400.
- En `main.go`: `humachi.New(router, config)` monta huma **sobre chi** — el router y su middleware no cambian. `config.CreateHooks = nil` quita el campo `$schema` que huma añade a cada respuesta.
- **Autorización por rol**: donde chi encerraba rutas en un `r.Group`, huma la pone en la operación que protege, con `Middlewares: huma.Middlewares{requireRole(api, "admin")}`. El middleware es `func(huma.Context, func(huma.Context))` y responde con `huma.WriteErr(api, ctx, 403, ...)`. Va junto al endpoint, no en el router — se lee sin reconstruir en qué grupo cayó cada ruta. **Ojo**: el servicio cree a ciegas en `X-User-Role`, así que su puerto no puede quedar expuesto; ahora ese header no solo revela datos, decide permisos.
- **El gateway no reescribe la ruta**: cada servicio declara la ruta pública completa, así que `proxy.To` reenvía tal cual.

- **Un campo sin `omitempty` es obligatorio para huma.** Si el spec lo declara opcional, marcarlo `required:"false"` — y que `Normalize()` le ponga el mismo valor por defecto que la columna.
- **Todo nombre JSON va en camelCase**, sin excepciones: `firstName`, `createdAt`, `pageSize`, `accessToken`. Se eligió para que los clientes Dart y Kotlin no necesiten anotar campo por campo. Un solo `snake_case` obliga al cliente a tratarlo aparte, así que no hay medias tintas.
- **No hay `openapi.yaml`.** El documento se arma en caliente: `gateway/internal/openapi` pide su spec a cada servicio (`/openapi.json`), lo fusiona con el del propio gateway y lo sirve en `/openapi.json`, que es lo que lee `/swagger`. Un servicio caído sale del documento con un warning en vez de romper la página. Se cachea 30s.
- **Consecuencia: el código es el contrato.** Agregar un campo o una operación aparece en la documentación sin tocar nada más. Ya no hay `contract_test.go` ni `pkg/apispec` — no hay dos verdades que reconciliar.
- Los servicios confían en `X-User-Id`/`X-User-Role`: **sus puertos no deben ser alcanzables desde fuera**. El gateway borra esos headers si vienen del cliente (probado en `gateway/internal/proxy/proxy_test.go`).
- En el mismo `handler.go`: `decode<X>` (JSON → dominio, construyendo los VOs), `to<X>VM` (dominio → JSON) y `writeDomainError` (error de dominio → código HTTP).
- Genérico y sin acople → `pkg/`. Referencia tipos del servicio → local.

### Idioma
- **Código, comentarios y errores: en INGLÉS.** Sin comentarios verbosos de IA.
- **Docs y comunicación con el usuario: en ESPAÑOL.**

### Git (obligatorio)
- Cuenta: **Jeudry**. Remoto: GitHub (`gh`). Base: `main`. Rama actual: `feature/1-Members_Management`.
- **Un commit atómico + push inmediato por cambio pedido.** Conventional commits. NUNCA `git add .` — stage explícito.
- Terminar mensajes con `Co-Authored-By: <modelo>`. Verificar con `go build ./...` / `go test ./...` / `go vet` antes de commitear.

## Estado por servicio

- **auth** ✅ — end-to-end (domain con `Email` VO + comportamiento, db, HTTP, gateway). Tests verdes.
- **members** ✅ — end-to-end (domain, db, HTTP server/wiring, ruta en el gateway). 100% completado.
- **notifications** — Redis pub/sub + mailer.

## Comandos

```bash
make tidy           # go mod tidy
go build ./...      # compila todo
go test ./...       # tests unitarios (sin DB)

make infra-up            # postgres + redis + mailhog
make migrate-all         # aplica las migraciones de los 4 servicios
make test-integration    # tests contra la DB real (build tag `integration`)
```

Los tests de integración viven junto a cada repositorio (`*_integration_test.go`, tag `integration`)
y usan `pkg/testdb`. Se limpian solos con `t.Cleanup`. `go test ./...` no los corre.

### Convenciones de testing

Los tests usan **testify**. Cuatro reglas, cada una medida contra la alternativa antes de adoptarla:

- **`require` cuando seguir no tiene sentido** (falló el setup, el `Create` devolvió error, vas a
  desreferenciar el resultado). **`assert` para las comprobaciones**, así una corrida muestra todos
  los fallos y no solo el primero. Es el `Fatalf`/`Errorf` de siempre con otro nombre.
- **Nunca `assert.True` si existe la assertion específica.** `True` recibe un bool ya evaluado y en
  el fallo solo dice *"Should be true"*. `Greater`, `Len`, `WithinDuration`, `ErrorIs`, `NotEmpty` y
  compañía imprimen los operandos.
- **Punteros: desreferenciar con `ptr.Deref` antes de comparar.** `assert.Equal` sobre un `*int`
  imprime la dirección de memoria, no el valor — es cómo Go formatea punteros, no un bug de testify.
  Si lo que se afirma es que algo quedó nil, pasar el valor en el mensaje:
  `assert.Nil(t, got.Location, "location: got %q", ptr.Deref(got.Location))`.
- **`t.Run` en vez de un parámetro `label`** cuando la misma comprobación corre sobre dos fuentes
  (lo que devolvió el `RETURNING` y lo que quedó en la tabla). El escenario pasa a ser el nombre del
  subtest, se puede correr suelto con `-run 'TestX/reread'`, y el `Error Trace` muestra las dos
  líneas: la assertion y el escenario que la llamó.

Una query solo para tests va como constante en el repositorio (`select<X>IncludingDeleted`) y se lee
con el mismo `scan<X>`, no como SQL suelto dentro del `_test.go`. La excepción son los `DELETE` de
`t.Cleanup`: no hay borrado duro en el repositorio a propósito —el proyecto entero se apoya en que
nada se borra de verdad— así que ese statement se queda como string en el test.

Tras cambiar una migración, actualizar a mano `services/<svc>/internal/db/models.go` (structs) y
`query.sql.go` (consts + `Scan`). No hay paso de generación.
