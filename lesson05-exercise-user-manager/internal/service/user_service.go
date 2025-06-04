package service

import (
	"strings"
	"user-management-api/internal/models"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAllUsers(search string, page, limit int) ([]models.User, error) {
	users, err := us.repo.FindAll()
	if err != nil {
		return nil, utils.WrapError(err, "failed to fetch users", utils.ErrCodeInternal)
	}

	var filteredUsers []models.User
	if search != "" {
		search = strings.ToLower(search)
		for _, user := range users {
			name := strings.ToLower(user.Name)
			email := strings.ToLower(user.Email)

			if strings.Contains(name, search) || strings.Contains(email, search) {
				filteredUsers = append(filteredUsers, user)
			}
		}
	} else {
		filteredUsers = users
	}

	start := (page - 1) * limit
	if  start >= len(filteredUsers) {
		return []models.User{}, nil
	}

	end := start + limit
	if end > len(filteredUsers) {
		end = len(filteredUsers)
	}

	return filteredUsers[start:end], nil
}

func (us *userService) CreateUser(user models.User) (models.User, error) {
	user.Email = utils.NormalizeString(user.Email)

	if _, exist := us.repo.FindByEmail(user.Email); exist {
		return models.User{}, utils.NewError("email already exist", utils.ErrCodeConflict)
	}

	user.UUID = uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, utils.WrapError(err, "faild to hash password", utils.ErrCodeInternal)
	}

	user.Password = string(hashedPassword)

	if err := us.repo.Create(user); err != nil {
		return models.User{}, utils.WrapError(err, "faild to create user", utils.ErrCodeInternal)
	}

	return user, nil
}

func (us *userService) GetUserByUUID(uuid string) (models.User, error) {
	user, found := us.repo.FindByUUID(uuid)
	if !found {
		return models.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
	}

	return user, nil
}

func (us *userService) UpdateUser(uuid string, updatedUser models.User) (models.User, error) {
	updatedUser.Email = utils.NormalizeString(updatedUser.Email)

	if _, exist := us.repo.FindByEmail(updatedUser.Email); exist {
		return models.User{}, utils.NewError("email already exist", utils.ErrCodeConflict)
	}

	currentUser, found := us.repo.FindByUUID(uuid)
	if !found {
		return models.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
	}

	currentUser.Name = updatedUser.Name
	currentUser.Email = updatedUser.Email
	currentUser.Age = updatedUser.Age
	currentUser.Status = updatedUser.Status
	currentUser.Level = updatedUser.Level

	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.User{}, utils.WrapError(err, "faild to hash password", utils.ErrCodeInternal)
		}

		currentUser.Password = string(hashedPassword)
	}

	if err := us.repo.Update(uuid, currentUser); err != nil {
		return models.User{}, utils.WrapError(err, "faild to update user", utils.ErrCodeInternal)
	}

	return currentUser, nil
}

func (us *userService) DeleteUser() {

}
