package middlewares

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"net/http"

	"github.com/gin-gonic/gin"

	"goyav/spotify"
)

// myMiddleware is the Gin middleware that checks headers and adds a new header
func CheckCredentials(spotifyService *spotify.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		accessToken := c.GetHeader("access_token")
		tokenType := c.GetHeader("token_type")
		refreshToken := c.GetHeader("refresh_token")
		expiryStr := c.GetHeader("expiry")

		if accessToken != "" && tokenType != "" && refreshToken != "" && expiryStr != "" {
			// parse time
			expiry, err := spotify.ParseExpiry(expiryStr)
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

			if !spotifyService.CheckUser(&tok) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				c.Abort() // Abort further processing
				return
			}
		}

		c.Next()
	}
}
