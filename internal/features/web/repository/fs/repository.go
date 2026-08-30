// Package fs реализует репозиторий для чтения статических файлов из файловой системы.
package web_fs_repository

// WebRepository — репозиторий для доступа к файлам веб-интерфейса.
type WebRepository struct{}

// NewWebRepository создаёт репозиторий для файловой системы.
func NewWebRepository() *WebRepository {
	return &WebRepository{}
}
