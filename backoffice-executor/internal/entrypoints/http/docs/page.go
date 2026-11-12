package docs

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"napkins-shop/backoffice-executor/api"
	httpentry "napkins-shop/backoffice-executor/internal/entrypoints/http"
)

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Backoffice API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      SwaggerUIBundle({ url: './openapi.yaml', dom_id: '#swagger-ui', deepLinking: true });
    };
  </script>
</body>
</html>`

func Register(e *echo.Echo) {
	e.GET(httpentry.BasePath+"/openapi.yaml", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/yaml", api.Spec)
	})
	e.GET(httpentry.BasePath+"/swagger", func(c echo.Context) error {
		return c.HTML(http.StatusOK, swaggerHTML)
	})
}
