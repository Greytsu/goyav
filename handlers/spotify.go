package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"

	"goyav/goyavUser"
	spotifyGoyav "goyav/spotify"
)

func RegisterSpotify(public *gin.RouterGroup, auth *spotifyAuth.Authenticator, userService *goyavUser.Service, spotifyService *spotifyGoyav.Service) {
	public.GET("/spotify/auth", getAuthUrl(auth))
	public.GET("/spotify/callback", authCallback(auth, userService, spotifyService))
}

func getAuthUrl(auth *spotifyAuth.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := auth.AuthURL(os.Getenv("STATE"))
		c.IndentedJSON(http.StatusOK, url)
	}
}

func authCallback(auth *spotifyAuth.Authenticator, userService *goyavUser.Service, spotifyService *spotifyGoyav.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		// authenticate user
		tok, err := auth.Token(c.Request.Context(), os.Getenv("STATE"), c.Request)
		if err != nil {
			log.Info().Err(err).Msg("Couldn't get token")
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}
		if st := c.Request.FormValue("state"); st != os.Getenv("STATE") {
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
