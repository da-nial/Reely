package server

import (
	"github.com/gin-gonic/gin"

	"IMDK/controllers"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.LoadHTMLGlob("web/templates/*")
	router.Static("/css", "web/assets/css")
	router.Static("/fonts", "web/assets/fonts")
	router.Static("/img", "web/assets/img")
	router.Static("/js", "web/assets/js")

	health := new(controllers.HealthController)
	router.GET("/health", health.Status)

	v1 := router.Group("/v1")
	{
		movieGroup := v1.Group("movies")
		{
			mc := new(controllers.MovieController)
			movieGroup.GET("/:movieID", mc.GetMovie)
			movieGroup.GET("/", mc.GetAllMovies)
		}
		reviewGroup := v1.Group("reviews")
		{
			review := new(controllers.ReviewController)
			reviewGroup.POST("movies/:movieID", review.AddReview)
		}
	}
	return router
}
