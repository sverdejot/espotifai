package controllers

import (
	"log"
	"net/http"

	"github.com/sverdejot/espotifai/internal/infrastructure/http/clients/spotify"
	"github.com/sverdejot/espotifai/internal/model"
	"github.com/sverdejot/espotifai/internal/views"
)

type profileRetriever interface {
	RequestToken(code string) (string, error)
	Me(token string) model.Profile
}

func Profile(pr profileRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, _ := spotify.GetToken(r.Context())
		token, err := pr.RequestToken(code)
		if err != nil {
			log.Fatal(err)
		}
		me := pr.Me(token)

		views.Profile(me).Render(r.Context(), w)
	}
}
