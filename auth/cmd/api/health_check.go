package main

import (
	"net/http"

	common "github.com/GolangProject/common"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "ok",
		"environment": app.config.env,
		"version":     "0.0.1",
	}

	common.WriteJSON(w, http.StatusOK, data)
}
