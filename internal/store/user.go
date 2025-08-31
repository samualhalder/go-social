package store

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users(username,email,password)
	VALUES($1,$2,$3) RETURNING id,created_at`
	err := u.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password).Scan(&user.Id, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrorNotFound
		default:
			return err
		}
	}
	return nil
}

func (u *UserStore) GetById(ctx context.Context, userId int64) (*User, error) {
	query := `SELECT id,username,email,password,created_at FROM users WHERE id=$1`
	user := &User{}
	err := u.db.
		QueryRowContext(ctx, query, userId).
		Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}
