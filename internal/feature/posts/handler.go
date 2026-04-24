package posts

import (
	"encoding/json"
	"net/http"

	"github.com/FranzSinaga/blogcms/internal/domain"
	"github.com/FranzSinaga/blogcms/internal/shared"
	"github.com/FranzSinaga/blogcms/internal/shared/validator"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	postService *Service
}

func NewPostHandler(postService *Service) *Handler {
	return &Handler{
		postService: postService,
	}
}

func (h *Handler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	posts, err := h.postService.GetAllPosts()
	if err != nil {
		shared.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	shared.WriteSuccess(w, "Get all post successfully", &posts)
}

func (h *Handler) GetPostBySlug(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	slug := chi.URLParam(r, "slug")
	posts, err := h.postService.GetPostBySlug(slug)
	if err != nil {
		shared.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	shared.WriteSuccess(w, "Get all post successfully", &posts)
}

func (h *Handler) CreateNewPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req domain.CreatePostRequest
	user := shared.GetUserFromContext(r)
	if user == nil {
		shared.WriteError(w, "Unauthorized: user not found in context", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&req); err != nil {
		validationErrors := validator.GetValidationErrors(err)
		w.Header().Set("Content-Type", "application/json")
		shared.WriteValidationError(w, http.StatusBadRequest, validationErrors)
		return
	}

	post, err := h.postService.CreateNewPost(&req, user.UserID)
	if err != nil {
		shared.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	shared.WriteSuccess(w, "Post created successfully", &post)
}
