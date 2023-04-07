package main

import (
	"log"

	config "github.com/SethukumarJ/CashierX-Auth-Service/pkg/config"
	di "github.com/SethukumarJ/CashierX-Auth-Service/pkg/di"
)

func main() {
	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	server, diErr := di.InitializeAPI(config)
	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start()
	}
}
