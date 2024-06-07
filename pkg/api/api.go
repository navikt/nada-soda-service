package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/navikt/nada-soda-service/pkg/bigquery"
	"github.com/navikt/nada-soda-service/pkg/models"
	"github.com/navikt/nada-soda-service/pkg/slack"
	"github.com/sirupsen/logrus"
)

type API struct {
	router      *gin.Engine
	bqClient    *bigquery.BigQueryClient
	slackClient *slack.SlackClient
	log         *logrus.Entry
}

func New(bqClient *bigquery.BigQueryClient, slackClient *slack.SlackClient, log *logrus.Logger) *API {
	r := gin.Default()
	a := &API{
		router:      r,
		bqClient:    bqClient,
		slackClient: slackClient,
		log:         logrus.WithField("subsystem", "api"),
	}
	a.addSODARouters()

	return a
}

func (a *API) Run() error {
	return a.router.Run()
}

func (a *API) addSODARouters() {
	a.router.POST("/soda/new", func(c *gin.Context) {
		sodaBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "reading request body",
			})
			return
		}

		sodaResults := models.SodaTest{}
		if err := json.Unmarshal(sodaBytes, &sodaResults); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "unmarshal request body",
			})
			return
		}

		if err := a.processSodaResults(c, sodaResults); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error processing soda results",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
}

func (a *API) processSodaResults(ctx context.Context, sodaTest models.SodaTest) error {
	if err := a.slackClient.Notify(sodaTest); err != nil {
		a.log.Errorf("sending slack notification: %v", err)
		return err
	}

	if err := a.bqClient.StoreSodaResults(ctx, sodaTest); err != nil {
		a.log.Errorf("storing soda results %v", err)
		return err
	}

	return nil
}
