package repository

import (
	"hoc-gin/internal/models"
	"log"
)

type SQLUserRepository struct {
}

func NewSQLUserRepository() UserRepository {
	return &SQLUserRepository{}
}

func (ur *SQLUserRepository) Create(user *models.User) {
	log.Println("Create")
}

func (ur *SQLUserRepository) FindById(id int) {
	log.Println("Find By Id")
}
