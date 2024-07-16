package auth

import (
	"net/http"
	"time"
)

const (
	cookieKey = "Authorization"
)

func SetCookie(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:    cookieKey,
		Value:   token,
		Expires: expires,
	})
}
