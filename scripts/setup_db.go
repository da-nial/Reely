package scripts

import (
	"fmt"
	"github.com/go-pg/pg/orm"
	"log"

	dbUtils "IMDK/db"
	"IMDK/models"
)

func SetupDB() {
	db := dbUtils.GetDB()

	for _, model := range []interface{}{&models.Movie{}, &models.Review{}} {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			log.Panicln(err)
		}
	}

	m := models.GetMovieModel()
	_, err := m.Add("James Bond - No Time to Die",
		"https://s6.uupload.ir/files/james-bond_4e7s.jpeg",
		"Sam Mendez",
		"James Bond is enjoying a tranquil life in Jamaica after leaving active service. However, his peace is short-lived as his old CIA friend, Felix Leiter, shows up and asks for help. The mission to rescue a kidnapped scientist turns out to be far more treacherous than expected, leading Bond on the trail of a mysterious villain who's armed with a dangerous new technology.\n")
	if err != nil {
		fmt.Println("SetupDB: Error while creating mock data: ", err)
		return
	}

	_, err = m.Add("Gone Girl",
		"https://s6.uupload.ir/files/gone-girl_gzjj.jpeg",
		"David Fincher",
		"Nick Dunne discovers that the entire media focus has shifted on him when his wife, Amy Dunne, mysteriously disappears on the day of their fifth wedding anniversary.\n",
	)
	if err != nil {
		fmt.Println("SetupDB: Error while creating mock data: ", err)
		return
	}

	_, err = m.Add("Under the Skin",
		"https://s6.uupload.ir/files/under-the-skin_zeib.jpeg",
		"Jonathan Glazer",
		"Disguising itself as a human female, an extraterrestrial drives around Scotland attempting to lure unsuspecting men into her van. Once there, she seduces and sends them into another dimension where they are nothing more than meat.\n")
	if err != nil {
		fmt.Println("SetupDB: Error while creating mock data: ", err)
		return
	}

	movies, err := m.GetAll()
	fmt.Println("All Movies in DB: ", movies)
}
