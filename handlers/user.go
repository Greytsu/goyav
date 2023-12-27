package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"goyav/goyavUser"
)

func RegisterUser(private *gin.RouterGroup, userService *goyavUser.Service) {
	private.GET("/users/:id", getCurrentUser(userService))
}

func getCurrentUser(userService *goyavUser.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		goyavCurrentUser, err := userService.GetUser(c.Request.Context(), id)
		log.Info().Interface("goyavCurrentUser", goyavCurrentUser).Msg("goyavCurrentUser")
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.IndentedJSON(http.StatusOK, goyavCurrentUser)
	}
}
