package repository

import "hoc-gin/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	FindById(user *models.User, id int) error
}
