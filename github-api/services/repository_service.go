package services

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/rwehresmann/go-microservices/github-api/config"
	"github.com/rwehresmann/go-microservices/github-api/domain/github"
	"github.com/rwehresmann/go-microservices/github-api/domain/repositories"
	github_provider "github.com/rwehresmann/go-microservices/github-api/providers"
	"github.com/rwehresmann/go-microservices/github-api/utils/errors"
)

type reposService struct{}

type reposServiceInterface interface {
	CreateRepo(request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
	CreateRepos(request []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError)
}

var RepositoryService reposServiceInterface

func init() {
	RepositoryService = &reposService{}
}

func (s *reposService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	request := github.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Private:     false,
	}

	fmt.Println("token")
	fmt.Println(os.Getenv("SECRET_GITHUB_ACCESS_TOKEN"))

	response, err := github_provider.CreateRepo(config.GetGithubAccessToken(), request)
	if err != nil {
		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}

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

	result, err := s.CreateRepo(input)

	if err != nil {
		ch <- repositories.CreateReposResult{Error: err}
		return
	}

	ch <- repositories.CreateReposResult{Response: result}
}
