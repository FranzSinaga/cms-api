package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID           uuid.UUID  `db:"id"`
	Title        string     `db:"title"`
	Slug         string     `db:"slug"`
	Description  *string    `db:"description"`
	ThumbnailURL *string    `db:"thumbnail_url"`
	Content      string     `db:"content"`
	Status       string     `db:"status"`
	CreatedBy    uuid.UUID  `db:"created_by"`
	UpdatedBy    uuid.UUID  `db:"updated_by"`
	PublishedAt  *time.Time `db:"published_at"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

type CreatePostRequest struct {
	Title        string     `json:"title" validate:"required"`
	Slug         string     `json:"slug" validate:"required"`
	Description  *string    `json:"description" db:"description"`
	ThumbnailURL *string    `json:"thumbnail_url" db:"thumbnail_url"`
	Content      string     `json:"content" validate:"required"`
	Status       string     `json:"status" validate:"required"`
	PublishedAt  *time.Time `json:"-" db:"published_at"`
	CreatedBy    string     `json:"-" db:"created_by"`
	UpdatedBy    string     `json:"-" db:"updated_by"`
}
