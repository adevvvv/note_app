package utils

import (
	"github.com/golang-jwt/jwt"
	"net/http"
	"note_app/internal/models"
	"time"
)

// GenerateToken создает и подписывает токен JWT для пользователя.
func GenerateToken(userID int, JWTKey []byte) (string, error) {
	// Устанавливаем время истечения токена на 24 часа от текущего времени
	expirationTime := time.Now().Add(24 * time.Hour)
	// Создаем кастомные claims для JWT, включая идентификатор пользователя и время истечения
	claims := &models.Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Создаем новый токен с указанными claims и методом подписи
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Подписываем токен с использованием ключа
	return token.SignedString(JWTKey)
}

// SetTokenCookie устанавливает токен в виде куки в ответе.
func SetTokenCookie(w http.ResponseWriter, tokenString string) {
	// Устанавливаем время истечения куки на 24 часа от текущего времени
	expirationTime := time.Now().Add(24 * time.Hour)
	// Создаем новую куку с именем "token", значением токена и временем истечения
	cookie := http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	}
	// Устанавливаем куку в ответ
	http.SetCookie(w, &cookie)
}
