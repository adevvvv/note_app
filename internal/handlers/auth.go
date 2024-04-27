package handlers

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"note_app/internal/models"
	"note_app/internal/repository"
	"regexp"
	"strconv"
	"time"
)

// SignupHandler обрабатывает запрос на регистрацию нового пользователя.
func SignupHandler(userRepository repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		// Проверка длины логина и пароля
		if err := validateUsername(user.Username); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := validatePassword(user.Password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Хеширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Ошибка хеширования пароля", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)

		// Создание пользователя в базе данных
		err = userRepository.CreateUser(&user)
		if err != nil {
			http.Error(w, "Пользователь с таким именем уже существует", http.StatusInternalServerError)
			return
		}

		// Отправка ответа об успешной регистрации
		json.NewEncoder(w).Encode("Пользователь успешно зарегистрирован")
	}
}

// LoginHandler обрабатывает запрос на аутентификацию пользователя.
func LoginHandler(userRepository repository.UserRepository, jwtKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		dbUser, err := userRepository.GetUserByUsername(user.Username)
		if err != nil {
			http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
			http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
			return
		}

		// Преобразуем тип int в uint64
		userID := uint64(dbUser.ID)

		token, err := generateToken(userID, jwtKey) // Используем userID
		if err != nil {
			http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
			return
		}

		response := map[string]string{"token": token}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
func validateUsername(username string) error {
	if len(username) < 4 || len(username) > 20 {
		return errors.New("имя пользователя должно содержать от 4 до 20 символов")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return errors.New("имя пользователя может содержать только буквы (латинские), цифры и символ подчеркивания")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 6 || len(password) > 20 {
		return errors.New("пароль должен содержать от 6 до 20 символов")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_!?@#$%^&*()-+=]+$`).MatchString(password) {
		return errors.New("пароль может содержать только буквы (латинские), цифры и специальные символы: ^[a-zA-Z0-9_!?@#$%^&*()-+=]+$")
	}
	return nil
}

func generateToken(userID uint64, jwtKey string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"userID": strconv.FormatUint(userID, 10),
		"exp":    expirationTime.Unix(),
		"iss":    "vkNotesSolid",
		"iat":    time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}
