package main

import (
	"fmt"

	"github.com/bbrighter/dreams-api/app"
	"github.com/bbrighter/dreams-api/config"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	cnf, err := config.NewConfig(logger)
	if err != nil {
		fmt.Print(err)
		return
	}
	app.Run(cnf, logger)
}
