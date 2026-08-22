package core_http_response

import (
	"net/http"
)

// StatusCodeUninitialized — сигнальное значение: WriteHeader ещё не вызывался.
// Если статус не был установлен явно, Go возвращает 200 — мы имитируем это поведение.
var (
	StatusCodeUninitialized = -1
)

// ResponseWriter — обёртка над http.ResponseWriter, которая запоминает статус-код.
// Стандартный http.ResponseWriter не даёт прочитать статус после WriteHeader,
// а он нужен middleware Trace для логирования.
//
// Встраивание http.ResponseWriter даёт все методы «бесплатно»;
// переопределяем только WriteHeader, чтобы перехватить код.
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter создаёт обёртку с незаписанным статусом.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeUninitialized,
	}
}

// WriteHeader перехватывает статус-код и передаёт его дальше в оригинальный ResponseWriter.
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

// GetStatusCode возвращает записанный статус-код.
// Если WriteHeader не вызывался — возвращает 200 (поведение по умолчанию в Go).
func (rw *ResponseWriter) GetStatusCode() int {
	if rw.statusCode == StatusCodeUninitialized {
		return http.StatusOK
	}

	return rw.statusCode
}
