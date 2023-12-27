package handlers

import (
	"golang.org/x/oauth2"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"

	"goyav/goyavUser"
	spotifyGoyav "goyav/spotify"
)

func RegisterSpotify(public *gin.RouterGroup, private *gin.RouterGroup, auth *spotifyAuth.Authenticator, userService *goyavUser.Service, spotifyService *spotifyGoyav.Service) {
	public.GET("/spotify/auth", getAuthUrl(auth))
	public.GET("/spotify/callback", authCallback(auth, userService, spotifyService))
	private.GET("/spotify/me", getCurrentUser(spotifyService))
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

func getCurrentUser(spotifyService *spotifyGoyav.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get current user token from headers
		accessToken := c.GetHeader("access_token")
		tokenType := c.GetHeader("token_type")
		refreshToken := c.GetHeader("refresh_token")
		expiryStr := c.GetHeader("expiry")

		// parse time
		expiry, err := spotifyGoyav.ParseExpiry(expiryStr)
		if err != nil {
			log.Info().Err(err).Msg("Error while parsing time")
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		// create oauth2 token needed for the spotify client
		var tok oauth2.Token
		tok.AccessToken = accessToken
		tok.TokenType = tokenType
		tok.RefreshToken = refreshToken
		tok.Expiry = expiry

		goyavCurrentUser, err := spotifyService.GetCurrentUser(&tok)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.IndentedJSON(http.StatusOK, goyavCurrentUser)
	}
}
