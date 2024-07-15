package bootstrap

import (
	"log"
	"net/url"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Host         string  `env:"HOST" envDefault:"localhost"`
	Port         int     `env:"PORT" envDefault:"8080"`
	ClientId     string  `env:"AUTH_CLIENT_ID,notEmpty"`
	ClientSecret string  `env:"AUTH_CLIENT_SECRET,notEmpty"`
	CallbackUrl  url.URL `env:"AUTH_CALLBACK_URL,expand" envDefault:"http://${HOST}:${PORT}/callback"`

	SpotifyAuthUrl url.URL
}

func ReadConfig() (config Config) {
	opts := env.Options{
		Prefix: "ESPOTIFAI_",
	}
	if err := env.ParseWithOptions(&config, opts); err != nil {
		log.Fatalf("%+v\n", err)
	}

	config.SpotifyAuthUrl = url.URL{
		Scheme: "https",
		Host: "accounts.spotify.com",
		Path: "authorize",

		RawQuery: url.Values{
			"response_type": { "code" },
			"client_id": { config.ClientId },
			"redirect_uri": { config.CallbackUrl.String() },
		}.Encode(),
	}

	return config
}
