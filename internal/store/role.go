package store

import (
	"context"
	"database/sql"
	"fmt"
)

type Role struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

type RoleStore struct {
	db *sql.DB
}

func (r *RoleStore) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `SELECT name,id,description,level FROM roles WHERE name=$1`
	role := &Role{}
	fmt.Printf("hit here")
	err := r.db.QueryRowContext(ctx, query, name).Scan(&role.Name, &role.Id, &role.Description, &role.Level)
	if err != nil {
		fmt.Print("new err:", err)
		return nil, err
	}
	fmt.Printf("role in get", &role)
	return role, nil
}
