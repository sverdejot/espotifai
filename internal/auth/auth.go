package auth

import (
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:    cookieKey,
		Value:   token,
		Expires: expires,
	})
}
