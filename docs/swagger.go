package api_docs

import (
	_ "embed"
	"io"
	"net/http"

	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
)

//go:embed openapi.yaml
var openAPISpec []byte

const swaggerUIHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Orange Team API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui.css">
  <style>
    html { box-sizing: border-box; overflow-y: scroll; }
    *, *::before, *::after { box-sizing: inherit; }
    body { margin: 0; background: #fafafa; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: "/openapi.yaml",
        dom_id: "#swagger-ui",
        deepLinking: true,
        displayRequestDuration: true,
        persistAuthorization: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: "StandaloneLayout"
      });
    };
  </script>
</body>
</html>`

func Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/openapi.yaml",
			Handler: openAPIHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/swagger",
			Handler: swaggerUIHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/swagger/",
			Handler: swaggerUIHandler,
		},
	}
}

func openAPIHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(openAPISpec)
}

func swaggerUIHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, swaggerUIHTML)
}
