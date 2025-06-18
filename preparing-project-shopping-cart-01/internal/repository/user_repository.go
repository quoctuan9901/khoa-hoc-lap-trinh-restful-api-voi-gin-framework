package repository

type SqlUserRepository struct {
}

func NewSqlUserRepository() UserRepository {
	return &SqlUserRepository{}
}

func (ur *SqlUserRepository) FindAll() {}

func (ur *SqlUserRepository) Create() {}

func (ur *SqlUserRepository) FindByUUID(uuid string) {}

func (ur *SqlUserRepository) Update(uuid string) {}

func (ur *SqlUserRepository) Delete(uuid string) {}

func (ur *SqlUserRepository) FindByEmail(email string) {}
