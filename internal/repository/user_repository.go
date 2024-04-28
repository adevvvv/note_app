package repository

import (
	"database/sql"
	"fmt"
	"note_app/internal/models"
)

// UserRepository представляет интерфейс для работы с пользователями.
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(userID int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
}

// UserRepositoryImpl представляет реализацию интерфейса UserRepository.
type UserRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository создает новый экземпляр UserRepositoryImpl.
func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

// CreateUser создает нового пользователя.
func (ur *UserRepositoryImpl) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (username, password)
		VALUES ($1, $2)
	`
	_, err := ur.db.Exec(query, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("ошибка при создании пользователя: %v", err)
	}
	return nil
}

// GetUserByID возвращает пользователя по его ID.
func (ur *UserRepositoryImpl) GetUserByID(userID int) (*models.User, error) {
	query := `
		SELECT id, username, password
		FROM users
		WHERE id = $1
	`
	var user models.User
	err := ur.db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователя по ID: %v", err)
	}
	return &user, nil
}

// GetUserByUsername возвращает пользователя по его имени пользователя.
func (ur *UserRepositoryImpl) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, password
		FROM users
		WHERE username = $1
	`
	var user models.User
	err := ur.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователя по имени пользователя: %v", err)
	}
	return &user, nil
}
