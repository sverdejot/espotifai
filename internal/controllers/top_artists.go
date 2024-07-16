package controllers

import (
	"log"
	"net/http"

	"github.com/sverdejot/espotifai/internal/infrastructure/http/clients/spotify"
	"github.com/sverdejot/espotifai/internal/model"
	"github.com/sverdejot/espotifai/internal/views"
)

type topArtistsFetcher interface {
	Artists(code string) model.TopArtists
	Tracks(code string) model.TopArtists
	RequestToken(code string) (string, error)
}

func TopArtists(af topArtistsFetcher) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		code, _ := spotify.GetToken(r.Context())
		token, err := af.RequestToken(code)
		if err != nil {
			log.Fatal(err)
		}
		topArtists := af.Artists(token)

		views.TopArtists(topArtists).Render(r.Context(), w)
	})
}

func TopTracks(af topArtistsFetcher) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		code, _ := spotify.GetToken(r.Context())
		token, err := af.RequestToken(code)
		if err != nil {
			log.Fatal(err)
		}
		topArtists := af.Tracks(token)

		views.TopArtists(topArtists).Render(r.Context(), w)
	})
}
