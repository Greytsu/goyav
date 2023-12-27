package spotify

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"goyav/goyavUser"
	"net/http"
	"os"
	"time"
)

func RegisterSpotify(router *gin.Engine, auth *spotifyauth.Authenticator, userService *goyavUser.Service) {
	router.GET("/spotify/auth", getAuthUrl(auth))
	router.GET("/spotify/callback", authCallback(auth, userService))
	router.GET("/spotify/me", getCurrentUser(auth))
}

func getAuthUrl(auth *spotifyauth.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := auth.AuthURL(os.Getenv("STATE"))
		c.IndentedJSON(http.StatusOK, url)
	}
}

func authCallback(auth *spotifyauth.Authenticator, userService *goyavUser.Service) gin.HandlerFunc {
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

		// create client for the current user
		client := spotify.New(auth.Client(c.Request.Context(), tok))

		// get current user info
		user, err := client.CurrentUser(context.Background())
		if err != nil {
			log.Info().Err(err).Msg("Error while fetching for current currentUser")
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		// upsert user
		goyavCurrentUser := goyavUser.MapSpotifyUserToGoyavUser(user)
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

func getCurrentUser(auth *spotifyauth.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get current user token from headers
		accessToken := c.GetHeader("access_token")
		tokenType := c.GetHeader("token_type")
		refreshToken := c.GetHeader("refresh_token")
		expiryStr := c.GetHeader("expiry")

		// parse time
		expiry, err := parseExpiry(expiryStr)
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

		// new spotify client
		client := spotify.New(auth.Client(c.Request.Context(), &tok))

		// retrieve current user info
		user, err := client.CurrentUser(context.Background())
		if err != nil {
			log.Info().Err(err).Msg("Error while fetching for current goyavUser")
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}
		goyavCurrentUser := goyavUser.MapSpotifyUserToGoyavUser(user)

		c.IndentedJSON(http.StatusOK, goyavCurrentUser)
	}
}

// Parses a string to Time using this layout : "2006-01-02T15:04:05.999999Z07:00"
func parseExpiry(expiryStr string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.999999Z07:00"

	expiry, err := time.Parse(layout, expiryStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Now(), err
	}

	return expiry, nil
}
