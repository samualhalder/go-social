package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samualhalder/go-social/internal/store"
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
			return
		}
		claims := jwtToken.Claims.(jwt.MapClaims)
		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.AuthorizationError(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.getUser(ctx, userId)
		if err != nil {
			app.AuthorizationError(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) checkPostOwnerShip(roleName string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r)
		post := getPostFromContext(r)
		// same user
		if user.Id == post.UserId {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		flag, err := app.checkRolePrecedence(ctx, user, roleName)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		if !flag {
			app.forbiddenError(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Role.GetByName(ctx, roleName)

	return user.Role.Level >= role.Level, err

}

func (app *application) getUser(ctx context.Context, userId int64) (*store.User, error) {

	user, err := app.cacheStorage.User.Get(ctx, userId)

	if user == nil {

		user, err = app.store.User.GetById(ctx, userId)

		if err != nil {
			return nil, err
		}
		if err := app.cacheStorage.User.Set(ctx, user); err != nil {
			return nil, err
		}

		return user, nil
	}
	return user, nil
}

func (app *application) RateTimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.ratelimiter.Enabled {
			if allowed, timeAfter := app.ratelimiter.Allow(r.RemoteAddr); !allowed {
				app.rateLimitExcedError(w, r, timeAfter.String())
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
