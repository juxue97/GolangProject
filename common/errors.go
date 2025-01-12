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
)

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	Logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

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
