package controllers

import (
	"net/http"
	"net/url"

	"github.com/a-h/templ"
	"github.com/sverdejot/espotifai/internal/views"
)

func Index(spotifyAuthEndpoint url.URL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.Index(spotifyAuthEndpoint.String())).ServeHTTP(w, r)
	}
}
