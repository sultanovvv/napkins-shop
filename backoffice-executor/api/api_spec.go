//go:generate npx -y @redocly/cli@2.53.3 bundle openapi.yaml -o openapi.gen.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate types,server -package api -o api.gen.go openapi.gen.yaml

package api

import _ "embed"

//go:embed openapi.gen.yaml
var Spec []byte
