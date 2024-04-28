package handlers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"note_app/internal/models"
	"note_app/internal/services"
	"note_app/pkg/utils"
)

// LoginHandler обрабатывает запросы на аутентификацию пользователя.
type LoginHandler struct {
	UserService *services.UserService
	JWTKey      []byte
}

// NewSignInHandler создает новый экземпляр LoginHandler для обработки запросов на аутентификацию.
func NewSignInHandler(userService *services.UserService, jwtKey string) *LoginHandler {
	return &LoginHandler{
		UserService: userService,
		JWTKey:      []byte(jwtKey),
	}
}

func (loginHandler *LoginHandler) SignIn(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное тело запроса"})
		return
	}

	dbUser, err := loginHandler.UserService.GetUserByUsername(user.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
		return
	}

	tokenString, err := utils.GenerateToken(dbUser.ID, loginHandler.JWTKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	utils.SetTokenCookie(c.Writer, tokenString)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
