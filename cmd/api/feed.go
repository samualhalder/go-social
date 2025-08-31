package main

import (
	"net/http"

	"github.com/samualhalder/go-social/internal/store"
)

func (app *application) GetFeedForUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pagination := store.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
	}
	p, err := pagination.Parse(r)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := Validate.Struct(pagination); err != nil {
		app.badRequest(w, r, err)
		return
	}
	posts, err := app.store.Post.GetUserFeedPosts(ctx, 151, p)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, posts); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
