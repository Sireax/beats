package handlers

import (
	"beats/db"
	"beats/db/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func Roles(c *gin.Context) {
	var roles []*models.Role
	err := db.DB.Find(&roles).Error
	if err != nil {
		log.Error().Err(err).Msg("Error getting roles")
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func Genres(c *gin.Context) {
	var genres []*models.Genre
	err := db.DB.Find(&genres).Error
	if err != nil {
		log.Error().Err(err).Msg("Error getting genres")
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, genres)
}
