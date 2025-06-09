package repository

import "log"

type SQLUserRepository struct {

}

func NewSQLUserRepository() UserRepository {
	return &SQLUserRepository{

	}
}

func (ur *SQLUserRepository) Create() {
	log.Println("Create")
}

func (ur *SQLUserRepository) FindById() {
	log.Println("Find By Id")
}