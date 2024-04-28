package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"note_app/internal/models"
	"note_app/internal/services"
	"note_app/pkg/utils"
	"time"
)

// NoteHandler обрабатывает запросы, связанные с заметками.
type NoteHandler struct {
	NoteService services.NoteService
	UserService *services.UserService
	JWTSecret   string
}

// NewNoteHandler создает новый экземпляр NoteHandler для обработки запросов, связанных с заметками.
func NewNoteHandler(noteService services.NoteService, userService *services.UserService, jwtSecret string) *NoteHandler {
	return &NoteHandler{
		NoteService: noteService,
		UserService: userService,
		JWTSecret:   jwtSecret,
	}
}

// AddNote обрабатывает запрос на добавление новой заметки.
func (noteHandler *NoteHandler) AddNote(c *gin.Context) {
	var note models.Note
	if err := c.BindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	// Проверяем длину заголовка и текста.
	if !utils.CheckNoteLength(note.Title, note.Text) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Превышена максимальная длина заголовка или текста"})
		return
	}

	// Получаем токен из куки запроса
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
		return
	}

	// Проверяем и расшифровываем токен
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(noteHandler.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
		return
	}

	// Получаем утверждения (claims) из токена
	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
		return
	}

	user, err := noteHandler.UserService.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении информации о пользователе"})
		return
	}

	note.UserID = claims.UserID
	note.CreatedAt = time.Now()
	note.Author = user.Username

	id, err := noteHandler.NoteService.AddNote(context.Background(), &note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении заметки"})
		return
	}

	note.ID = id

	c.JSON(http.StatusOK, note)
}
