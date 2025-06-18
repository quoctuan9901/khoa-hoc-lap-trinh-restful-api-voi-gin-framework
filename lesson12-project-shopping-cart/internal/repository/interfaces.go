package repository

type UserRepository interface {
	FindAll()
	Create()
	FindByUUID(uuid string)
	Update(uuid string)
	Delete(uuid string)
	FindByEmail(email string)
}
