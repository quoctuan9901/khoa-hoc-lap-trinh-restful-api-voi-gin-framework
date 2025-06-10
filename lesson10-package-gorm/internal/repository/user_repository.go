package repository

import (
	"fmt"
	"hoc-gin/internal/models"

	"gorm.io/gorm"
)

type SQLUserRepository struct {
	db *gorm.DB
}

func NewSQLUserRepository(db *gorm.DB) UserRepository {
	return &SQLUserRepository{
		db: db,
	}
}

func (ur *SQLUserRepository) Create(user *models.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (ur *SQLUserRepository) FindById(user *models.User, id int) error {
	if err := ur.db.First(user, id).Error; err != nil {
		return err
	}

	return nil
}
