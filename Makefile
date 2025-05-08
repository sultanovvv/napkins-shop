apis:public-openapi-gen

public-openapi-gen:
	cd public-executor/internal/api &&  go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config/generator.config.yaml openapi.yaml