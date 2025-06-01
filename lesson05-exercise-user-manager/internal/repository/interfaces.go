package repository

type UserRepository interface {
	FindAll()
	Create()
	FindByUUID()
	Update()
	Delete()
}