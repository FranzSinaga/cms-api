package posts

import (
	"github.com/FranzSinaga/blogcms/internal/domain"
	"github.com/jmoiron/sqlx"
)

type RepositoryInterface interface {
	CreateNewPost(post *domain.CreatePostRequest) error
	FindBySlug(slug string) (*domain.Post, error)
	GetAllPosts() ([]*domain.Post, error)
}

type Repository struct {
	db *sqlx.DB
}

func NewPostsRepository(db *sqlx.DB) RepositoryInterface {
	return &Repository{db: db}
}

func (r *Repository) GetAllPosts() ([]*domain.Post, error) {
	var posts []*domain.Post
	query := `SELECT * FROM posts`
	err := r.db.Select(&posts, query)

	return posts, err
}

func (r *Repository) CreateNewPost(post *domain.CreatePostRequest) error {
	query := `
		INSERT INTO posts (title, slug, description, thumbnail_url, content, status, published_at, created_by, updated_by)
		VALUES (:title, :slug, :description, :thumbnail_url, :content, :status, :published_at, :created_by, :updated_by)
	`
	_, err := r.db.NamedExec(query, post)
	return err
}

func (r *Repository) FindBySlug(slug string) (*domain.Post, error) {
	post := &domain.Post{}
	query := `SELECT * FROM posts WHERE slug = $1`
	err := r.db.Get(post, query, slug)
	if err != nil {
		return nil, err
	}
	return post, nil
}
