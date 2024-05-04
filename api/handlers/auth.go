package handlers

import (
	"beats/api/requests"
	"beats/db"
	"beats/db/models"
	"beats/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

const SecretKey = "secret"

func Register(c *gin.Context) {
	var r requests.RegisterRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if r.Username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
	}
	if r.Email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
	}
	if r.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
	}
	if r.RoleID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "RoleID is required"})
	}
	var userExist int
	err := db.DB.Raw("SELECT COUNT(id) from users where email=?", r.Email).Scan(&userExist).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if userExist > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Пользователь уже существует",
		})
		return
	}

	var roleExists int
	err = db.DB.Raw("SELECT COUNT(id) from roles where id=?", r.RoleID).Scan(&roleExists).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if roleExists == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Роль не существует",
		})
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(r.Password), 14)
	user := models.User{
		Email:    r.Email,
		Username: r.Username,
		Password: password,
		RoleID:   r.RoleID,
	}
	db.DB.Create(&user)

	c.JSON(200, gin.H{})
}

func Login(c *gin.Context) {
	var r requests.LoginRequest
	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	if r.Email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
	}
	if r.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
	}

	var user models.User
	db.DB.Where("email = ?", r.Email).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(r.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect password"})
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})
	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "token": token})
}

func User(c *gin.Context) {
	user, err := util.ExtractUserFromRequest(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, user)
}

func EditProfile(c *gin.Context) {
	var r requests.EditProfileRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	user, err := util.ExtractUserFromRequest(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if r.Password != "" {
		password, _ := bcrypt.GenerateFromPassword([]byte(r.Password), 14)
		err = db.DB.
			Exec("UPDATE users SET username = ?, email = ?, password = ? WHERE id = ?",
				r.Username, r.Email, password, user.ID).Error
	} else {
		err = db.DB.
			Exec("UPDATE users SET username = ?, email = ? WHERE id = ?",
				r.Username, r.Email, user.ID).Error
	}
	if err != nil {
		log.Error().Err(err).Msg("error updating user")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// обновляем модель пользователя
	db.DB.Raw("SELECT * FROM users WHERE id = ?", user.ID).Scan(&user)

	c.JSON(http.StatusOK, user)
}
