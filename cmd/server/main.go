package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sverdejot/espotifai/cmd/server/bootstrap"
	"github.com/sverdejot/espotifai/internal/auth"
	"github.com/sverdejot/espotifai/internal/controllers"
)


func main() {
	config := bootstrap.ReadConfig()
	mux := http.NewServeMux()

	auth_client := auth.NewSpotifyAuth(config.ClientId, config.ClientSecret)

	mux.HandleFunc("GET /index", controllers.Index(auth_client))
	mux.HandleFunc("GET /callback", controllers.Profile(auth_client))

	server := &http.Server{
		Addr: ":8080",

		ReadTimeout:       2 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       2 * time.Minute,
		ReadHeaderTimeout: time.Second,

		Handler: mux,
	}

	registerShutdown(server)
	log.Println("starting server at ", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func registerShutdown(server *http.Server) {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		<-ctx.Done()
		cancel()
		log.Println("gracefully shutting down server at", server.Addr)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("error while shutting down server: ", err)
		}
		log.Println("server shutdown successfully")
	}()
}
