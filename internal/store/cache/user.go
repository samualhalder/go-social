package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/samualhalder/go-social/internal/store"
)

type UserStore struct {
	rdb *redis.Client
}

func (u *UserStore) Get(ctx context.Context, userId int64) (*store.User, error) {

	cacheKey := fmt.Sprintf("user-%v", userId)
	data, err := u.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var user *store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), user)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}
func (u *UserStore) Set(ctx context.Context, user *store.User) error {

	cacheKey := fmt.Sprintf("user-%v", user.Id)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return u.rdb.SetEX(ctx, cacheKey, string(json), time.Minute).Err()
}
