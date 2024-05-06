package handlers

import (
	"beats/db"
	"beats/db/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func Roles(c *gin.Context) {
	roles := make([]*models.Role, 0)
	err := db.DB.Find(&roles).Error
	if err != nil {
		log.Error().Err(err).Msg("Error getting roles")
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func Genres(c *gin.Context) {
	genres := make([]models.Genre, 0)
	err := db.DB.Find(&genres).Error
	if err != nil {
		log.Error().Err(err).Msg("error getting genres")
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, genres)
}

func Tags(c *gin.Context) {
	tags := make([]*models.Tag, 0)
	err := db.DB.Raw("SELECT * FROM tags ORDER BY id DESC").Scan(&tags).Error
	if err != nil {
		log.Error().Err(err).Msg("error getting tags")
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func Artists(c *gin.Context) {
	var artists []*models.User
	err := db.DB.
		Raw("select users.* from users join roles on users.role_id = roles.id where roles.type = ?", models.ArtistRoleType).
		Scan(&artists).Error
	if err != nil {
		log.Error().Err(err).Msg("error getting artists")
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, artists)
}

func Artist(c *gin.Context) {
	artistID := c.Param("artist")
	var artist *models.User
	err := db.DB.Raw("select * from users where id = ? limit 1", artistID).
		Scan(&artist).Error
	if err != nil {
		log.Error().Err(err).Msg("error getting artist")
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if artist == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	beats := make([]*models.Beat, 0)
	err = db.DB.Raw("select * from beats where user_id = ? order by id desc", artist.ID).
		Scan(&beats).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}
	for _, beat := range beats {
		db.DB.Raw("select * from genres where id = ? limit 1", beat.GenreID).Scan(&beat.Genre)
		db.DB.Raw("select * from demos where beat_id = ? order by id desc", beat.ID).Scan(&beat.Demos)
		db.DB.Raw("select * from snippets where beat_id = ? order by id desc", beat.ID).Scan(&beat.Snippets)
	}
	artist.Beats = beats

	c.JSON(http.StatusOK, artist)
}
