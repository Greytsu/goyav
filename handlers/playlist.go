package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"goyav/playlist"
	"goyav/spotify"
)

func RegisterPlaylist(private *gin.RouterGroup, playlistService *playlist.Service) {
	private.GET("/playlists", getPlaylists(playlistService))
	private.GET("/playlists/:id", getPlaylist(playlistService))
	private.POST("/playlists", createPlaylist(playlistService))
	private.POST("/playlists/:id/add-tracks", addRecommendationsToPlaylist(playlistService))
	private.PUT("/playlists/:id/contributors", updateContributors(playlistService))
	private.POST("/playlists/:id/update", updateSpotifyPlaylist(playlistService))
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
		playlistID := c.Param("id")
		userID := c.GetHeader("user_id")

		foundPlaylist, err := playlistService.GetPlaylist(c.Request.Context(), playlistID, userID)
		if err != nil {
			switch {
			case errors.Is(playlist.UserNotContributorError, err):
				c.IndentedJSON(http.StatusUnauthorized, err.Error())
			default:
				c.IndentedJSON(http.StatusInternalServerError, err.Error())
			}
			return
		}

		c.IndentedJSON(http.StatusOK, foundPlaylist)
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

func addRecommendationsToPlaylist(playlistService *playlist.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		playlistID := c.Param("id")
		userID := c.GetHeader("user_id")

		var recommendationRequest *spotify.RecommendationRequest
		if err := c.BindJSON(&recommendationRequest); err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Error while parsing JSON")
			return
		}

		// Create token from headers
		tok, err := createTokenFromHeader(c)
		if err != nil {
			log.Info().Err(err).Msg("Error while creating oAuth2 token from headers")
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		err = playlistService.AddRecommendationsToPlaylist(c.Request.Context(), &tok, playlistID, userID, recommendationRequest)
		if err != nil {
			switch {
			case errors.Is(playlist.UserNotContributorError, err):
				c.IndentedJSON(http.StatusUnauthorized, err.Error())
			case errors.Is(spotify.InvalidRecommendationRequestError, err):
				c.IndentedJSON(http.StatusBadRequest, err.Error())
			default:
				c.IndentedJSON(http.StatusInternalServerError, err.Error())
			}
			return
		}

		c.IndentedJSON(http.StatusOK, nil)
	}
}

func updateContributors(playlistService *playlist.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		playlistID := c.Param("id")
		userID := c.GetHeader("user_id")

		var contributors []string
		if err := c.BindJSON(&contributors); err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Error while parsing JSON")
			return
		}

		err := playlistService.UpdateContributors(c.Request.Context(), playlistID, userID, contributors)
		if err != nil {
			switch {
			case errors.Is(playlist.UserNotOwnerError, err):
				c.IndentedJSON(http.StatusUnauthorized, err.Error())
			default:
				c.IndentedJSON(http.StatusBadRequest, err.Error())
			}
			return
		}
		c.IndentedJSON(http.StatusOK, nil)
	}
}

func updateSpotifyPlaylist(playlistService *playlist.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		playlistID := c.Param("id")
		userID := c.GetHeader("user_id")

		// Create token from headers
		tok, err := createTokenFromHeader(c)
		if err != nil {
			log.Info().Err(err).Msg("Error while creating oAuth2 token from headers")
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		err = playlistService.UpdateSpotifyPlaylistTracks(c.Request.Context(), &tok, playlistID, userID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.IndentedJSON(http.StatusOK, nil)
	}
}
