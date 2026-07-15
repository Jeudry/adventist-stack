// Package database maneja la conexión a PostgreSQL (pool con pgx) y la
// ejecución de migraciones embebidas con golang-migrate.
package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // driver database/sql para migraciones
	pgxuuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

// Connect abre un pool de conexiones a PostgreSQL y verifica conectividad.
func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("database: parse dsn: %w", err)
	}
	cfg.MaxConns = 10
	cfg.MaxConnLifetime = time.Hour

	// Registra el tipo google/uuid en cada conexión para que las columnas
	// `uuid` de Postgres se lean/escriban como uuid.UUID de forma nativa.
	cfg.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("database: connect: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database: ping: %w", err)
	}
	return pool, nil
}

// Migrate aplica las migraciones embebidas (carpeta migrations/ de cada
// servicio). Es idempotente: si no hay cambios, no hace nada.
//
// table es la tabla donde se rastrea la versión. Cada servicio usa la suya
// (ej. "auth_schema_migrations") para que sus versiones no choquen con las de
// otros servicios que comparten la misma base de datos.
func Migrate(dsn, table string, fsys embed.FS, dir string) error {
	src, err := iofs.New(fsys, dir)
	if err != nil {
		return fmt.Errorf("database: migrate source: %w", err)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("database: migrate open: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{MigrationsTable: table})
	if err != nil {
		return fmt.Errorf("database: migrate driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("database: migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("database: migrate up: %w", err)
	}
	return nil
}
