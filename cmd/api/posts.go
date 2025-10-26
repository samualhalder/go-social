package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samualhalder/go-social/internal/store"
)

type PostKey string

var postCtx PostKey = "post"

type PostData struct {
	Title   string   `json:"title" validate:"required"`
	Content string   `json:"content" validate:"required"`
	Tags    []string `json:"tags" validate:"required"`
}

// createPost godoc
// @Summary Create a new post
// @Description Creates a new post for the authenticated user
// @Tags posts
// @Accept json
// @Produce json
// @Param post body map[string]string true "Post content"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/posts/create [post]
func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	var payload PostData
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	user := getUserFromContext(r)
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserId:  user.Id,
	}
	ctx := r.Context()

	if err := app.store.Post.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	comments, err := app.store.Comment.GetCommentByPostId(r.Context(), post.Id)
	if err != nil {
		app.internalServerError(w, r, err)
	}
	post.Comments = comments
	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postId := chi.URLParam(r, "postId")
	id64, err := strconv.ParseInt(postId, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.store.Post.DeletePostById(ctx, id64); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.notFound(w, r, fmt.Errorf("No row founded"))
			return
		default:
			app.internalServerError(w, r, err)
		}
	}
	if err := app.jsonResponse(w, http.StatusOK, "Post deleted successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty"`
	Content *string `json:"content" validate:"omitempty"`
}

func (app *application) updatePostById(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)
	var payload UpdatePostPayload
	fmt.Print("hit here")
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if err := app.store.Post.UpdatePostById(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		postId := chi.URLParam(r, "postId")
		postIdint, err := strconv.ParseInt(postId, 10, 64)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}
		post, err := app.store.Post.GetPostById(ctx, postIdint)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrorNotFound):
				app.notFound(w, r, err)
			default:
				app.internalServerError(w, r, err)

			}
			return
		}
		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromContext(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
