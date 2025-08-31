package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	Id        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserId    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	Version   int       `json:"version"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetaData struct {
	Post
	CommentCount int `json:"comment_count"`
}

type PostStore struct {
	db *sql.DB
}

func (p *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content,title,user_id,tags)
	 VALUES($1,$2,$3,$4) RETURNING id , created_at, updated_at`
	err := p.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserId,
		pq.Array(post.Tags),
	).Scan(&post.Id, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostStore) GetPostById(ctx context.Context, postId int64) (*Post, error) {
	query := `SELECT id,title,content,tags,user_id,created_at,updated_at,version FROM posts where id=$1`
	var post Post
	err := p.db.
		QueryRowContext(ctx, query, postId).
		Scan(&post.Id, &post.Title, &post.Content, pq.Array(&post.Tags), &post.UserId, &post.CreatedAt, &post.UpdatedAt, &post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (p *PostStore) DeletePostById(ctx context.Context, postId int64) error {
	query := `DELETE FROM posts WHERE id=$1`
	res, err := p.db.ExecContext(ctx, query, postId)

	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrorNotFound
	}
	return nil
}

func (p *PostStore) UpdatePostById(ctx context.Context, post *Post) error {
	query := `UPDATE posts SET title=$1,content=$2,version=version+1 WHERE id=$3 AND version=$4 returning version`
	err := p.db.QueryRowContext(ctx, query, post.Title, post.Content, post.Id, post.Version).Scan(&post.Version)
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

func (p *PostStore) GetUserFeedPosts(ctx context.Context, userId int64, filters PaginatedFeedQuery) ([]PostWithMetaData, error) {
	query := `SELECT
				p.id,p.title,p.user_id,p.content,p.tags,p.created_at,COUNT(distinct c.id) AS comment_count
				FROM Posts p
				JOIN COMMENTS c ON c.post_id=p.id
				JOIN followers f ON p.user_id=f.follower_id OR p.user_id=$1
				WHERE f.user_id=$1 OR p.user_id=$1
				GROUP BY p.id
				ORDER BY p.created_at ` + filters.Sort +
		` LIMIT $2
				OFFSET $3
				`
	var posts []PostWithMetaData
	rows, err := p.db.QueryContext(ctx, query, userId, filters.Limit, filters.Offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post PostWithMetaData
		err := rows.Scan(&post.Id, &post.Title, &post.UserId, &post.Content, pq.Array(&post.Tags), &post.CreatedAt, &post.CommentCount)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
