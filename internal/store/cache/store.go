package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/samualhalder/go-social/internal/store"
)

type Store struct {
	User interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}
}

func NewRedisStore(rdb *redis.Client) Store {
	return Store{
		User: &UserStore{rdb: rdb},
	}
}
