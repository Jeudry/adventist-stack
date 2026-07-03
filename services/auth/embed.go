package auth

import "embed"

// MigrationsFS contiene las migraciones SQL embebidas del servicio auth.
//
//go:embed migrations/*.sql
var MigrationsFS embed.FS
