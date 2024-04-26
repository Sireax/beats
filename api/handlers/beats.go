package handlers

import (
	"beats/api/requests"
	"beats/db"
	"beats/db/models"
	"beats/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func CreateSnippet(c *gin.Context) {
	var r requests.CreateSnippetRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var beat *models.Beat
	err := db.DB.Where("id = ?", r.BeatID).First(&beat).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if beat == nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	snippet := models.Snippet{
		Start:  r.Start,
		End:    r.End,
		BeatID: r.BeatID,
	}
	err = db.DB.Create(&snippet).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusCreated, snippet)
}

func CreateDemo(c *gin.Context) {
	var r requests.CreateDemoRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var beat *models.Beat
	err := db.DB.Where("id = ?", r.BeatID).First(&beat).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if beat == nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	demo := models.Demo{
		Start:  r.Start,
		End:    r.End,
		BeatID: r.BeatID,
	}
	err = db.DB.Create(&demo).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusCreated, demo)
}

func CreateBeat(c *gin.Context) {
	var r requests.CreateBeatRequest
	err := c.ShouldBindJSON(&r)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := util.ExtractUserFromRequest(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var genre *models.Genre
	err = db.DB.Where("id = ?", r.GenreID).First(&genre).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if genre == nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	beat := &models.Beat{
		ReleaseDate: time.Now(),
		Photo:       r.Photo,
		Title:       r.Title,
		Link:        r.Link,
		GenreID:     genre.ID,
		UserID:      user.ID,
		IsHide:      false,
	}
	err = db.DB.Create(&beat).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, beat)
}

func Beats(c *gin.Context) {
	user, err := util.ExtractUserFromRequest(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var beats []*models.Beat

	switch user.RoleID {
	case models.ArtistRoleID:
		err = db.DB.Raw("SELECT * FROM beats WHERE user_id = ? ORDER BY id DESC", user.ID).
			Scan(&beats).Error
		if err != nil {
			log.Error().Err(err).Msg("error getting user beats")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	case models.ClientRoleID:
		err = db.DB.Raw("SELECT * FROM beats ORDER BY id DESC").
			Scan(&beats).Error
		if err != nil {
			log.Error().Err(err).Msg("error getting all beats")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	// иначе отдает null вместо пустого массива
	if len(beats) == 0 {
		beats = []*models.Beat{}
	}

	c.JSON(http.StatusOK, beats)
}
