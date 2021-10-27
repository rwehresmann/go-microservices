package repositories

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rwehresmann/go-microservices/github-api/domain/repositories"
	"github.com/rwehresmann/go-microservices/github-api/services"
	"github.com/rwehresmann/go-microservices/github-api/utils/errors"
)

func CreateRepo(c *gin.Context) {
	var request repositories.CreateRepoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := errors.NewBadRequestError("Invalid json body.")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	result, err := services.RepositoryService.CreateRepo(request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, result)
}
