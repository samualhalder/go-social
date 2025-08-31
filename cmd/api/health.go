package main

import (
	"net/http"
)

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"message": "ok",
	}
	writeJSONError(w, http.StatusBadRequest, "cheking error")

	if err := writeJSON(w, http.StatusOK, data); err == nil {
		writeJSONError(w, http.StatusBadRequest, "cheking error")
	}
}
