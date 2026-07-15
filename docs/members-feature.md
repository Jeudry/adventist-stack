# Feature: Gestión de Miembros (`services/members`)

Guía paso a paso para construir la primera feature vertical del proyecto.
El objetivo no es solo tener miembros funcionando, sino **aprender el patrón**
que vas a repetir en cada feature: `.proto → repositorio → servicio → gRPC →
gateway REST → prueba`.

> 💡 **Referencia viva:** el servicio `auth/` ya sigue exactamente esta estructura.
> Cuando dudes de "cómo se hace X", abrí el archivo equivalente en `services/auth/`.

- **Rama:** `feature/1-Members_Management`
- **Definición de terminado:** poder hacer CRUD de miembros por REST a través del
  gateway (`/api/v1/members`), protegido con JWT, y probado con Swagger/curl.

---

## Paso 0 — Preparar el entorno

Antes de escribir código, tené la infraestructura arriba.

```bash
cd ~/Documents/Coding/adventist-stack
git branch --show-current          # confirmá: feature/1-Members_Management
make infra-up                      # Postgres, Redis, MailHog
```

**Checkpoint:** `docker ps` muestra los 3 contenedores `adventist_*` healthy.

---

## Paso 1 — Diseñar el contrato gRPC (`.proto`)

El `.proto` es la **fuente de verdad** de la API. Todo lo demás se deriva de él.

1. Crear `proto/members/v1/members.proto`.
2. Definir el mensaje `Member` y el servicio `MemberService` con CRUD.

**Decisiones a entender (no copiar a ciegas):**
- **`Create` vs `Update` con mensajes distintos:** el request de crear NO lleva
  `id` (lo genera la DB); el de actualizar SÍ. Mensajes separados = contratos claros.
- **`List` con paginación:** nunca devuelvas "todos" sin límite. Campos típicos:
  `page`, `page_size` en el request; `items` + `total` en el response.
- **Campos opcionales:** en proto3, los escalares tienen valor cero por defecto.
  Para "email no enviado" vs "email vacío", usá `optional` (genera un puntero en Go).
- **Fechas:** usamos `google.protobuf.Timestamp` (tipo estándar de Google, fecha/hora
  UTC). En Go se convierte con `timestamppb.New(t)` (desde `time.Time`) y `.AsTime()`
  (hacia `time.Time`).

**Campos propuestos para `Member`** (los afinamos juntos):
`id, first_name, last_name, email, phone, birth_date, gender, address,
baptism_date, status (active|inactive|visitor), created_at, updated_at`.

**RPCs:** `CreateMember, GetMember, ListMembers, UpdateMember, DeleteMember`.

3. Generar el código Go:

```bash
make proto
```

**Checkpoint:** aparecen `gen/members/v1/members.pb.go` y `members_grpc.pb.go`.

---

## Paso 2 — Migración + modelos (core y DB)

Estructura por capas (mirá `services/auth/internal/`):

1. **Migración** `services/members/migrations/0001_init.up.sql` (+ `.down.sql`)
   - Tabla `members` con las columnas del paso 1.
   - Índices donde tenga sentido (ej. `email`).
2. **`embed.go`** en `services/members/` que embebe `migrations/*.sql`
   (copiá el patrón de `services/auth/embed.go`).
3. **Modelo de dominio (core)** `internal/domain/member.go`
   - El `Member` "puro", sin tags de DB ni JSON. Acá viven también los errores
     de dominio (`ErrMemberNotFound`, etc.) y tipos como `Status`.
4. **Modelo de DB** `internal/models/member.go`
   - Mapea las columnas de la tabla. Método `ToDomain()` para convertir.

**Por qué dos modelos:** el dominio no debe "saber" de columnas ni de JSON. Esto
te deja cambiar la DB sin tocar la lógica, y testear la lógica sin DB.

**Checkpoint:** `go build ./services/members/...` compila (aunque falten capas,
estos paquetes deben compilar solos).

---

## Paso 3 — Repositorio (acceso a datos)

`internal/repository/member_repository.go` — mirá `auth`'s `user_repository.go`.

- Recibe un `*pgxpool.Pool`.
- Métodos: `Create, GetByID, List (con limit/offset), Update, Delete, Count`.
- Devuelve **modelos de dominio**, no de DB (convertí con `ToDomain()`).
- Mapeá `pgx.ErrNoRows` → `domain.ErrMemberNotFound`.

**Checkpoint:** compila. (Todavía no hay cómo llamarlo; eso viene ahora.)

---

## Paso 4 — Servicio (lógica de negocio)

`internal/service/member_service.go` — mirá `auth`'s `auth_service.go`.

- Define una **interface** `memberRepository` (para poder testear con un mock).
- Validaciones: nombre obligatorio, email con formato válido si viene, `status`
  dentro de los permitidos, `page_size` con tope (ej. máx 100).
