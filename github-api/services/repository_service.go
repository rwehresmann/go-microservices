package services

import (
	"net/http"
	"sync"

	"github.com/rwehresmann/go-microservices/github-api/config"
	"github.com/rwehresmann/go-microservices/github-api/domain/github"
	"github.com/rwehresmann/go-microservices/github-api/domain/repositories"
	"github.com/rwehresmann/go-microservices/github-api/log/zap_logger"
	github_provider "github.com/rwehresmann/go-microservices/github-api/providers"
	"github.com/rwehresmann/go-microservices/github-api/utils/errors"
)

type reposService struct{}

type reposServiceInterface interface {
	CreateRepo(clientId string, request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
	CreateRepos(request []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError)
}

var RepositoryService reposServiceInterface

func init() {
	RepositoryService = &reposService{}
}

func (s *reposService) CreateRepo(clientId string, input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	request := github.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Private:     false,
	}

	zap_logger.Info(
		"About to send request to external API",
		zap_logger.Field("client_id", clientId),
		zap_logger.Field("satus", "pending"),
	)
	response, err := github_provider.CreateRepo(config.GetGithubAccessToken(), request)
	if err != nil {
		zap_logger.Error(
			"Response from external API",
			zap_logger.Field("client_id", clientId),
			zap_logger.Field("satus", "error"),
		)

		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}

	zap_logger.Info(
		"Response from external API",
		zap_logger.Field("client_id", clientId),
		zap_logger.Field("satus", "success"),
	)
	result := repositories.CreateRepoResponse{
		Id:    response.Id,
		Name:  response.Name,
		Owner: response.Owner.Login,
	}

	return &result, nil
}

func (s *reposService) CreateRepos(requests []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError) {
	chInput := make(chan repositories.CreateReposResult)
	chOutput := make(chan repositories.CreateReposResponse)
	defer close(chOutput)

	var wg sync.WaitGroup
	go s.handleRepoResults(&wg, chInput, chOutput)

	for _, request := range requests {
		wg.Add(1)
		go s.createRepoConcurrent(request, chInput)
	}

	wg.Wait()
	close(chInput)

	result := <-chOutput

	successCreations := 0
	for _, current := range result.Results {
		if current.Response != nil {
			successCreations++
		}
	}

	if successCreations == 0 {
		result.StatusCode = result.Results[0].Error.Status()
	} else if successCreations == len(result.Results) {
		result.StatusCode = http.StatusCreated
	} else {
		result.StatusCode = http.StatusPartialContent
	}

	return result, nil
}

func (s *reposService) handleRepoResults(wg *sync.WaitGroup, inputCh chan repositories.CreateReposResult, outputCh chan repositories.CreateReposResponse) {
	var results repositories.CreateReposResponse

	for result := range inputCh {
		results.Results = append(results.Results, repositories.CreateReposResult{
			Response: result.Response,
			Error:    result.Error,
		})
		wg.Done()
	}

	outputCh <- results
}

func (s *reposService) createRepoConcurrent(input repositories.CreateRepoRequest, ch chan repositories.CreateReposResult) {
	if err := input.Validate(); err != nil {
		ch <- repositories.CreateReposResult{Error: err}
		return
	}

	result, err := s.CreateRepo("", input)

	if err != nil {
		ch <- repositories.CreateReposResult{Error: err}
		return
	}

	ch <- repositories.CreateReposResult{Response: result}
}
