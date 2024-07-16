package middleware

import (
	"net/http"
	"time"

	"github.com/sverdejot/espotifai/internal/infrastructure/http/clients/spotify"
)


type AuthMiddleware struct {
	client *spotify.Client

	sessions map[string]spotify.Session
}

func NewAuthMiddleware(client *spotify.Client) *AuthMiddleware {
	return &AuthMiddleware{
		client:   client,
		sessions: make(map[string]spotify.Session),
	}
}

func (m *AuthMiddleware) SetSession(key, token string) {
	m.sessions[key] = spotify.Session{
		Token: token, 
		ExpireAt: time.Now().Add(time.Hour),
	}
}

func (m *AuthMiddleware) Use(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieKey)

		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		session, ok := m.sessions[cookie.Value]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if session.IsExpired() {
			delete(m.sessions, cookie.Value)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := spotify.SetToken(r.Context(), session.Token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
