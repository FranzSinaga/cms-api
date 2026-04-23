package authentication

import (
	appMiddleware "github.com/FranzSinaga/blogcms/internal/shared/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.AuthMiddleware)
			r.Post("/logout", h.Logout)
		})
	}
}
