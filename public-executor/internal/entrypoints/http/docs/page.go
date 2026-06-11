package docs

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"napkins-shop/public-executor/api"
)

// swaggerHTML — минимальная страница swagger-ui-dist с CDN. Грузит spec
// по /openapi.yaml, который мы отдаём из embed'нутого api.Spec.
const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Decoupage Napkins Shop API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '/openapi.yaml',
        dom_id: '#swagger-ui',
        deepLinking: true,
      });
    };
  </script>
</body>
</html>`

func Register(e *echo.Echo) {
	e.GET("/openapi.yaml", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/yaml", api.Spec)
	})
	e.GET("/swagger", func(c echo.Context) error {
		return c.HTML(http.StatusOK, swaggerHTML)
	})
	// Заглушка, чтобы браузер не спамил 404 в логах при открытии /swagger.
	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
}
