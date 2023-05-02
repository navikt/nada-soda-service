package main

import (
	"github.com/navikt/nada-soda-service/pkg/api"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	router := api.New(log)
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
