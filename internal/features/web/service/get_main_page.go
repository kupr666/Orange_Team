package web_service

import (
	"fmt"
	"path"

	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

func (s *WebService) GetMainPage() (core_http_response.File, error) {
	htmlFilePath := path.Join(s.projectRoot, "frontend/dist/index.html")
	htmlFile, err := s.webRepository.GetFile(htmlFilePath)
	if err != nil {
		return core_http_response.File{}, fmt.Errorf("get file from repository: %w", err)
	}
	return htmlFile, nil
}
