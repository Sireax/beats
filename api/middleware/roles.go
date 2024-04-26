package middleware

import (
	"beats/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ArtistMiddleware(c *gin.Context) {
	user, err := util.ExtractUserFromRequest(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if user.Role != nil && !user.Role.IsArtist() {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}

func ClientMiddleware(c *gin.Context) {
	user, err := util.ExtractUserFromRequest(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if user.Role != nil && !user.Role.IsClient() {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}
