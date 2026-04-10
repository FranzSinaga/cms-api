package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/FranzSinaga/blogcms/internal/handler"
	"github.com/FranzSinaga/blogcms/internal/repository"
	"github.com/FranzSinaga/blogcms/internal/service"
	"github.com/FranzSinaga/blogcms/internal/shared/middleware"
	"github.com/FranzSinaga/blogcms/pkg/config"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	fmt.Println("Database connected!")

	// Init Layer
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "BlogCMS API Is Running")
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Ok")
	})

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/protected-test", func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value(middleware.UserContextKey).(*middleware.UserClaim)
			fmt.Fprintf(w, "Halo %s, kamu berhasil masuk protected route!", user.Email)
		})
	})

	fmt.Printf("Server starting in %s...", port)
	http.ListenAndServe(":"+port, r)
}
