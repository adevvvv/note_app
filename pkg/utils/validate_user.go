package utils

import (
	"net/http"
	"note_app/internal/models"
	"regexp"
)

// ValidateUser проверяет валидность данных пользователя.
func ValidateUser(user *models.User) *HTTPError {
	const (
		minUsernameLength = 4
		maxUsernameLength = 20
		minPasswordLength = 6
		maxPasswordLength = 20
	)

	// Проверка длины имени пользователя
	usernameLength := len(user.Username)
	if usernameLength < minUsernameLength || usernameLength > maxUsernameLength {
		return &HTTPError{
			Message: "Имя пользователя должно быть от 4 до 20 символов",
			Code:    http.StatusBadRequest,
		}
	}

	// Проверка на допустимые символы в имени пользователя
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(user.Username) {
		return &HTTPError{
			Message: "Имя пользователя может содержать только буквы (латинские), цифры и символ подчеркивания",
			Code:    http.StatusBadRequest,
		}
	}

	// Проверка длины пароля
	passwordLength := len(user.Password)
	if passwordLength < minPasswordLength || passwordLength > maxPasswordLength {
		return &HTTPError{
			Message: "Пароль должен быть от 6 до 20 символов",
			Code:    http.StatusBadRequest,
		}
	}

	// Проверка на допустимые символы в пароле
	if !regexp.MustCompile(`^[a-zA-Z0-9_!?@#$%^&*()-+=]+$`).MatchString(user.Password) {
		return &HTTPError{
			Message: "Пароль может содержать только буквы (латинские), цифры и следующие специальные символы: !?@#$%^&*()-+=",
			Code:    http.StatusBadRequest,
		}
	}

	return nil
}
