package controllers

import (
	"log"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"

	"IMDK/models"
)

type MovieController struct{}

func (mc MovieController) GetAllMovies(c *gin.Context) {
	m := models.GetMovieModel()

	movies, err := m.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving movies", "error": err})
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "movie_list.html", gin.H{"movies": movies})
}

func (mc MovieController) GetMovie(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("movieID"))
	lang := c.DefaultQuery("lang", "en")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid MovieID", "error": err})
		c.Abort()
		return
	}

	m := models.GetMovieModel()
	var movie models.Movie
	err = m.Get(movieID, &movie)
	if err != nil {
		log.Println("Err: ", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Requested MovieID does not exist", "error": err})
		c.Abort()
		return
	}

	r := models.GetReviewModel()
	reviews, err := r.Get(movieID, lang)
	if err != nil {
		log.Println("Err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong while retrieving comments for the requested move", "error": err})
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "movie.html", gin.H{
		"movie":   movie,
		"reviews": reviews,
	})
}
