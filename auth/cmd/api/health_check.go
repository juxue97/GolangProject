package main

import (
	"net/http"

	"github.com/juxue97/common"
)

// healthcheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	To perform server health check
//	@Tags			Health Check
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Failure		500	{object}	string	"error"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "ok",
		"environment": app.config.env,
		"version":     app.config.version,
	}
	common.WriteJSON(w, http.StatusOK, data)
}
