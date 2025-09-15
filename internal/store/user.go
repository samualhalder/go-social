package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64        `json:"id"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Password  PasswordType `json:"-"`
	CreatedAt string       `json:"created_at"`
}

type PasswordType struct {
	password *string
	hash     []byte
}

func (p *PasswordType) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.password = &text
	p.hash = hash
	return nil
}

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `INSERT INTO users(username,email,password)
	VALUES($1,$2,$3) RETURNING id,created_at`
	err := tx.QueryRowContext(
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

func (u *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, envExpriry time.Duration) error {
	return WithTx(u.db, ctx, func(tx *sql.Tx) error {
		if err := u.Create(ctx, tx, user); err != nil {
			return err
		}
		if err := u.createUserInvitation(ctx, tx, token, envExpriry, user.Id); err != nil {
			return err
		}
		return nil
	})
}

func (u *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userId int64) error {
	query := `INSERT INTO user_invitations (user_id,token,expiry) values($1,$2,$3)`

	_, err := tx.ExecContext(ctx, query, userId, token, time.Now().Add(exp))
	if err != nil {
		return err
	}
	return nil
}
