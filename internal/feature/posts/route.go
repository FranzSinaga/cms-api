package posts

import "github.com/go-chi/chi/v5"

func (h *Handler) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/get", h.GetAllPosts)
		r.Get("/get/{slug}", h.GetPostBySlug)
		r.Post("/create", h.CreateNewPost)
	}
}
