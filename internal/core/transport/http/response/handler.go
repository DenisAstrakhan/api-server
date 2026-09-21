package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/DenisAstrakhan/api-server/internal/core/errors"
	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *core_logger.Logge
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(log *core_logger.Logge, rw http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}

}
func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.WriteHeader(statusCode)
	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP respons", zap.Error(err))
	}
}

// Метод для отправки HTTP ответа в случаи ошибки
func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
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
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}
	logFunc(msg, zap.Error(err))
	h.errorResponse(statusCode, err, msg)
}

// Метод для отправки HTTP ответа в случаи возникновения паники
func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	err := fmt.Errorf("unexpected panic: %v", p)
	h.log.Error(msg, zap.Error(err))

	h.errorResponse(http.StatusInternalServerError, err, msg)

}

func (h *HTTPResponseHandler) errorResponse(statusCode int, err error, msg string) {
	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}
	h.JSONResponse(response, statusCode)

}
