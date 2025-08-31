package store

import (
	"context"
	"database/sql"
	"errors"
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
		Create(context.Context, *User) error
		GetById(context.Context, int64) (*User, error)
	}
	Comment interface {
		GetCommentByPostId(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
	}
	Follower interface {
		Follow(context.Context, int64, int64) error
		UnFollow(context.Context, int64, int64) error
	}
}

func NewStore(db *sql.DB) Store {
	return Store{
		Post:     &PostStore{db},
		User:     &UserStore{db},
		Comment:  &CommentStore{db},
		Follower: &FollowerStore{db},
	}
}
