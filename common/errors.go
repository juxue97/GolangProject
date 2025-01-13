package common

import (
	"errors"
	"net/http"
)

var (
	ErrApiKey             = errors.New("api key is required")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNotFound           = errors.New("user not found")
	ErrExceededLimit      = errors.New("rate limit exceeded")
	ErrContextNotFound    = errors.New("target user not found in context")
)

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	Logger.Warnf("Bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusBadRequest, err.Error())
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	Logger.Warnf("Internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	Logger.Warnf("Resource not found", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusNotFound, err.Error())
}

func UnauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	Logger.Warnf("Invalid credentials", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusUnauthorized, err.Error())
}

func UnauthorizedMiddlewareError(w http.ResponseWriter, r *http.Request, err error) {
	Logger.Warnf("Unauthorized request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusUnauthorized, err.Error())
}

func ForbiddenError(w http.ResponseWriter, r *http.Request, err error) {
	Logger.Warnf("Forbidden action", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusForbidden, err.Error())
}

func TooManyRequestsError(w http.ResponseWriter, r *http.Request, retryAfter string) {
	Logger.Warnf("Rate limit exceeded", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Retry-After", retryAfter)

	WriteError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter+"s")
}
