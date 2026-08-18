package sabbath_school

import "embed"

//go:embed migrations/*.sql
var MigrationsFS embed.FS
