package handlers

import (
	"beats/api/requests"
	"beats/db"
	"beats/db/models"
	"beats/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"
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
	err = db.DB.Clauses(clause.Returning{}).Create(&beat).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	for _, lType := range models.LicenseTypesList() {
		err = db.DB.Debug().
			Exec(
				"INSERT INTO licenses (price, rental_time, license_type_id, beat_id, is_active) VALUES (?, ?, (select id from license_types where type = ?), ?, ?)",
				0, "", lType, beat.ID, false,
			).Error
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			log.Error().Err(err).Msg("error inserting license type")
			return
		}
	}

	c.JSON(http.StatusCreated, beat)
}

func Beats(c *gin.Context) {
	beats := make([]*models.Beat, 0)

	err := db.DB.Raw("SELECT * FROM beats ORDER BY id DESC").
		Scan(&beats).Error
	if err != nil {
		log.Error().Err(err).Msg("error getting all beats")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, beats)
}

func Beat(c *gin.Context) {
	beatId := c.Param("beat")

	var beat *models.Beat
	err := db.DB.Raw("select * from beats where id = ?", beatId).Scan(&beat).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Error().Err(err).Msg("error getting beat")
		return
	}
	if beat == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	err = db.DB.
		Raw("select * from licenses where beat_id = ?", beat.ID).
		Scan(&beat.Licenses).Error

	c.JSON(http.StatusOK, beat)
}

func PurchaseBeat(c *gin.Context) {
	beatID := c.Param("beat")
	if beatID == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	user, err := util.ExtractUserFromRequest(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var beat *models.Beat
	err = db.DB.Where("id = ?", beatID).First(&beat).Error
	if err != nil {
		log.Error().Err(err).Msg("error getting beat")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if beat == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err = db.DB.
		Exec("INSERT INTO purchases (user_id, beat_id) VALUES (?, ?)", user.ID, beatID).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func PurchasedBeats(c *gin.Context) {

}
