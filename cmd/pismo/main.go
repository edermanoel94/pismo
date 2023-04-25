package main

import (
	"github.com/edermanoel94/pismo/internal/api"
	"github.com/edermanoel94/pismo/internal/infra/config"
	"github.com/edermanoel94/pismo/internal/infra/logger"
	"github.com/sirupsen/logrus"
)

func main() {

	config.Init()

	logger.Init(logrus.StandardLogger())

	if err := api.Start(); err != nil {
		logrus.Fatal(err)
	}
}
