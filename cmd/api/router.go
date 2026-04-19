package main

import (
	"fmt"
	"net/http"

	"github.com/FranzSinaga/blogcms/internal/shared"
	appMiddleware "github.com/FranzSinaga/blogcms/internal/shared/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func setupRouter(c *Container) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "api-cms is running")
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// Public Routes
	r.Route("/auth", c.AuthHandler.Routes())

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware)

		r.Get("/protected-test", func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value(appMiddleware.UserContextKey).(*appMiddleware.UserClaim)
			fmt.Fprintf(w, "Hello %s, you have successfully accessed the protected route!", user.Email)
		})

		r.Get("/check-login", func(w http.ResponseWriter, r *http.Request) {
			shared.WriteSuccess(w, "User is logged in", true)
		})
	})
	return r
}
