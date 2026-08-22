// Package core_errors содержит sentinel-ошибки, общие для всего приложения.
//
// Sentinel-ошибки (переменные типа error) используются для проверки через
// errors.Is(). Это позволяет слоям приложения принимать решение о
// своих дальнейших действиях на основе полученной ошибки,
// не зная деталей реализации низлежащих слоёв.
//
// Пример: если репозиторий вернул errors.Is(err, ErrNotFound) == true,
// HTTP-обработчик ответит статусом 404.
package core_errors

import "errors"

var (
	// ErrNotFound — запрашиваемая сущность не найдена (HTTP 404).
	ErrNotFound = errors.New("not found")

	// ErrInvalidArgument — переданы некорректные данные (HTTP 400).
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrConflict — конфликт при обновлении, обычно из-за конкурентного
	// изменения той же записи (HTTP 409).
	ErrConflict = errors.New("conflict")
)
