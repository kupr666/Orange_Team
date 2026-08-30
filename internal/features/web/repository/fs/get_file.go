package web_fs_repository

import (
	"errors"
	"fmt"
	"os"

	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

func (r *WebRepository) GetFile(filePath string) (core_http_response.File, error) {
	buffer, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return core_http_response.File{}, fmt.Errorf("file: %s: %w", filePath, core_errors.ErrNotFound)
		}
		return core_http_response.File{}, fmt.Errorf("get file: %s: %w", filePath, err)
	}
	return core_http_response.NewFile(buffer), nil
}
