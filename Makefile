# При явно заданном GOROOT используем go из того же SDK, что и компилятор.
# Выбор можно переопределить: make GO=/path/to/go test.
GO ?= $(if $(GOROOT),$(GOROOT)/bin/go,go)

.PHONY: apis public-openapi-bundle public-openapi-gen backoffice-openapi-gen test build backoffice-run go-version

apis: public-openapi-bundle public-openapi-gen backoffice-openapi-gen

# Собирает openapi.yaml (с возможными $ref) в плоский openapi.gen.yaml,
# который потом embed'ится в api_spec.go.
public-openapi-bundle:
	cd public-executor/api && npx -y @redocly/cli@latest bundle openapi.yaml -o openapi.gen.yaml

# Генерит api.gen.go из bundled openapi.gen.yaml.
public-openapi-gen:
	cd public-executor/api && "$(GO)" run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config/generator.config.yaml openapi.gen.yaml

backoffice-openapi-gen:
	cd backoffice-executor/api && "$(GO)" generate .

test:
	"$(GO)" test ./shared/... ./public-executor/... ./backoffice-executor/...

build:
	"$(GO)" build ./shared/... ./public-executor/... ./backoffice-executor/...

backoffice-run:
	"$(GO)" run ./backoffice-executor/cmd

go-version:
	"$(GO)" version
	"$(GO)" tool compile -V
