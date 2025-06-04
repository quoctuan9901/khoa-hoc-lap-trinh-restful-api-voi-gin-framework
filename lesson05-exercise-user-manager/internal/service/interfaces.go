package service

import "user-management-api/internal/models"

type UserService interface {
	GetAllUsers() ([]models.User, error)
	CreateUser(user models.User) (models.User, error)
	GetUserByUUID(uuid string) (models.User, error)
	UpdateUser()
	DeleteUser()
}
