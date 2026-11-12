//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate types,server -package api -o api.gen.go openapi.yaml

package api

import _ "embed"

//go:embed openapi.yaml
var Spec []byte
