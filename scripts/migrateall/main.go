// Command migrateall applies every service's embedded migrations against one
// database, using the same code path the servers run on startup.
package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/Jeudry/adventist-stack/pkg/database"
	"github.com/Jeudry/adventist-stack/services/auth"
	"github.com/Jeudry/adventist-stack/services/members"
	"github.com/Jeudry/adventist-stack/services/prayers"
	"github.com/Jeudry/adventist-stack/services/sabbath_school"
)

type service struct {
	name  string
	table string
	fs    embed.FS
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5433/adventist_stack?sslmode=disable"
	}

	services := []service{
		{"auth", "auth_schema_migrations", auth.MigrationsFS},
		{"members", "members_schema_migrations", members.MigrationsFS},
		{"prayers", "prayers_schema_migrations", prayers.MigrationsFS},
		{"sabbath_school", "sabbath_school_schema_migrations", sabbath_school.MigrationsFS},
	}

	failed := false
	for _, s := range services {
		if err := database.Migrate(dsn, s.table, s.fs, "migrations"); err != nil {
			fmt.Printf("%-16s FAIL  %v\n", s.name, err)
			failed = true
			continue
		}
		fmt.Printf("%-16s ok\n", s.name)
	}
	if failed {
		os.Exit(1)
	}
}
