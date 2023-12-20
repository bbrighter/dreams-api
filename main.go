package main

import (
	"log"
)

// var ipAddress string

// func init() {
// 	flag.StringVar(&ipAddress, "ipAddress", "192.168.178.133", "IP address to run")
// 	flag.Parse()
// }

func main() {
	repo := InitRepo("dreams.sqlite")
	con := InitController(repo)

	if err := Migration(repo.db); err != nil {
		log.Fatal(err)
	}
	r := SetupRouter(con)
	r.Run("localhost" + ":5000")
}
