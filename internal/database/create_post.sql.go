// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: create_post.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, created_at, updated_at, title, url, description, published_at, feed_id
`

type CreatePostParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description string
	PublishedAt time.Time
	FeedID      uuid.NullUUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}
