package bootstrap

import (
	"log"

	"github.com/caarlos0/env/v11"
)


type Config struct {
	ClientId     string `env:"AUTH_CLIENT_ID"`
	ClientSecret string `env:"AUTH_CLIENT_SECRET"`
}

func ReadConfig() (config Config) {
	opts := env.Options{
		Prefix: "ESPOTIFAI_",
	}
	if err := env.ParseWithOptions(&config, opts); err != nil {
		log.Fatal(err)
	}

	return config
}
