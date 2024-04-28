package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword хэширует пароль пользователя.
func HashPassword(password string) (string, error) {
	// Генерируем хэш пароля с использованием bcrypt и стандартной стоимости
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
