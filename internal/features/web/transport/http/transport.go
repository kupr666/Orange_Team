package web_transport_http

import (
	"net/http"

	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
)

type WebHTTPHandler struct {
	webService WebService
}

type WebService interface {
	GetMainPage() (core_http_response.File, error)
}

func NewWebHTTPHandler(webService WebService) *WebHTTPHandler {
	return &WebHTTPHandler{webService: webService}
}

func (h *WebHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: h.GetMainPage,
		},
	}
}
