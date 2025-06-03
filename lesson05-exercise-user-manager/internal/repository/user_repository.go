package repository

import (
	"log"
	"user-management-api/internal/models"
)

type InMemoryUserRepository struct {
	users []models.User
}

func NewInMemoryUserRepository() UserRepository {
	return &InMemoryUserRepository{
		users: make([]models.User, 0),
	}
}

func (ur *InMemoryUserRepository) FindAll() {
	log.Println("GetAllUsers into UserRepository")
}

func (ur *InMemoryUserRepository) Create(user models.User) error {
	ur.users = append(ur.users, user)
	return nil
}

func (ur *InMemoryUserRepository) FindByUUID() {

}

func (ur *InMemoryUserRepository) Update() {

}

func (ur *InMemoryUserRepository) Delete() {

}

func (ur *InMemoryUserRepository) FindByEmail(email string) (models.User, bool) {
	for _, user := range ur.users {
		if user.Email == email {
			return user, true
		}
	}

	return models.User{}, false
}
