package web_service

import core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"

type WebService struct {
	webRepository WebRepository
	projectRoot   string
}

type WebRepository interface {
	GetFile(filePath string) (core_http_response.File, error)
}

func NewWebService(webRepository WebRepository, projectRoot string) *WebService {
	return &WebService{
		webRepository: webRepository,
		projectRoot:   projectRoot,
	}
}
