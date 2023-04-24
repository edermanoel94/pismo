package main

import (
	"github.com/edermanoel94/pismo/internal/api"
	"github.com/edermanoel94/pismo/internal/infra/config"
	"log"
)

func main() {

	config.Init()

	if err := api.Start(); err != nil {
		log.Fatal(err)
	}
}
