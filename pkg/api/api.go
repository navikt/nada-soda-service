package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type API struct {
	router *gin.Engine
	log    *logrus.Entry
}

func New(log *logrus.Logger) *API {
	r := gin.Default()
	a := &API{
		router: r,
		log:    logrus.WithField("subsystem", "api"),
	}
	a.addSODARouters(r)

	return a
}

func (a *API) Run() error {
	return a.router.Run()
}

func (a *API) addSODARouters(r *gin.Engine) {
	a.router.POST("/soda/new", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
}
