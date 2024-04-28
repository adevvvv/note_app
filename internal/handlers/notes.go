package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"note_app/internal/models"
	"note_app/internal/services"
	"note_app/pkg/utils"
	"strconv"
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

// EditNote обрабатывает запрос на редактирование существующей заметки.
func EditNote(ns services.NoteService, us *services.UserService, jwtKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлечение токена из куки запроса
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			return
		}

		// Проверка токена
		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			return
		}

		claims, ok := token.Claims.(*models.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			return
		}
		userID := claims.UserID

		// Получение идентификатора заметки из параметров URL
		noteIDStr := c.Param("id")
		noteID, err := strconv.Atoi(noteIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор заметки"})
			return
		}

		// Получение заметки по идентификатору
		note, err := ns.GetNoteByID(context.Background(), noteID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Заметка не найдена"})
			return
		}

		// Проверка, является ли пользователь владельцем заметки
		if userID != note.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Не имеющий права редактировать эту заметку"})
			return
		}

		// Проверка, можно ли редактировать заметку (в течение 24 часов с момента создания).
		if time.Since(note.CreatedAt) > 24*time.Hour {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Заметка не может быть отредактирована по истечении 1 дня"})
			return
		}

		var updatedNote models.Note
		if err := c.BindJSON(&updatedNote); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимый формат запроса"})
			return
		}

		// Проверяем длину заголовка и текста.
		if !utils.CheckNoteLength(updatedNote.Title, updatedNote.Text) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Превышена максимальная длина заголовка или текста"})
			return
		}

		// Получение информации об авторе заметки
		author, err := us.GetUserByID(note.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении информации об авторе"})
			return
		}

		updatedNote.ID = noteID
		updatedNote.UserID = userID
		updatedNote.CreatedAt = note.CreatedAt
		updatedNote.Author = author.Username

		if err := ns.UpdateNote(context.Background(), noteID, &updatedNote); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Примечание об ошибке при обновлении"})
			return
		}

		c.JSON(http.StatusOK, updatedNote)
	}
}

// DeleteNote обрабатывает запрос на удаление заметки.
func DeleteNote(ns services.NoteService, jwtKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлечение токена из куки запроса
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			return
		}

		// Проверка токена
		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			return
		}

		claims, ok := token.Claims.(*models.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			return
		}
		userID := claims.UserID

		// Получение идентификатора заметки из параметров URL
		noteIDStr := c.Param("id")
		noteID, err := strconv.Atoi(noteIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор заметки"})
			return
		}

		// Получение заметки по идентификатору
		note, err := ns.GetNoteByID(context.Background(), noteID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Заметка не найдена"})
			return
		}

		// Проверка, является ли пользователь владельцем заметки
		if userID != note.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Несанкционированное удаление этой заметки"})
			return
		}

		// Удаление заметки
		if err := ns.DeleteNote(context.Background(), noteID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении заметки"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Заметка успешно удалена"})
	}
}
