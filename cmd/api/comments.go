package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samualhalder/go-social/internal/store"
)

type CommentPaylod struct {
	UserId  int64  `json:"user_id" validate:"required"`
	PostId  int64  `json:"post_id" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (app *application) createComment(w http.ResponseWriter, r *http.Request) {
	var commentPayload CommentPaylod
	postId := chi.URLParam(r, "postId")
	postIdint64, err := strconv.ParseInt(postId, 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	commentPayload.PostId = postIdint64
	if err := readJSON(w, r, &commentPayload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := Validate.Struct(commentPayload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	comment := &store.Comment{
		PostId:  commentPayload.PostId,
		UserId:  commentPayload.UserId,
		Content: commentPayload.Content,
	}
	if err := app.store.Comment.Create(r.Context(), comment); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusCreated, comment); err != nil {
		app.badRequest(w, r, err)
		return
	}
}
