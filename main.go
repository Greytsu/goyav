package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	spotifyGoyav "goyav/handlers/spotify"
	"os"
)

var (
	auth *spotifyauth.Authenticator
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Debug().Err(err).Msg("Error loading .env file")
	}
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(os.Getenv("REDIRECT_URL")), spotifyauth.WithScopes(spotifyauth.ScopePlaylistReadPrivate, spotifyauth.ScopeUserTopRead))
}

func main() {
	router := gin.Default()

	//Routes
	spotifyGoyav.RegisterSpotify(router, auth)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server. Exiting")
	}

}
