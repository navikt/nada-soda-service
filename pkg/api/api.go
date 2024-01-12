package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/navikt/nada-soda-service/pkg/bigquery"
	"github.com/navikt/nada-soda-service/pkg/models"
	"github.com/navikt/nada-soda-service/pkg/slack"
	"github.com/sirupsen/logrus"
)

type API struct {
	router   *gin.Engine
	bigQuery *bigquery.BigQueryClient
	slack    *slack.Client
	log      *logrus.Entry
}

func New(bqClient *bigquery.BigQueryClient, slackClient *slack.Client, log *logrus.Entry) *API {
	api := &API{
		router:   gin.Default(),
		bigQuery: bqClient,
		slack:    slackClient,
		log:      log,
	}
	api.addSodaRouters()

	return api
}

func (a *API) Run() error {
	return a.router.Run()
}

func (a *API) addSodaRouters() {
	a.router.POST("/soda/new", func(c *gin.Context) {
		sodaBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "reading request body",
				"error":   err.Error(),
			})
			return
		}

		sodaResults := models.SodaTest{}
		if err := json.Unmarshal(sodaBytes, &sodaResults); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "unmarshal request body",
				"error":   err.Error(),
			})
			return
		}

		if err := a.processSodaResults(c, sodaResults); err != nil {
			a.log.WithError(err).Error("processing Soda results")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error processing Soda results",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
}

func (a *API) processSodaResults(ctx context.Context, sodaTest models.SodaTest) error {
	if err := a.slack.NotifyOnDiscrepancies(sodaTest); err != nil {
		return fmt.Errorf("sending Slack notification: %w", err)
	}

	if err := a.bigQuery.StoreSodaResults(ctx, sodaTest); err != nil {
		return fmt.Errorf("storing Soda results: %w", err)
	}

	return nil
}
