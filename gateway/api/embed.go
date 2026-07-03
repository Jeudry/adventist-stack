// Package api expone la especificación OpenAPI y el HTML de Swagger UI
// embebidos, para que el gateway los sirva sin archivos externos.
package api

import _ "embed"

//go:embed openapi.yaml
var OpenAPISpec []byte

//go:embed swagger.html
var SwaggerHTML []byte
