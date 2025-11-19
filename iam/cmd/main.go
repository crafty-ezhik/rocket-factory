package main

import (
	"log"

	"github.com/crafty-ezhik/rocket-factory/iam/internal/config"
)

const configPath = "../deploy/compose/iam/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		log.Println(err)
		return
	}
}
