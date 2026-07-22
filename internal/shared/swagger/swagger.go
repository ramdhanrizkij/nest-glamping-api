package swagger

import (
	"io/fs"
	"strings"

	"github.com/gofiber/fiber/v3"
	swaggerFiles "github.com/swaggo/files/v2"
)

func Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		path := c.Path()

		if path == "/swagger" || path == "/swagger/" {
			return c.Redirect().Status(fiber.StatusMovedPermanently).To("/swagger/index.html")
		}

		stripPath := strings.TrimPrefix(path, "/swagger/")
		if stripPath == "" {
			stripPath = "index.html"
		}

		f, err := fs.ReadFile(swaggerFiles.FS, stripPath)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("not found")
		}

		contentType := "text/plain"
		switch {
		case strings.HasSuffix(stripPath, ".html"):
			contentType = "text/html"
		case strings.HasSuffix(stripPath, ".js"):
			contentType = "application/javascript"
		case strings.HasSuffix(stripPath, ".css"):
			contentType = "text/css"
		case strings.HasSuffix(stripPath, ".json"):
			contentType = "application/json"
		case strings.HasSuffix(stripPath, ".png"):
			contentType = "image/png"
		}

		c.Set("Content-Type", contentType)
		return c.Send(f)
	}
}

func SwaggerUI(title, specURL string) fiber.Handler {
	return func(c fiber.Ctx) error {
		html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>` + title + ` - Swagger UI</title>
  <link rel="stylesheet" href="/swagger/swagger-ui.css">
  <style>
    body { margin: 0; padding: 0; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="/swagger/swagger-ui-bundle.js"></script>
  <script src="/swagger/swagger-ui-standalone-preset.js"></script>
  <script>
    SwaggerUIBundle({
      url: "` + specURL + `",
      dom_id: '#swagger-ui',
      deepLinking: true,
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
      layout: "BaseLayout"
    });
  </script>
</body>
</html>`

		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	}
}
