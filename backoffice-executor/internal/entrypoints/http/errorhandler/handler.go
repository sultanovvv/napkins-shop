package errorhandler

import (
	"errors"
	"net/http"

	"napkins-shop/backoffice-executor/internal/customerr"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *Handler {
	return &Handler{logger: logger}
}

func Register(e *echo.Echo, handler *Handler) {
	e.HTTPErrorHandler = handler.Handle
}

// Handle вызывается Echo для ошибок обработчиков, маршрутизации и middleware.
// После начала ответа (например, скачивания файла) заменить его JSON уже нельзя.
func (h *Handler) Handle(err error, c echo.Context) {
	if err == nil || c.Response().Committed {
		return
	}

	status, body := classify(err)
	if status >= http.StatusInternalServerError {
		h.logger.Error("http request failed",
			zap.Error(err),
			zap.Int("status", status),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().URL.Path),
			zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
		)
	}

	var writeErr error
	if c.Request().Method == http.MethodHead {
		writeErr = c.NoContent(status)
	} else {
		writeErr = c.JSON(status, body)
	}
	if writeErr != nil {
		h.logger.Error("write http error response", zap.Error(writeErr))
	}
}

func classify(err error) (int, Response) {
	var appErr *customerr.Error
	if errors.As(err, &appErr) {
		status := http.StatusInternalServerError
		switch appErr.Code {
		case customerr.Validation:
			status = http.StatusBadRequest
		case customerr.NotFound:
			status = http.StatusNotFound
		case customerr.Conflict:
			status = http.StatusConflict
		case customerr.Forbidden:
			status = http.StatusForbidden
		case customerr.Unavailable:
			status = http.StatusServiceUnavailable
		case customerr.Internal:
		default:
			return http.StatusInternalServerError, response(string(customerr.Internal), "Internal Server Error")
		}
		details := make([]Detail, 0, len(appErr.Details))
		for _, detail := range appErr.Details {
			details = append(details, Detail{Code: detail.Code, Message: detail.Message, Target: detail.Target})
		}
		return status, response(string(appErr.Code), appErr.Message, details...)
	}
	return http.StatusInternalServerError, response(string(customerr.Internal), "Internal Server Error")
}
