package services

import "note_app/internal/models"

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	CreateUser(user *models.User) error
}

// UserService реализация интерфейса UserRepository
type UserService struct {
	userRepository UserRepository
}

// NewUserService создает новый экземпляр UserService
func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

// CreateUser создает нового пользователя
func (us *UserService) CreateUser(user *models.User) error {
	return us.userRepository.CreateUser(user)
}
