package main

import (
	"net/http"
)

// healthCheck godoc
// @Summary Health Check
// @Description Responds with status OK to show server is alive
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/health [get]
func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"message": "ok",
	}
	writeJSONError(w, http.StatusBadRequest, "cheking error")

	if err := writeJSON(w, http.StatusOK, data); err == nil {
		writeJSONError(w, http.StatusBadRequest, "cheking error")
	}
}
