package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rwehresmann/go-microservices/github-api/log/zap_logger"
)

var router *gin.Engine

func init() {
	router = gin.Default()
}

func StartApp() {
	zap_logger.Info("Mapping urls", zap_logger.Field("step", "1"), zap_logger.Field("status", "pending"))
	mapUrls()

	zap_logger.Info("Mapping urls", zap_logger.Field("step", "2"), zap_logger.Field("status", "success"))

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
