package shared

import (
	"net/http"
	"time"
)

func SetAuthCookie(w http.ResponseWriter, token string, expiresIn time.Duration, isProduction bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthTokenCookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   int(expiresIn.Seconds()),
	})
}

func ClearAuthCookie(w http.ResponseWriter, isProduction bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthTokenCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Unix(0, 0), // ← set ke masa lalu → browser hapus cookie
		MaxAge:   -1,              // ← pastikan langsung dihapus
	})
}
