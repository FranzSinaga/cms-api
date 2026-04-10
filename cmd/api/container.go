package main

import (
	"github.com/FranzSinaga/blogcms/internal/feature/authentication"
	"github.com/jmoiron/sqlx"
)

type Container struct {
	AuthHandler *authentication.Handler
}

func NewContainer(db *sqlx.DB) *Container {
	// Authentication
	authRepo := authentication.NewUserRepository(db)
	authService := authentication.NewAuthService(authRepo)
	authHandler := authentication.NewAuthHandler(authService)

	return &Container{
		AuthHandler: authHandler,
	}
}
