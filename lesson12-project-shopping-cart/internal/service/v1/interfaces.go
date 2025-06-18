package v1service

type UserService interface {
	GetAllUsers(search string, page, limit int)
	CreateUser()
	GetUserByUUID(uuid string)
	UpdateUser(uuid string)
	DeleteUser(uuid string)
}
