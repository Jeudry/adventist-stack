package api

import _ "embed"

// The API document is not a file: it is built at runtime from what each service
// declares about itself, and served at /openapi.json. Only the viewer is static.
//
//go:embed swagger.html
var SwaggerHTML []byte
