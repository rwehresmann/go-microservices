package domain

import (
	"fmt"
	"net/http"

	"github.com/rwehresmann/go-microservices/mvc/utils"
)

var (
	users = map[int64]*User{
		1: {Id: 1, FirstName: "Geralt", LastName: "of Rivia", Email: "gr@thewitcher.com"},
	}
)

func GetUser(userId int64) (*User, *utils.ApplicationError) {
	if user := users[userId]; user != nil {
		return user, nil
	}

	return nil, &utils.ApplicationError{
		Message:    fmt.Sprintf("User %v was not found", userId),
		StatusCode: http.StatusNotFound,
		Code:       "not_found",
	}
}
