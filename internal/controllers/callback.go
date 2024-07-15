package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sverdejot/espotifai/internal/auth"
	"github.com/sverdejot/espotifai/internal/views"
)

type sessionStorage interface {
	SetSession(key, value string)
}

func Callback(s sessionStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		codeParam, ok := r.URL.Query()["code"]
		if !ok || len(codeParam) == 0 {
			log.Fatal("cannot get params")
		}

		key, _ := uuid.NewRandom()
		s.SetSession(key.String(), codeParam[0])
		auth.SetCookie(w, key.String(), time.Now().Add(time.Hour))
		views.Callback().Render(r.Context(), w)
	}
}
