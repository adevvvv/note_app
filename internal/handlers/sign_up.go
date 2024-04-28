package handlers

import (
	"net/http"
	"note_app/internal/models"
	"note_app/internal/services"
	"note_app/pkg/utils"

	"github.com/gin-gonic/gin"
)

// UserHandler обрабатывает запросы пользователя.
type UserHandler struct {
	UserService *services.UserService
}

// NewSignupHandler создает новый экземпляр UserHandler для обработки запросов на регистрацию пользователя.
func NewSignupHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (userHandler *UserHandler) SignUp(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	if err := utils.ValidateUser(&user); err != nil {
		c.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хэшировании пароля"})
		return
	}
	user.Password = hashedPassword

	if err := userHandler.UserService.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно зарегистрирован"})
}
