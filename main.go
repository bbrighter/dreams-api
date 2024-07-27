package main

import (
	"github.com/bbrighter/dreams-api/app"
	"github.com/bbrighter/dreams-api/config"
)

func main() {
	cnf, err := config.NewConfig()
	if err != nil {
		return
	}
	app.Run(cnf)
}
