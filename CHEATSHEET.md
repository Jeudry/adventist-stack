# Cheat Sheet — Comandos Rápidos

## 🚀 Infraestructura & Entorno

```bash
make infra-up                         # Levanta Postgres, Redis y MailHog (Docker)
make infra-down                       # Detiene la infraestructura
```

---

## 🛠️ Generación de Código

```bash
make proto                            # Genera Go desde archivos .proto
make tidy                             # Sincroniza dependencias (go mod tidy)
```

---

## 🗄️ Migraciones de Base de Datos

```bash
# Crear migración vacía (.up.sql / .down.sql)
make migration svc=sabbath_school name=create_sabbath_school_table

# Aplicar migraciones
make migrate-up svc=sabbath_school

# Revertir última migración
make migrate-down svc=sabbath_school n=1

# Ver versión actual de migración
make migrate-version svc=sabbath_school

# Desbloquear migración 'dirty' tras un fallo
make migrate-force svc=sabbath_school version=1
```

---

## 🧪 Compilación & Tests

```bash
go build ./...                        # Compila todo el monorepo
go test ./...                         # Ejecuta todos los unit tests
go vet ./...                          # Análisis estático de código
```

---

## 🟢 Ejecución de Servicios

```bash
go run ./services/auth/cmd/server           # Auth (:50051)
go run ./services/members/cmd/server        # Members (:50052)
go run ./services/prayers/cmd/server        # Prayers (:50055)
go run ./services/sabbath_school/cmd/server # Sabbath School (:50056)
make run-gateway                            # REST Gateway (:8080)
```
