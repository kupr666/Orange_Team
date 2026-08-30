package web_transport_http_spa

import (
	"net/http"
	"strings"

	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
)

// RegisterSPA регистрирует обработчик для SPA (одностраничного приложения).
// Все запросы, кроме /api/v1/*, отдают статику из папки distPath,
// а если файл не найден — index.html (SPA-fallback).
func RegisterSPA(server *core_http_server.HTTPServer, distPath string) {
	staticFS := http.FileServer(http.Dir(distPath))

	spaHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Если запрос к API – не трогаем (они уже обработаны выше)
		if strings.HasPrefix(r.URL.Path, "/api/v1/") {
			http.NotFound(w, r)
			return
		}

		// Если запрос к статике – отдаём файл
		if strings.HasPrefix(r.URL.Path, "/static/") ||
			strings.HasPrefix(r.URL.Path, "/assets/") ||
			strings.Contains(r.URL.Path, ".") {
			staticFS.ServeHTTP(w, r)
			return
		}

		// Все остальные запросы → index.html (SPA)
		http.ServeFile(w, r, distPath+"/index.html")
	})

	server.RegisterRoutes(core_http_server.Route{
		Method:  http.MethodGet,
		Path:    "/",
		Handler: spaHandler,
	})
}
