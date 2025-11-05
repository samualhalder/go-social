package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samualhalder/go-social/internal/store"
)

type UserType string

var userCtx UserType = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	user, err := app.store.User.GetById(ctx, id)
	if err != nil {
		switch err {
		case store.ErrorNotFound:
			app.notFound(w, r, err)
		default:
			app.badRequest(w, r, err)
		}
	}
	if err := writeJSON(w, http.StatusOK, user); err != nil {
		app.badRequest(w, r, err)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {

	user := getUserFromContext(r)
	followedId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	ctx := r.Context()

	if err := app.store.Follower.Follow(ctx, followedId, user.Id); err != nil {
		app.badRequest(w, r, err)
		return
	}
	app.jsonResponse(w, http.StatusOK, "followed")
}

func (app *application) unFollowUserHandler(w http.ResponseWriter, r *http.Request) {

	user := getUserFromContext(r)
	ctx := r.Context()
	followedId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := app.store.Follower.UnFollow(ctx, followedId, user.Id); err != nil {
		switch err {
		case store.ErrConflict:
			app.ConflictError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
		}
		app.badRequest(w, r, err)
		return
	}
	app.jsonResponse(w, http.StatusOK, "unfollowed")
}

func (app *application) activateUserHanlder(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	ctx := r.Context()
	err := app.store.User.Activate(ctx, token)
	if err != nil {
		switch err {
		case store.ErrorNotFound:
			app.badRequest(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, ""); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
func (app *application) getUserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}
		user, err := app.store.User.GetById(ctx, id)
		if err != nil {
			switch err {
			case store.ErrorNotFound:
				app.notFound(w, r, err)
			default:
				app.badRequest(w, r, err)
			}
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
