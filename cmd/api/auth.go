package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samualhalder/go-social/internal/mailer"
	"github.com/samualhalder/go-social/internal/store"
)

type RegisterUserPayloadType struct {
	Username string `json:"username" validate:"required,max=24"`
	Email    string `json:"email" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=24"`
}
type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var registerPayload RegisterUserPayloadType
	if err := readJSON(w, r, &registerPayload); err != nil {
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

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	if err := app.store.User.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp); err != nil {
		app.badRequest(w, r, err)
		return
	}

	isProdEnv := app.config.env == "production"

	activationURL := fmt.Sprint("%s/confirm/%s", app.config.frontEndURL, plainToken)
	vars := struct {
		ActivationURL string
		Username      string
	}{
		ActivationURL: activationURL,
		Username:      user.Username,
	}

	// sending mail

	err := app.mailer.Send(mailer.UserRegisterMailTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {

		if err := app.store.User.Delete(ctx, user.Id); err != nil {
			app.logger.Errorw("Error Deleting", "user", err.Error())
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type tokenPayloadType struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var tokenPayload tokenPayloadType
	if err := readJSON(w, r, &tokenPayload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := Validate.Struct(tokenPayload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err := app.store.User.GetByEmail(r.Context(), tokenPayload.Email)
	if err != nil {
		switch err {
		case store.ErrorNotFound:
			app.AuthorizationError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := user.Password.Check(tokenPayload.Password); err != nil {
		app.forbiddenError(w, r)
		return
	}

	claims := jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(app.config.auth.token.expiry).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.issuer,
		"aud": app.config.auth.token.issuer,
	}
	token, err := app.authenticator.GenarateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, token); err != nil {
		app.internalServerError(w, r, err)
	}
}
