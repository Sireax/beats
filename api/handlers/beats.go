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
	var r requests.UpdateBeatRequest
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

func EditBeat(c *gin.Context) {
	var r requests.UpdateBeatRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	beatID := c.Param("beat")
	err := db.DB.
		Exec("UPDATE beats SET title = ?, photo = ?, link = ?, genre_id = ? WHERE id = ?",
			r.Title, r.Photo, r.Link, r.GenreID, beatID).Error
	if err != nil {
		log.Error().Err(err).Msg("error updating beat")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var beat *models.Beat
	db.DB.Raw("SELECT * FROM beats WHERE id = ?", beatID).Scan(&beat)
	c.JSON(200, beat)
}

func DeleteBeat(c *gin.Context) {
	beatID := c.Param("beat")

	var purchasedCount int
	err := db.DB.Raw("SELECT count(*) FROM purchases WHERE beat_id = ?", beatID).
		Scan(&purchasedCount).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Error().Err(err).Msg("error getting purchased count")
		return
	}
	if purchasedCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete already purchased beat"})
		return
	}

	err = db.DB.Exec("DELETE FROM snippets WHERE beat_id = ?", beatID).Error
	if err != nil {
		log.Error().Err(err).Msg("error deleting snippets")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = db.DB.Exec("DELETE FROM demos WHERE beat_id = ?", beatID).Error
	if err != nil {
		log.Error().Err(err).Msg("error deleting demos")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = db.DB.Exec("DELETE FROM song_tags WHERE beat_id = ?", beatID).Error
	if err != nil {
		log.Error().Err(err).Msg("error deleting song_tags")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = db.DB.Exec("DELETE FROM licenses WHERE beat_id = ?", beatID).Error
	if err != nil {
		log.Error().Err(err).Msg("error deleting licenses")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = db.DB.Exec("DELETE FROM reviews WHERE beat_id = ?", beatID).Error
	if err != nil {
		log.Error().Err(err).Msg("error deleting reviews")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = db.DB.Exec("DELETE FROM beats WHERE id = ?", beatID).Error
	if err != nil {
		log.Error().Err(err).Msg("error deleting beat")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func HideBeat(c *gin.Context) {
	beatID := c.Param("beat")

	err := db.DB.Exec("UPDATE beats SET is_hide = true WHERE id = ?", beatID).Error
	if err != nil {
		log.Error().Err(err).Msg("error hiding beat")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{})
}

func UnhideBeat(c *gin.Context) {
	beatID := c.Param("beat")

	err := db.DB.Exec("UPDATE beats SET is_hide = false WHERE id = ?", beatID).Error
	if err != nil {
		log.Error().Err(err).Msg("error unhiding beat")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{})
}

func Beats(c *gin.Context) {
	beats := make([]*models.Beat, 0)

	err := db.DB.Raw("SELECT * FROM beats WHERE is_hide = false ORDER BY id DESC").
		Scan(&beats).Error
	if err != nil {
		log.Error().Err(err).Msg("error getting all beats")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	for _, beat := range beats {
		db.DB.Raw("select * from users where id = ? limit 1", beat.UserID).Scan(&beat.User)
		db.DB.Raw("select * from genres where id = ? limit 1", beat.GenreID).Scan(&beat.Genre)
		db.DB.Raw("select * from tags join public.song_tags st on tags.id = st.tag_id WHERE st.beat_id = ?", beat.ID).Scan(&beat.Tags)
	}

	c.JSON(http.StatusOK, beats)
}

func Snippets(c *gin.Context) {
	snippets := make([]*models.Snippet, 0)
	err := db.DB.Raw("SELECT snippets.* FROM snippets join beats b on b.id = snippets.beat_id WHERE b.is_hide = false ORDER BY id DESC").Error
	if err != nil {
		log.Error().Err(err).Msg("error getting all snippets")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	for _, snippet := range snippets {
		db.DB.Raw("SELECT * FROM beats WHERE id = ?", snippet.BeatID).Scan(&snippet.Beat)
	}

	c.JSON(200, &snippets)
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
	db.DB.
		Raw("select * from licenses where beat_id = ?", beat.ID).
		Scan(&beat.Licenses)
	for _, license := range beat.Licenses {
		db.DB.Raw("select * from license_types where id = ?", license.LicenseTypeID).
			Scan(&license.LicenseType)
	}
	db.DB.
		Raw("select * from users where id = ?", beat.UserID).
		Scan(&beat.User)
	db.DB.
		Raw("select * from snippets where beat_id = ?", beat.ID).
		Scan(&beat.Snippets)
	db.DB.
		Raw("select * from demos where beat_id = ?", beat.ID).
		Scan(&beat.Demos)
	db.DB.
		Raw("select * from tags join public.song_tags st on tags.id = st.tag_id WHERE st.beat_id = ?", beat.ID).
		Scan(&beat.Tags)
	db.DB.
		Raw("select * from genres where id = ? limit 1", beat.GenreID).
		Scan(&beat.Genre)

	c.JSON(http.StatusOK, beat)
}

func EditLicense(c *gin.Context) {
	var r requests.EditLicenseRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	beatID := c.Param("beat")
	licenseID := c.Param("license")

	var beat *models.Beat
	err := db.DB.Raw("SELECT * FROM beats WHERE id = ?", beatID).Scan(&beat).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Error().Err(err).Msg("error getting beat")
		return
	}
	if beat == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var license *models.License
	err = db.DB.Raw("SELECT * FROM licenses WHERE id = ?", licenseID).Scan(&license).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Error().Err(err).Msg("error getting license")
		return
	}
	if license == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err = db.DB.
		Exec("UPDATE licenses SET price = ?, rental_time = ?, is_active = ? WHERE id = ?",
			r.Price, r.RentalTime, r.IsActive, licenseID).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Error().Err(err).Msg("error updating license")
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func PurchaseBeat(c *gin.Context) {
	var r requests.PurchaseBeatRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	var license *models.License
	err = db.DB.Raw("SELECT * FROM licenses WHERE id = ? LIMIT 1", r.LicenseID).Scan(&license).Error
	if err != nil {
		log.Error().Err(err).Msg("error getting license")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if license == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err = db.DB.
		Exec("INSERT INTO purchases (user_id, beat_id, license_id) VALUES (?, ?, ?)", user.ID, beatID, license.ID).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func PurchasedBeats(c *gin.Context) {
	user, err := util.ExtractUserFromRequest(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	beats := make([]*models.Beat, 0)
	err = db.DB.
		Raw("SELECT beats.* FROM beats join public.purchases p on beats.id = p.beat_id where p.user_id = ? order by created_at desc", user.ID).
		Scan(&beats).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Error().Err(err).Msg("error getting purchased beats")
		return
	}
	for _, beat := range beats {
		db.DB.Raw("select * from users where id = ?", beat.UserID).Scan(&beat.User)
		db.DB.Raw("select * from licenses where beat_id = ?", beat.ID).Scan(&beat.Licenses)
	}

	c.JSON(http.StatusOK, beats)
}
