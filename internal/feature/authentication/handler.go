package authentication

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/FranzSinaga/blogcms/internal/domain"
	"github.com/FranzSinaga/blogcms/internal/service"
	"github.com/FranzSinaga/blogcms/internal/shared"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	if err := h.authService.Register(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Registrasi Berhasil",
	})
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteError(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	token, userResponse, err := h.authService.Login(&req)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		shared.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // true di production (https)
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	w.Header().Set("Content-Type", "application/json")
	shared.WriteSuccess(w, "Berhasil Login", domain.LoginResponse{
		Token: token,
		User:  userResponse,
	})
}
