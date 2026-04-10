package authentication

import "github.com/go-chi/chi/v5"

func (h *Handler) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
	}
}
