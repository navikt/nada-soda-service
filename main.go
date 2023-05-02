package main

import (
	"context"
	"os"

	"github.com/navikt/nada-soda-service/pkg/api"
	"github.com/navikt/nada-soda-service/pkg/bigquery"
	"github.com/navikt/nada-soda-service/pkg/slack"
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

	slackToken := os.Getenv("SLACK_TOKEN")
	slackClient := slack.New(slackToken, log.WithField("subsystem", "slack"))

	router := api.New(bqClient, slackClient, log)
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
