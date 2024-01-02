package main

import (
	"flag"
	"log"
)

var ipAddress string

func main() {
	flag.StringVar(&ipAddress, "ipAddress", "localhost", "IP address to run")
	flag.Parse()

	repo := InitRepo("dreams.sqlite")
	con := InitController(repo)

	if err := Migration(repo.db); err != nil {
		log.Fatal(err)
	}
	r := SetupRouter(con)

	r.Run(ipAddress + ":5005")
}
