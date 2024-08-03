package main

import (
	"fmt"

	"github.com/bbrighter/dreams-api/app"
	"github.com/bbrighter/dreams-api/config"
)

func main() {
	cnf, err := config.NewConfig()
	if err != nil {
		fmt.Print(err)
		return
	}
	app.Run(cnf)
}
