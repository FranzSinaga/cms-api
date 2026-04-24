package posts

import (
	"time"

	"github.com/FranzSinaga/blogcms/internal/domain"
)

type Service struct {
	postRepo RepositoryInterface
}

func NewPostService(postRepo RepositoryInterface) *Service {
	return &Service{
		postRepo: postRepo,
	}
}

func (s *Service) GetAllPosts() ([]*domain.Post, error) {
	posts, err := s.postRepo.GetAllPosts()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *Service) GetPostBySlug(slug string) (*domain.Post, error) {
	post, err := s.postRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Service) CreateNewPost(req *domain.CreatePostRequest, userID string) (*domain.Post, error) {
	var published_at *time.Time
	if req.Status == "published" {
		now := time.Now()
		published_at = &now
	}

	post := &domain.CreatePostRequest{
		Title:        req.Title,
		Slug:         req.Slug,
		Description:  req.Description,
		ThumbnailURL: req.ThumbnailURL,
		Content:      req.Content,
		Status:       req.Status,
		PublishedAt:  published_at,
		CreatedBy:    userID,
		UpdatedBy:    userID,
	}

	if err := s.postRepo.CreateNewPost(post); err != nil {
		return nil, err
	}

	createdPost, err := s.postRepo.FindBySlug(post.Slug)
	if err != nil {
		return nil, err
	}
	return createdPost, nil
}
