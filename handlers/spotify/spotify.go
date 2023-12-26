package spotify

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

func RegisterSpotify(router *gin.Engine, auth *spotifyauth.Authenticator) {
	router.GET("/spotify/auth", getAuthUrl(auth))
	router.GET("/spotify/callback", authCallback(auth))
	router.GET("/spotify/me", getCurrentUser(auth))
}

func getAuthUrl(auth *spotifyauth.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := auth.AuthURL(os.Getenv("STATE"))
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
		c.IndentedJSON(http.StatusOK, url)
	}
}

func authCallback(auth *spotifyauth.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// return the token
		c.IndentedJSON(http.StatusOK, tok)
	}
}

func getCurrentUser(auth *spotifyauth.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {

		var tok *oauth2.Token
		if err := c.BindJSON(&tok); err != nil {
			log.Info().Err(err).Msg("Error while parsing JSON")
			c.IndentedJSON(http.StatusBadRequest, "Error while parsing JSON")
			return
		}

		client := spotify.New(auth.Client(c.Request.Context(), tok))

		user, err := client.CurrentUser(context.Background())
		if err != nil {
			log.Info().Err(err).Msg("Error while fetching for current user")
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}
		log.Info().Str("user id", user.ID).Msg("User logged in")

		c.IndentedJSON(http.StatusOK, user.ID)
	}
}