- **Normalización:** trim de strings, email a minúsculas (igual que en auth).
- Reglas del dominio viven acá, NO en el repositorio ni en el handler.

**Checkpoint:** compila. Si querés, escribí un test unitario del servicio con un
repo falso (ver skill/ejemplo de TDD) — es el mejor momento porque no depende de DB.

---

## Paso 5 — Servidor gRPC (transporte)

`internal/grpc/server.go` — mirá `auth`'s `internal/grpc/server.go`.

- Implementa `membersv1.MemberServiceServer` (embebé `UnimplementedMemberServiceServer`).
- Es un **adaptador delgado:** traduce request proto → llamada al servicio →
  response proto. Sin lógica de negocio acá.
- Traducí errores de dominio a códigos gRPC (`NotFound`, `InvalidArgument`, …)
  con una función `toStatus()` como la de auth.

---

## Paso 6 — `main.go` del servicio

`services/members/cmd/server/main.go` — casi idéntico a `auth`'s main.

- Cargar config (`ENV`, `MEMBERS_GRPC_PORT` ej. `50053`, `Postgres`).
- Correr migraciones si `AUTO_MIGRATE`.
- Conectar Postgres → repo → servicio → servidor gRPC.
- Registrar el server + reflection, escuchar, apagado graceful.

Agregar a `.env.example` y `.env`:
```
MEMBERS_GRPC_PORT=50053
MEMBERS_GRPC_ADDR=localhost:50053
```
Y un target en el `Makefile`: `run-members`.

**Checkpoint (servicio aislado):**
```bash
make run-members
# en otra terminal, probá con grpcurl:
grpcurl -plaintext -d '{"first_name":"Juan","last_name":"Pérez","status":"active"}' \
  localhost:50053 members.v1.MemberService/CreateMember
```

---

## Paso 7 — Exponerlo por REST en el gateway

1. **Cliente gRPC:** en `gateway/internal/clients/clients.go`, agregá el
   `MemberServiceClient` y su conexión (junto a auth y notifications).
2. **Handler REST:** `gateway/internal/handlers/member_handler.go`
   - Viewmodels (DTOs) para request/response (mirá `auth_handler.go`).
   - Métodos: `Create, Get, List, Update, Delete`.
3. **Rutas:** en `gateway/internal/router/router.go`, dentro de `/api/v1`,
   **protegidas con JWT** (grupo con `middleware.Auth`):
   ```
   POST   /api/v1/members
   GET    /api/v1/members          (con ?page=&page_size=)
   GET    /api/v1/members/{id}
   PUT    /api/v1/members/{id}
   DELETE /api/v1/members/{id}
   ```
4. **Swagger:** agregá estos endpoints a `gateway/api/openapi.yaml`.

---

## Paso 8 — Prueba end-to-end

Con `make infra-up` y los servicios corriendo
(`run-auth`, `run-members`, `run-gateway`):

```bash
BASE=http://localhost:8090

# 1. Login para obtener token (auth ya existe)
TOKEN=$(curl -s -X POST $BASE/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"pastor@iglesia.org","password":"secreto123"}' \
  | python3 -c "import sys,json;print(json.load(sys.stdin)['access_token'])")

# 2. Crear miembro
curl -s -X POST $BASE/api/v1/members \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"first_name":"Ana","last_name":"García","status":"active"}'

# 3. Listar
curl -s "$BASE/api/v1/members?page=1&page_size=10" -H "Authorization: Bearer $TOKEN"
```

**Definición de terminado ✅**
- [ ] `make proto` genera el código de members
- [ ] `go build ./...` y `go vet ./...` limpios
- [ ] CRUD completo funciona por REST a través del gateway
- [ ] Endpoints protegidos con JWT (sin token → 401)
- [ ] Paginación funciona en `List`
- [ ] Endpoints documentados en Swagger
- [ ] (ideal) al menos un test unitario del servicio

---

## Orden mental para no perderte

```
proto  →  migración+modelos  →  repositorio  →  servicio  →  grpc  →  main
                                                                        │
                                       gateway (cliente + handler + rutas + swagger)
                                                                        │
                                                                     probar
```

Siempre de **adentro hacia afuera**: primero el dato y la lógica, al final la
exposición REST. Cada paso compila antes de avanzar al siguiente.

---

## Commits sugeridos (los hacés vos)

Podés hacer un commit por bloque lógico, por ejemplo:
1. `feat(members): proto + código generado`
2. `feat(members): migración, modelos y repositorio`
3. `feat(members): servicio con validaciones + tests`
4. `feat(members): servidor gRPC y main`
5. `feat(members): endpoints REST en el gateway + swagger`

Al terminar, PR desde `feature/1-Members_Management` hacia `main`.
```bash
gh pr create --base main --head feature/1-Members_Management --web
```
