package main

import (
	"log"

	"github.com/bbrighter/dreams-api/controller"
	"github.com/bbrighter/dreams-api/store"
)

func main() {
	var dbName string = "dreams.sqlite"

	var repo store.Repo = store.InitRepo(dbName)
	var con controller.Controller = controller.InitController(repo)

	if err := store.Migration(repo); err != nil {
		log.Fatal(err)
	}

	con.RunRouter()
}
