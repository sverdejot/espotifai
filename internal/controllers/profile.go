package controllers

import (
	"log"
	"net/http"

	"github.com/sverdejot/espotifai/internal/model"
	"github.com/sverdejot/espotifai/internal/views"
)

type profileRetriever interface {
	RequestToken(code string)	(string, error) 
	Me(token string) model.Profile
}

func Profile(pr profileRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		codeParam, ok := r.URL.Query()["code"]
		if !ok || len(codeParam) == 0 {
			log.Fatal("cannot get params")
		}
		token, err := pr.RequestToken(codeParam[0])
		if err != nil {
			log.Fatal(err)
		}
		me := pr.Me(token)

		views.Profile(me).Render(r.Context(), w)
	}
}
