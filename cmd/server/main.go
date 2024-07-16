package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sverdejot/espotifai/cmd/server/bootstrap"
	"github.com/sverdejot/espotifai/internal/controllers"
	"github.com/sverdejot/espotifai/internal/controllers/middleware"
	"github.com/sverdejot/espotifai/internal/infrastructure/http/clients/spotify"
)

func main() {
	config := bootstrap.ReadConfig()
	mux := http.NewServeMux()

	client := spotify.NewClient(config.ClientId, config.ClientSecret)

	authMiddleware := middleware.NewAuthMiddleware(client)

	mux.HandleFunc("GET /", controllers.Index(config.SpotifyAuthUrl))
	mux.HandleFunc("GET /callback", controllers.Callback(authMiddleware))
	mux.Handle("GET /me", authMiddleware.Use(controllers.Profile(client)))
	mux.Handle("GET /me/top/artists", authMiddleware.Use(controllers.TopArtists(client)))
	mux.Handle("GET /me/top/tracks", authMiddleware.Use(controllers.TopTracks(client)))

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Port),

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
