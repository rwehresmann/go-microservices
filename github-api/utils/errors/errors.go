package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ApiError interface {
	Status() int
	Message() string
	Error() string
}

type apiError struct {
	Xstatus  int    `json:"status"`
	Xmessage string `json:"message"`
	Xerror   string `json:"error,omitempty"`
}

func (e *apiError) Status() int {
	return e.Xstatus
}

func (e *apiError) Message() string {
	return e.Xmessage
}

func (e *apiError) Error() string {
	return e.Xerror
}

func NewNotFoundApiError(message string) ApiError {
	return &apiError{
		Xstatus:  http.StatusNotFound,
		Xmessage: message,
	}
}

func NewInternalServerError(message string) ApiError {
	return &apiError{
		Xstatus:  http.StatusInternalServerError,
		Xmessage: message,
	}
}

func NewBadRequestError(message string) ApiError {
	return &apiError{
		Xstatus:  http.StatusBadRequest,
		Xmessage: message,
	}
}

func NewApiError(statusCode int, message string) ApiError {
	return &apiError{
		Xstatus:  statusCode,
		Xmessage: message,
	}
}

func NewApiErrFromBytes(body []byte) (ApiError, error) {
	var result apiError
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.New("Invalid json body.")
	}

	return &result, nil
}
