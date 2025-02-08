package common

import (
	"errors"
	"net/http"
	"os"

	"go.uber.org/zap"
)

var (
	ErrApiKey             = errors.New("api key is required")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNotFound           = errors.New("user not found")
	ErrExceededLimit      = errors.New("rate limit exceeded")
	ErrContextNotFound    = errors.New("target user not found in context")
	ErrNoQuantity         = errors.New("quantity cannot less than 1")
	ErrConvertID          = errors.New("failed to convert inserted ID to primitive.ObjectID")
	ErrNoDoc              = errors.New("no document found")
)

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	var logger *zap.SugaredLogger = NewLogger(os.Getenv("ENV"))

	logger.Warnf("Bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusBadRequest, err.Error())
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	var logger *zap.SugaredLogger = NewLogger(os.Getenv("ENV"))

	logger.Warnf("Internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func DuplicateErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var logger *zap.SugaredLogger = NewLogger(os.Getenv("ENV"))

	logger.Warnf("Duplicate error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusConflict, err.Error())
}

func NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	var logger *zap.SugaredLogger = NewLogger(os.Getenv("ENV"))

	logger.Warnf("Resource not found", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusNotFound, err.Error())
}

func UnauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	var logger *zap.SugaredLogger = NewLogger(os.Getenv("ENV"))

	logger.Warnf("Invalid credentials", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusUnauthorized, err.Error())
}

func UnauthorizedMiddlewareError(w http.ResponseWriter, r *http.Request, err error) {
	var logger *zap.SugaredLogger = NewLogger(os.Getenv("ENV"))

	logger.Warnf("Unauthorized request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusUnauthorized, err.Error())
}

func ForbiddenError(w http.ResponseWriter, r *http.Request, err error) {
	var logger *zap.SugaredLogger = NewLogger(os.Getenv("ENV"))

	logger.Warnf("Forbidden action", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusForbidden, err.Error())
}

func TooManyRequestsError(w http.ResponseWriter, r *http.Request, retryAfter string) {
	var logger *zap.SugaredLogger = NewLogger(os.Getenv("ENV"))

	logger.Warnf("Rate limit exceeded", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Retry-After", retryAfter)

	WriteError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter+"s")
}

func UnprocessableEntityResponse(w http.ResponseWriter, r *http.Request, err error) {
	var logger *zap.SugaredLogger = NewLogger(os.Getenv("ENV"))

	logger.Warnf("Unprocessable Entity", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusUnprocessableEntity, err.Error())
}
