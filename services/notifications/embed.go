package notifications

import "embed"

// TemplatesFS contiene las plantillas de correo embebidas.
//
//go:embed templates/*.html
var TemplatesFS embed.FS
