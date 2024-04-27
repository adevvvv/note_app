package repository

import (
	"database/sql"
	"fmt"
	"note_app/internal/models"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	CreateUser(user *models.User) error
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
