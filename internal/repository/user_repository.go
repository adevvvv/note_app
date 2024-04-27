package repository

import (
	"database/sql"
	"fmt"
	"note_app/internal/models"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
}

// UserRepositoryImpl реализация интерфейса UserRepository
type UserRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository создает новый экземпляр UserRepositoryImpl
func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

// CreateUser создает нового пользователя
func (ur *UserRepositoryImpl) CreateUser(user *models.User) error {
	_, err := ur.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("ошибка при создании пользователя: %v", err)
	}
	return nil
}

// GetUserByUsername возвращает пользователя по его имени пользователя
func (ur *UserRepositoryImpl) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := ur.db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка при вводе пользователя по имени пользователя: %v", err)
	}
	return &user, nil
}
