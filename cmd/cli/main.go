package main

import (
	"log"

	"github.com/schnit7el/goAnsible/internal/config"
	"github.com/schnit7el/goAnsible/internal/deployer"
)

func main() {
	config, err := config.LoadConfig("example.yml")
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	err = deployer.RunDeployment(config)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

}
