package repositories

import (
	"strings"

	"github.com/rwehresmann/go-microservices/github-api/utils/errors"
)

type CreateRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *CreateRepoRequest) Validate() errors.ApiError {
	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return errors.NewBadRequestError("Invalid repository name.")
	}

	return nil
}

type CreateRepoResponse struct {
	Id    int64  `json:"id"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type CreateReposResponse struct {
	StatusCode int                 `json:"status"`
	Results    []CreateReposResult `json:"results"`
}

type CreateReposResult struct {
	Response *CreateRepoResponse `json:"repo"`
	Error    errors.ApiError     `json:"error"`
}
