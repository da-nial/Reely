package main

import (
	"IMDK/config"
	"IMDK/db"
	"IMDK/proxies"
	"IMDK/server"
	"flag"
	"fmt"
	"os"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()

	config.Init(*environment)

	db.Init()
	//scripts.SetupDB()
	//return

	proxies.InitKateb()
	proxies.InitHuman()
	proxies.InitDilmaj()

	server.Init()
}
