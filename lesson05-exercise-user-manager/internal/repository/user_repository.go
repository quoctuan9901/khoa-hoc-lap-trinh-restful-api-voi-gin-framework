package repository

import (
	"fmt"
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

func (ur *InMemoryUserRepository) FindAll() ([]models.User, error) {
	return ur.users, nil
}

func (ur *InMemoryUserRepository) Create(user models.User) error {
	ur.users = append(ur.users, user)
	return nil
}

func (ur *InMemoryUserRepository) FindByUUID(uuid string) (models.User, bool) {
	for _, user := range ur.users {
		if user.UUID == uuid {
			return user, true
		}
	}

	return models.User{}, false
}

func (ur *InMemoryUserRepository) Update(uuid string, user models.User) error {
	for i, u := range ur.users {
		if u.UUID == uuid {
			ur.users[i] = user
			return nil
		}
	}

	return fmt.Errorf("user not found")
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
