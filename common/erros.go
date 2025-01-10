package common

import (
	"net/http"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func init() {
	Logger = newLogger()
}

func newLogger() *zap.SugaredLogger {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	return logger
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	Logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusBadRequest, err.Error())
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	Logger.Warnf("Internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteError(w, http.StatusInternalServerError, "the server encountered a problem")
}
