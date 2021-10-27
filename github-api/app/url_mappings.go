package app

import (
	"github.com/rwehresmann/go-microservices/github-api/controllers/polo"
	repositories "github.com/rwehresmann/go-microservices/github-api/controllers/repositories"
)

func mapUrls() {
	router.GET("/marco", polo.Polo)
	router.POST("/repositories", repositories.CreateRepo)
}
