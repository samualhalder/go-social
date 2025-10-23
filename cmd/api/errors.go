package main

import (
	"net/http"
)

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("Bad Request: %s,err: %s,path: %s", r.Method, r.URL.Path, err.Error())
	app.logger.Errorw("Bad Request", "Methode", r.Method, "Path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}
func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("Not Found: %s,err: %s,path: %s", r.Method, r.URL.Path, err.Error())
	app.logger.Errorw("Not Found", "Methode", r.Method, "Path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusNotFound, "Record not found")
}
func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("Internal Server Error: %s,err: %s,path: %s", r.Method, r.URL.Path, err.Error())
	app.logger.Errorw("Internal Server Error", "Methode", r.Method, "Path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusBadRequest, "Something went wrong")
}
func (app *application) ConflictError(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("Conflict Error: %s,err: %s,path: %s", r.Method, r.URL.Path, err.Error())
	app.logger.Errorw("Conflict Error", "Methode", r.Method, "Path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}
func (app *application) AuthorizationError(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("Conflict Error: %s,err: %s,path: %s", r.Method, r.URL.Path, err.Error())
	app.logger.Errorw("Authorization Error", "Methode", r.Method, "Path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusUnauthorized, err.Error())
}
