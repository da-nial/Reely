package models

import (
	dbUtils "IMDK/db"
	"log"
)

var m = &Movie{}

type Movie struct {
	ID int `pg:",pk"`

	Title       string
	PosterURL   string
	Director    string
	Description string

	Reviews []*Review `pg:"rel:has-many"`
}

func (m Movie) Add(title, posterURL, director, description string) (*Movie, error) {
	db := dbUtils.GetDB()

	movie := &Movie{
		Title:       title,
		PosterURL:   posterURL,
		Director:    director,
		Description: description,
	}
	err := db.Insert(movie)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (m Movie) GetAll() ([]Movie, error) {
	db := dbUtils.GetDB()

	var movies []Movie
	err := db.Model(&movies).Select()
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (m Movie) Get(movieID int, movie *Movie) error {
	db := dbUtils.GetDB()

	err := db.Model(movie).Where("ID = ?", movieID).Select()
	if err != nil {
		log.Println("MovieModel: Error retrieving movie: ", err)
		return err
	}

	return nil
}

func GetMovieModel() *Movie {
	return m
}
