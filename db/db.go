package db

import (
	"github.com/go-pg/pg"
	"log"
	"os"

	"IMDK/config"
)

var db *pg.DB

func Init() {
	c := config.GetConfig()

	opts := &pg.Options{
		User:     c.GetString("db.user"),
		Password: c.GetString("db.password"),
		Addr:     c.GetString("db.host"),
		Database: c.GetString("db.name"),
	}
	db = pg.Connect(opts)
	if db == nil {
		log.Printf("Failed to connect")
		os.Exit(1)
	}
	log.Printf("Connected to db")
}

func GetDB() *pg.DB {
	return db
}
