package store

import (
	"context"
	"database/sql"
)

type CommentStore struct {
	db *sql.DB
}
type Comment struct {
	Id        int64  `json:"id"`
	PostId    int64  `json:"post_id"`
	UserId    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

func (c *CommentStore) GetCommentByPostId(ctx context.Context, postId int64) ([]Comment, error) {
	query := `SELECT a.id,a.post_id,a.content,a.created_at,b.id,b.username
			FROM comments a
			JOIN users b
			ON a.user_id=b.id
			WHERE a.post_id=$1
			ORDER BY a.created_at DESC;`
	comments := []Comment{}
	rows, err := c.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		var comment Comment
		err := rows.Scan(&comment.Id, &comment.PostId, &comment.Content, &comment.CreatedAt, &user.Id, &user.Username)
		if err != nil {
			return nil, err
		}
		comment.User = user
		comments = append(comments, comment)
	}
	return comments, nil
}

func (c *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `INSERT INTO comments (post_id,user_id,content) VALUES($1,$2,$3) RETURNING id,created_at`
	err := c.db.QueryRowContext(ctx, query, comment.PostId, comment.UserId, comment.Content).Scan(&comment.Id, &comment.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
