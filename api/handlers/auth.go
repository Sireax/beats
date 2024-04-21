package handlers

import (
	"beats/api/requests"
	"beats/db"
	"beats/db/models"
	"beats/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
	var userExist int
	err := db.DB.Raw("SELECT COUNT(id) from users where email=?", r.Email).Scan(&userExist).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if userExist > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "User already exist",
		})
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(r.Password), 14)
	user := models.User{
		Email:    r.Email,
		Username: r.Username,
		Photo:    r.Photo,
		Password: password,
	}
	db.DB.Create(&user)
	c.JSON(200, user)
}

func Login(c *gin.Context) {
	var r requests.LoginRequest
	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
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

	c.JSON(http.StatusOK, gin.H{"status": "success", "token": token, "user": user})
}

func User(c *gin.Context) {
	user, err := util.ExtractUserFromRequest(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, user)
}

func Logout(c *gin.Context) {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
