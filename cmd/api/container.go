package main

import (
	"github.com/FranzSinaga/blogcms/internal/feature/authentication"
	"github.com/FranzSinaga/blogcms/pkg/config"
	"github.com/jmoiron/sqlx"
)

type Container struct {
	AuthHandler *authentication.Handler
}

func NewContainer(db *sqlx.DB, cfg *config.Config) *Container {
	// Authentication
	authRepo := authentication.NewUserRepository(db)
	authService := authentication.NewAuthService(authRepo, cfg.JWT)
	authHandler := authentication.NewAuthHandler(authService, cfg.App, cfg.JWT)

	return &Container{
		AuthHandler: authHandler,
	}
}
