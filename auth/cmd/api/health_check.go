package main

import (
	"net/http"

	"github.com/juxue97/common"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "ok",
		"environment": app.config.env,
		"version":     app.config.version,
	}

	common.WriteJSON(w, http.StatusOK, data)
}
