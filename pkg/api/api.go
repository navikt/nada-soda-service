package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/navikt/nada-soda-service/pkg/bigquery"
	"github.com/sirupsen/logrus"
)

type API struct {
	router   *gin.Engine
	bqClient *bigquery.NadaBigQuery
	log      *logrus.Entry
}

func New(bqClient *bigquery.NadaBigQuery, log *logrus.Logger) *API {
	r := gin.Default()
	a := &API{
		router:   r,
		bqClient: bqClient,
		log:      logrus.WithField("subsystem", "api"),
	}
	a.addSODARouters(r)

	return a
}

func (a *API) Run() error {
	return a.router.Run()
}

func (a *API) addSODARouters(r *gin.Engine) {
	a.router.POST("/soda/new", func(c *gin.Context) {
		sodaBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "no go",
			})
			return
		}

		sodaResults := []bigquery.SodaResult{}
		if err := json.Unmarshal(sodaBytes, &sodaResults); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "no go",
			})
			return
		}

		if err := a.bqClient.StoreSodaResults(c, sodaResults); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "no go",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
}
