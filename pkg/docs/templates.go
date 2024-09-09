package docs

// templates.go provides templates load at compilation time.

import _ "embed"

//go:embed templates/openapi.yaml
var docs []byte

//go:embed templates/elements.tpl
var elements string
