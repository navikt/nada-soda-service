package main

import (
	"context"
	"os"

	"github.com/navikt/nada-soda-service/pkg/api"
	"github.com/navikt/nada-soda-service/pkg/bigquery"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	log := logrus.New()

	project := os.Getenv("GCP_TEAM_PROJECT_ID")
	dataset := os.Getenv("BIGQUERY_DATASET")
	table := os.Getenv("BIGQUERY_TABLE")

	bqClient, err := bigquery.New(ctx, project, dataset, table, log.WithField("subsystem", "bigquery"))
	if err != nil {
		log.Fatal(err)
	}

	router := api.New(bqClient, log)
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
