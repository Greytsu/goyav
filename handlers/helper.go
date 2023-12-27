package handlers

import (
	"golang.org/x/oauth2"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	spotifyGoyav "goyav/spotify"
)

func createTokenFromHeader(c *gin.Context) (oauth2.Token, error) {
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
		return oauth2.Token{}, err
	}

	// create oauth2 token needed for the spotify client
	var tok oauth2.Token
	tok.AccessToken = accessToken
	tok.TokenType = tokenType
	tok.RefreshToken = refreshToken
	tok.Expiry = expiry

	return tok, nil
}
