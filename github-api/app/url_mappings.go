package app

import (
	"github.com/rwehresmann/go-microservices/github-api/controllers/polo"
	repositories "github.com/rwehresmann/go-microservices/github-api/controllers/repositories"
)

func mapUrls() {
	router.GET("/marco", polo.Marco)
	router.POST("/repository", repositories.CreateRepo)
	router.POST("/repositories", repositories.CreateRepos)
}
