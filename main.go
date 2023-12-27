package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"

	"goyav/config"
	"goyav/goyavUser"
	spotifyGoyav "goyav/handlers/spotify"
	mongoGoyav "goyav/mongo"
)

var (
	auth *spotifyAuth.Authenticator
	conf *config.Config
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Debug().Err(err).Msg("Error loading .env file")
	}

	conf = config.NewConfig()
	auth = spotifyAuth.New(spotifyAuth.WithRedirectURL(os.Getenv("REDIRECT_URL")), spotifyAuth.WithScopes(conf.Spotify.Scopes...))
}

func main() {
	// database connection
	mongoService, err := mongoGoyav.NewService(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to mongo. Exiting")
	}

	// goyavUser
	userRepo := goyavUser.NewRepository(mongoService)
	userServ := goyavUser.NewService(userRepo)

	router := gin.Default()

	//Routes
	spotifyGoyav.RegisterSpotify(router, auth, userServ)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server. Exiting")
	}
}
