package service

import "user-management-api/internal/models"

type UserService interface {
	GetAllUsers()
	CreateUser(user models.User) (models.User, error)
	GetUserByUUID()
	UpdateUser()
	DeleteUser()
}
