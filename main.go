package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"

	"goyav/config"
	"goyav/goyavUser"
	handlers "goyav/handlers"
	"goyav/middlewares"
	mongoGoyav "goyav/mongo"
	"goyav/playlist"
	spotifyGoyav "goyav/spotify"
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
	auth = spotifyAuth.New(spotifyAuth.WithRedirectURL(conf.Spotify.RedirectURL), spotifyAuth.WithScopes(conf.Spotify.Scopes...))
}

func main() {
	// database connection
	mongoService, err := mongoGoyav.NewService(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to mongo. Exiting")
	}

	// spotify
	spotifyServ := spotifyGoyav.NewService(auth)

	// goyavUser
	userRepo := goyavUser.NewRepository(mongoService)
	userServ := goyavUser.NewService(userRepo)

	// goyavPlaylist
	playlistRepo := playlist.NewRepository(mongoService)
	playlistServ := playlist.NewService(playlistRepo)

	router := gin.Default()

	// public routes
	public := router.Group("/api/v1")

	// private routes
	private := router.Group("/api/v1")
	private.Use(middlewares.CheckCredentials(spotifyServ))

	//Routes
	handlers.RegisterSpotify(public, private, auth, conf, userServ, spotifyServ)
	handlers.RegisterPlaylist(private, playlistServ)
	handlers.RegisterUser(private, userServ)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server. Exiting")
	}
}
