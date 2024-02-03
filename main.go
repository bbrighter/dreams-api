package main

import (
	"flag"
	"log"

	"github.com/bbrighter/dreams-api/controller"
	"github.com/bbrighter/dreams-api/store"
)

var ipAddress string

func main() {
	flag.StringVar(&ipAddress, "ipAddress", "localhost", "IP address to run")
	flag.Parse()

	repo := store.InitRepo("dreams.sqlite")
	con := controller.InitController(repo)

	if err := store.Migration(repo); err != nil {
		log.Fatal(err)
	}

	r := controller.SetupRouter(con)
	store.Rollback(repo)

	r.Run(ipAddress + ":5005")
}
