// Package testdb wires integration tests to a real PostgreSQL instance.
//
// Tests using it are guarded by the `integration` build tag, so `go test ./...`
// stays offline. Run them against the dev database with:
//
//	make infra-up
//	go run ./scripts/migrateall
//	go test -tags=integration ./...
package testdb

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Jeudry/adventist-stack/pkg/database"
)

// DefaultDSN points at the dev database from deploy/docker-compose.yml.
const DefaultDSN = "postgres://postgres:postgres@localhost:5433/adventist_stack?sslmode=disable"

// DSN returns DATABASE_URL when set, otherwise DefaultDSN.
func DSN() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}
	return DefaultDSN
}

// Connect opens a pool for the test and closes it on cleanup. It fails the test
// when the database is unreachable rather than skipping: an integration run that
// silently passes with nothing exercised is worse than a red one.
func Connect(t *testing.T) *pgxpool.Pool {
	t.Helper()

	pool, err := database.Connect(context.Background(), DSN())
	if err != nil {
		t.Fatalf("testdb: connect to %s: %v\n"+
			"is the dev database up? (make infra-up && go run ./scripts/migrateall)", DSN(), err)
	}
	t.Cleanup(pool.Close)
	return pool
}

// Exec runs a statement, reporting failures without aborting the test — it is
// meant for cleanup, where a failure should be visible but not mask the result.
func Exec(t *testing.T, pool *pgxpool.Pool, sql string, args ...any) {
	t.Helper()

	if _, err := pool.Exec(context.Background(), sql, args...); err != nil {
		t.Errorf("testdb: exec %q: %v", sql, err)
	}
}
