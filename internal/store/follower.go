package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type FollowerStore struct {
	db *sql.DB
}
type Follower struct {
	UserId     int64  `json:"usre_id"`
	FollowerId int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

// TODO: pick the userId from token(authintication flow)
func (f *FollowerStore) Follow(ctx context.Context, follower int64, userId int64) error {
	query := `INSERT INTO followers(user_id,follower_id) VALUES($1,$2)`
	_, err := f.db.ExecContext(ctx, query, userId, follower)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
		return err
	}
	return nil

}

// TODO: pick the userId from token(authintication flow)
func (f *FollowerStore) UnFollow(ctx context.Context, follower int64, userId int64) error {
	query := `DELETE FROM  followers WHERE user_id=$1 AND follower_id=$2`
	_, err := f.db.ExecContext(ctx, query, userId, follower)
	if err != nil {
		return err
	}
	return err
}
