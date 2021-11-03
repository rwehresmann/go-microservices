package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rwehresmann/go-microservices/github-api/log/logrus_logger"
)

var router *gin.Engine

func init() {
	router = gin.Default()
}

func StartApp() {
	logrus_logger.Info("Mapping urls", "step:1", "status:pending")
	mapUrls()

	logrus_logger.Info("Mapping urls", "step:2", "status:success")

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
