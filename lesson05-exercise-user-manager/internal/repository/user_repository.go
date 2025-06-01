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

func (ur *InMemoryUserRepository) Create() {

}

func (ur *InMemoryUserRepository) FindByUUID() {

}

func (ur *InMemoryUserRepository) Update() {

}

func (ur *InMemoryUserRepository) Delete() {

}
