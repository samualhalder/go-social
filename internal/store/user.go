package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64        `json:"id"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Password  PasswordType `json:"password"`
	CreatedAt string       `json:"created_at"`
	IsActive  bool         `json:"is_active"`
	Role      Role         `json:"role"`
	RoleId    int64        `json:"role_id"`
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
func (p *PasswordType) Check(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `INSERT INTO users(username,email,password,role_id)
	VALUES($1,$2,$3,$4) RETURNING id,created_at`
	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password.hash, user.Role.Id).Scan(&user.Id, &user.CreatedAt)
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
	query := `SELECT a.id,a.username,a.email,a.password,a.created_at,b.id,b.name,b.description,b.level FROM users a JOIN roles b ON a.role_id=b.id WHERE a.id=$1`
	user := &User{}
	err := u.db.
		QueryRowContext(ctx, query, userId).
		Scan(&user.Id, &user.Username, &user.Email, &user.Password.hash, &user.CreatedAt, &user.Role.Id, &user.Role.Name, &user.Role.Description, &user.Role.Level)
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

func (u *UserStore) Activate(ctx context.Context, token string) error {

	return WithTx(u.db, ctx, func(tx *sql.Tx) error {
		//1.find the user for the token
		user, err := u.findUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}
		//2.activate user
		user.IsActive = true
		if err := u.update(ctx, tx, user); err != nil {
			return err
		}
		//3.delete the invitiation from db
		if err := u.deleteUserInvitation(ctx, tx, user.Id); err != nil {
			return err
		}
		return nil
	})
}

func (u *UserStore) findUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `SELECT a.id,a.username,a.email,a.is_active,a.created_at FROM
			users a JOIN user_invitations b
			ON a.id=b.user_id
			WHERE b.token=$1 AND b.expiry>$2`

	user := &User{}

	hash := sha256.Sum256([]byte(token)) // not readble by human so cant store in sql
	hashToken := hex.EncodeToString(hash[:])
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(&user.Id, &user.Username, &user.Email, &user.IsActive, &user.CreatedAt)
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

func (u *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id,email,username,created_at,password FROM users where email=$1 AND is_active=true`
	user := &User{}
	err := u.db.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Email, &user.Username, &user.CreatedAt, &user.Password.hash)
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

func (u *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET username=$1,email=$2,is_active=$3 WHERE id=$4`
	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.Id)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userId int64) error {
	query := `DELETE FROM user_invitations where user_id=$1`
	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) delete(ctx context.Context, tx *sql.Tx, userId int64) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserStore) Delete(ctx context.Context, userId int64) error {
	return WithTx(u.db, ctx, func(tx *sql.Tx) error {
		err := u.delete(ctx, tx, userId)
		if err != nil {
			return err
		}
		err = u.deleteUserInvitation(ctx, tx, userId)
		if err != nil {
			return err
		}
		return nil
	})
}
