package controllers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/sverdejot/espotifai/internal/views"
)

type authEndpointRetriever interface {
	AuthEndpoint() string
}

func Index(auth authEndpointRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.Index(auth.AuthEndpoint())).ServeHTTP(w, r)
	}
}
