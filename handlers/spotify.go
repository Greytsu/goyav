package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"

	"goyav/config"
	"goyav/goyavUser"
	spotifyGoyav "goyav/spotify"
)

func RegisterSpotify(public *gin.RouterGroup, private *gin.RouterGroup, auth *spotifyAuth.Authenticator, conf *config.Config, userService *goyavUser.Service, spotifyService *spotifyGoyav.Service) {
	public.GET("/spotify/auth", getAuthUrl(auth, conf))
	public.GET("/spotify/callback", authCallback(auth, conf, userService, spotifyService))
	private.GET("/spotify/me/top/tracks", getTopTracks(spotifyService))
	private.GET("/spotify/me/top/artists", getTopArtist(spotifyService))
	private.GET("/spotify/me/recommendations", getRecommendations(spotifyService))
}

func getAuthUrl(auth *spotifyAuth.Authenticator, conf *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := auth.AuthURL(conf.Spotify.State)
		c.IndentedJSON(http.StatusOK, url)
	}
}

func authCallback(auth *spotifyAuth.Authenticator, conf *config.Config, userService *goyavUser.Service, spotifyService *spotifyGoyav.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		// authenticate user
		tok, err := auth.Token(c.Request.Context(), conf.Spotify.State, c.Request)
		if err != nil {
			log.Info().Err(err).Msg("Couldn't get token")
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}
		if st := c.Request.FormValue("state"); st != conf.Spotify.State {
			log.Info().Err(err).Msg("State mismatch")
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		goyavCurrentUser, err := spotifyService.GetCurrentUser(tok)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		// upsert user
		err = userService.CreateUser(
			c.Request.Context(),
			goyavCurrentUser)
		if err != nil {
			log.Info().Err(err).Msg("Error while creating user")
		}

		// return the token
		c.IndentedJSON(http.StatusOK, tok)
	}
}

func getTopTracks(spotifyService *spotifyGoyav.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		tok, err := createTokenFromHeader(c)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		tracks, err := spotifyService.GetCurrentUserTopTracks(&tok)
		c.IndentedJSON(http.StatusOK, tracks)
	}
}

func getTopArtist(spotifyService *spotifyGoyav.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		tok, err := createTokenFromHeader(c)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		artists, err := spotifyService.GetCurrentUserTopArtists(&tok)
		c.IndentedJSON(http.StatusOK, artists)
	}
}

func getRecommendations(spotifyService *spotifyGoyav.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		tok, err := createTokenFromHeader(c)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		recommendations, err := spotifyService.GetFullRecommendations(&tok)
		c.IndentedJSON(http.StatusOK, recommendations)
	}
}
