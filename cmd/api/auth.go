package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"
	"github.com/samualhalder/go-social/internal/store"
)

type RegisterUserPayloadType struct {
	Username string `json:"username" validate:"required,max=24"`
	Email    string `json:"email" validate:"required,max=24"`
	Password string `json:"password" validate:"required,max=24"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var registerPayload RegisterUserPayloadType
	if err := readJSON(w, r, registerPayload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := Validate.Struct(registerPayload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	user := &store.User{
		Username: registerPayload.Username,
		Email:    registerPayload.Email,
	}
	if err := user.Password.Set(registerPayload.Password); err != nil {
		app.badRequest(w, r, err)
	}
	ctx := r.Context()

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken)) // not readble by human so cant store in sql
	hashToken := hex.EncodeToString(hash[:])  // readable by human a string

	if err := app.store.User.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
