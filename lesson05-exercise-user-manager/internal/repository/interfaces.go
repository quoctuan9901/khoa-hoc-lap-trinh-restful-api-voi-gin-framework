package repository

import (
	"user-management-api/internal/models"
)

type UserRepository interface {
	FindAll() ([]models.User, error)
	Create(user models.User) error
	FindByUUID(uuid string) (models.User, bool)
	Update(uuid string, user models.User) error
	Delete(uuid string) error
	FindByEmail(email string) (models.User, bool)
}
