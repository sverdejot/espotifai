package auth

import (
	"net/http"
	"time"

	"github.com/sverdejot/espotifai/internal/infrastructure/http/clients"
)

const (
	cookieKey = "Authorization"
)

type AuthMiddleware struct {
	client *clients.SpotifyClient

	sessions map[string]session
}

func NewMiddleware(client *clients.SpotifyClient) *AuthMiddleware {
	return &AuthMiddleware{
		client:   client,
		sessions: make(map[string]session),
	}
}

func (m *AuthMiddleware) SetSession(key, token string) {
	m.sessions[key] = session{
		token, time.Now().Add(time.Hour),
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

		ctx := SetSpotifyToken(r.Context(), session.Token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
