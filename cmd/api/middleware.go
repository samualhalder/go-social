package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) BaiscAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.AuthorizationError(w, r, fmt.Errorf("authorization token is missing"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.AuthorizationError(w, r, fmt.Errorf("authorization token form is not correct"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.AuthorizationError(w, r, err)
				return
			}

			creds := strings.SplitN(string(decoded), ":", 2)

			username := app.config.auth.basic.username
			password := app.config.auth.basic.pass
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.AuthorizationError(w, r, fmt.Errorf("invalid credentials"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.AuthorizationError(w, r, fmt.Errorf("authorization token is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.AuthorizationError(w, r, fmt.Errorf("authorization token form is not correct"))
			return
		}

		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.AuthorizationError(w, r, err)
		}
		claims := jwtToken.Claims.(jwt.MapClaims)
		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.AuthorizationError(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.store.User.GetById(ctx, userId)
		if err != nil {
			app.AuthorizationError(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
