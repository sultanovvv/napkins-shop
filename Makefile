apis: public-openapi-bundle public-openapi-gen

# Собирает openapi.yaml (с возможными $ref) в плоский openapi.gen.yaml,
# который потом embed'ится в api_spec.go.
public-openapi-bundle:
	cd public-executor/api && npx -y @redocly/cli@latest bundle openapi.yaml -o openapi.gen.yaml

# Генерит api.gen.go из bundled openapi.gen.yaml.
public-openapi-gen:
	cd public-executor/api && go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config/generator.config.yaml openapi.gen.yaml
