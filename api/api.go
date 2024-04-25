package api

import (
	"beats/api/handlers"
	"beats/api/middleware"
	"github.com/gin-gonic/gin"
)
import "github.com/rs/zerolog/log"

type API struct {
	Addr string
	r    *gin.Engine
}

func NewAPI(addr string) *API {
	r := gin.Default()
	s := &API{
		Addr: addr,
		r:    r,
	}
	s.UseRoutes()
	return s
}

func (a *API) UseRoutes() {
	cors := func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
	a.r.Use(cors)

	a.r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	a.r.POST("/api/register", handlers.Register)
	a.r.POST("/api/login", handlers.Login)
	a.r.GET("/api/user", middleware.AuthMiddleware, handlers.User)

	a.r.GET("/api/roles", handlers.Roles)

	a.r.POST("/api/beats/create", middleware.AuthMiddleware, handlers.CreateBeat)
	a.r.POST("/api/snippets/create", middleware.AuthMiddleware, handlers.CreateSnippet)
	a.r.POST("/api/demo/create", middleware.AuthMiddleware, handlers.CreateDemo)
}

func (a *API) Start() {
	if err := a.r.Run(a.Addr); err != nil {
		log.Fatal().Err(err).Msg("api start error")
	}
}
