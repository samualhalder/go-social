package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrorNotFound = errors.New("record not round")
	ErrConflict   = errors.New("record already exists")
)

type Store struct {
	Post interface {
		Create(context.Context, *Post) error
		GetPostById(ctx context.Context, postId int64) (*Post, error)
		DeletePostById(ctx context.Context, postId int64) error
		UpdatePostById(ctx context.Context, post *Post) error
		GetUserFeedPosts(ctx context.Context, userId int64, p PaginatedFeedQuery) ([]PostWithMetaData, error)
	}
	User interface {
		Create(context.Context, *sql.Tx, *User) error
		GetById(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
	}
	Comment interface {
		GetCommentByPostId(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
	}
	Follower interface {
		Follow(context.Context, int64, int64) error
		UnFollow(context.Context, int64, int64) error
	}
	Role interface {
		GetByName(context.Context, string) (*Role, error)
	}
}

func NewStore(db *sql.DB) Store {
	return Store{
		Post:     &PostStore{db},
		User:     &UserStore{db},
		Comment:  &CommentStore{db},
		Follower: &FollowerStore{db},
		Role:     &RoleStore{db},
	}
}

func WithTx(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
