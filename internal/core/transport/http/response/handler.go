package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

type HTTPResponseHandler struct {
	log *slog.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(
	log *slog.Logger,
	rw http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

func (h *HTTPResponseHandler) HTMLResponse(html []byte) {
	h.rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.rw.WriteHeader(http.StatusOK)
	if _, err := h.rw.Write(html); err != nil {
		h.log.Error("write HTML HTTP response", "error", err)
	}
}

func (h *HTTPResponseHandler) JsonResponse(responseBody any, statusCode int) {
	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", "error", err)
	}
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...any)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}
	logFunc(msg, "error", err)

	h.errorResponse(statusCode, err, msg)
}

// method for sending http response in case of panic
func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)
	h.log.Error(msg, "error", err)

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
) {

	response := ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	}

	h.JsonResponse(response, statusCode)
}
