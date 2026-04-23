package authentication

import (
	"encoding/json"
	"net/http"

	"github.com/FranzSinaga/blogcms/internal/domain"
	"github.com/FranzSinaga/blogcms/internal/shared"
	"github.com/FranzSinaga/blogcms/internal/shared/validator"
	"github.com/FranzSinaga/blogcms/pkg/config"
)

type Handler struct {
	authService *Service
	appConfig   config.AppConfig
	jwtConfig   config.JWTConfig
}

func NewAuthHandler(authService *Service, appConfig config.AppConfig, jwtConfig config.JWTConfig) *Handler {
	return &Handler{
		authService: authService,
		appConfig:   appConfig,
		jwtConfig:   jwtConfig,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		validationErrors := validator.GetValidationErrors(err)
		w.Header().Set("Content-Type", "application/json")
		shared.WriteValidationError(w, http.StatusBadRequest, validationErrors)
		return
	}

	token, _, err := h.authService.Register(&req)
	if err != nil {
		shared.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shared.SetAuthCookie(w, token, h.jwtConfig.ExpiresIn, h.appConfig.Env == "production")

	w.Header().Set("Content-Type", "application/json")
	shared.WriteSuccess(w, "User registered successfully", map[string]string{
		"email": req.Email,
		"name":  req.Name,
		"token": token,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		validationErrors := validator.GetValidationErrors(err)
		w.Header().Set("Content-Type", "application/json")
		shared.WriteValidationError(w, http.StatusBadRequest, validationErrors)
		return
	}

	token, userResponse, err := h.authService.Login(&req)
	if err != nil {
		shared.WriteError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Set secure cookie based on environment
	shared.SetAuthCookie(w, token, h.jwtConfig.ExpiresIn, h.appConfig.Env == "production")

	w.Header().Set("Content-Type", "application/json")
	shared.WriteSuccess(w, "Login successful", domain.LoginResponse{
		Token: token,
		User:  userResponse,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	shared.ClearAuthCookie(w, h.appConfig.Env == "production")
}
