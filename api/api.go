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
	api := a.r.Group("/api")

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
	api.Use(cors)

	api.POST("/register", handlers.Register)
	api.POST("/login", handlers.Login)
	api.GET("/user", middleware.AuthMiddleware, handlers.User)
	api.PUT("/user", middleware.AuthMiddleware, handlers.EditProfile)

	api.GET("/roles", handlers.Roles)
	api.GET("/genres", handlers.Genres)
	api.GET("/tags", handlers.Tags)

	api.GET("/artists", handlers.Artists)
	api.GET("/artists/:artist", handlers.Artist)

	api.GET("/beats", handlers.Beats)
	api.GET("/snippets", handlers.Snippets)
	api.GET("/beats/purchased", middleware.AuthMiddleware, middleware.ClientMiddleware, handlers.PurchasedBeats)
	api.GET("/beats/:beat", handlers.Beat)
	api.PUT("/beats/:beat", handlers.EditBeat)
	api.POST("/beats/:beat/hide", handlers.HideBeat)
	api.POST("/beats/:beat/unhide", handlers.UnhideBeat)
	api.DELETE("/beats/:beat", handlers.DeleteBeat)
	api.PUT("/beats/:beat/licenses/:license", handlers.EditLicense)
	api.POST("/beats/:beat/purchase", middleware.AuthMiddleware, middleware.ClientMiddleware, handlers.PurchaseBeat)
	api.POST("/beats/create", middleware.AuthMiddleware, middleware.ArtistMiddleware, handlers.CreateBeat)

	api.POST("/snippets/create", middleware.AuthMiddleware, middleware.ArtistMiddleware, handlers.CreateSnippet)
	api.POST("/demo/create", middleware.AuthMiddleware, middleware.ArtistMiddleware, handlers.CreateDemo)
}

func (a *API) Start() {
	if err := a.r.Run(a.Addr); err != nil {
		log.Fatal().Err(err).Msg("api start error")
	}
}
