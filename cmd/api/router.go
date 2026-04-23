package main

import (
	"fmt"
	"net/http"

	"github.com/FranzSinaga/blogcms/internal/shared"
	appMiddleware "github.com/FranzSinaga/blogcms/internal/shared/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func setupRouter(c *Container) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		shared.WriteSuccess(w, "Api is healthy", "API is healthy")
	})

	// Public Routes
	r.Route("/api/auth", c.AuthHandler.Routes())

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware)

		r.Get("/api/protected-test", func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value(shared.UserContextKey).(*appMiddleware.UserClaim)
			shared.WriteSuccess(w, "Successfully accessed the protected route", fmt.Sprintf("Hello %s, you have successfully accessed the protected route!", user.Name))
		})

		r.Get("/api/check-login", func(w http.ResponseWriter, r *http.Request) {
			shared.WriteSuccess(w, "User is logged in", true)
		})
	})

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            false,
	})

	return corsHandler.Handler(r)
}
