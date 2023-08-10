package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"IMDK/models"
	"IMDK/proxies"
)

type ReviewController struct{}

func (rc ReviewController) AddReview(c *gin.Context) {
	var movieModel = models.GetMovieModel()
	var reviewModel = models.GetReviewModel()

	movieID, err := strconv.Atoi(c.Param("movieID"))
	if err != nil {
		log.Println("error on review movieID")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid MovieID", "error": err})
		c.Abort()
		return
	}

	var movie models.Movie
	err = movieModel.Get(movieID, &movie)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Requested MovieID does not exist", "error": err})
		c.Abort()
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["review-audio"]

	for _, file := range files {
		log.Println(file.Filename)
	}

	reviewVoiceFile, reviewVoiceFileHeader, err := c.Request.FormFile("review-audio")
	if err != nil {
		log.Println("Error Parsing FormFile: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error parsing review file", "error": err})
		c.Abort()
		return
	}

	log.Println(reviewVoiceFileHeader.Filename)

	kateb := proxies.GetKateb()
	reviewText, err := kateb.Transcribe(&reviewVoiceFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		c.Abort()
		return
	}

	human := proxies.GetHuman()
	isAppropriate, err := human.Check(reviewText)
	if !isAppropriate {
		c.JSON(http.StatusOK, gin.H{"message": "Your review is violent, therefore cannot be submitted."})
		c.Abort()
		return
	}

	dilmaj := proxies.GetDilmaj()

	frenchTranslation, err := dilmaj.Translate(reviewText, models.LANG_FR)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error Translating your review to french"})
		c.Abort()
		return
	}

	spanishTranslation, err := dilmaj.Translate(reviewText, models.LANG_ES)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error Translating your review to spanish"})
		c.Abort()
		return
	}

	_, err = reviewModel.Add(movieID, reviewText, models.LANG_EN)
	if err != nil {
		return
	}

	_, err = reviewModel.Add(movieID, frenchTranslation, models.LANG_FR)
	if err != nil {
		return
	}

	_, err = reviewModel.Add(movieID, spanishTranslation, models.LANG_ES)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your review was submitted successfully."})
}
