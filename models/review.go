package models

import (
	dbUtils "IMDK/db"
	"log"
	"time"
)

var r = &Review{}

type Review struct {
	ID int `pg:",pk"`

	MovieID int
	Movie   *Movie `pg:"rel:has-one"`

	Text      string
	Lang      string
	CreatedAt time.Time `pg:"default:now()"`
}

func (r Review) Add(movieID int, text, lang string) (*Review, error) {
	db := dbUtils.GetDB()

	review := &Review{
		MovieID: movieID,

		Text: text,
		Lang: lang,
	}
	err := db.Insert(review)
	if err != nil {
		return nil, err
	}

	return review, nil

}

func (r Review) Get(movieID int, lang string) ([]Review, error) {
	db := dbUtils.GetDB()

	var reviews []Review

	err := db.Model(&reviews).Where("movie_id = ?", movieID).Where("lang = ?", lang).Select()
	if err != nil {
		log.Println("ReviewModel: Error retrieving review: ", err)
		return nil, err
	}

	return reviews, nil
}

func GetReviewModel() *Review {
	return r
}
