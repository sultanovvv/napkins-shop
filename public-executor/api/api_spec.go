//go:generate oapi-codegen --config=config/generator.config.yaml openapi.gen.yaml

package api

import _ "embed"

//go:embed openapi.gen.yaml
var Spec []byte
