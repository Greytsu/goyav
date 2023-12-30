package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"goyav/playlist"
)

func RegisterPlaylist(private *gin.RouterGroup, playlistService *playlist.Service) {
	private.GET("/playlists", getPlaylists(playlistService))
	private.GET("/playlists/:id", getPlaylist(playlistService))
	private.POST("/playlists", createPlaylist(playlistService))
}

func getPlaylists(playlistService *playlist.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		playlists, err := playlistService.GetPlaylists(c.Request.Context())
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.IndentedJSON(http.StatusOK, playlists)
	}
}

func getPlaylist(playlistService *playlist.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		newPlaylist, err := playlistService.GetPlaylist(c.Request.Context(), id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.IndentedJSON(http.StatusOK, newPlaylist)
	}
}

func createPlaylist(playlistService *playlist.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		owner := c.GetHeader("user_id")

		var newPlaylist *playlist.Playlist
		if err := c.BindJSON(&newPlaylist); err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Error while parsing JSON")
			return
		}

		// Setup owner
		newPlaylist.Owner = owner

		// Create token from headers
		tok, err := createTokenFromHeader(c)
		if err != nil {
			log.Info().Err(err).Msg("Error while creating oAuth2 token from headers")
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		// Create goyav playlist
		createdPlaylist, err := playlistService.CreatePlaylist(c.Request.Context(), newPlaylist, &tok)

		if err != nil {
			switch {
			case errors.Is(playlist.InvalidAttributesError, err) || errors.Is(playlist.UnknownAttributeError, err):
				c.IndentedJSON(http.StatusBadRequest, err.Error())
			default:
				c.IndentedJSON(http.StatusInternalServerError, err.Error())
			}
			return
		}

		c.IndentedJSON(http.StatusCreated, createdPlaylist)
	}
}
