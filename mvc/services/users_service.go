package services

import (
	"github.com/rwehresmann/go-microservices/mvc/domain"
	"github.com/rwehresmann/go-microservices/mvc/utils"
)

func GetUser(userId int64) (*domain.User, *utils.ApplicationError) {
	return domain.GetUser(userId)
}
