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

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		validationErrors := validator.GetValidationErrors(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Validation failed",
			"errors":  validationErrors,
		})
		return
	}

	if err := h.authService.Register(&req); err != nil {
		shared.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	shared.WriteSuccess(w, "User registered successfully", map[string]string{
		"email": req.Email,
		"name":  req.Name,
	})
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		validationErrors := validator.GetValidationErrors(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Validation failed",
			"errors":  validationErrors,
		})
		return
	}

	token, userResponse, err := h.authService.Login(&req)
	if err != nil {
		shared.WriteError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Set secure cookie based on environment
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   h.appConfig.Env == "production",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   int(h.jwtConfig.ExpiresIn.Seconds()),
	})

	w.Header().Set("Content-Type", "application/json")
	shared.WriteSuccess(w, "Login successful", domain.LoginResponse{
		Token: token,
		User:  userResponse,
	})
}
