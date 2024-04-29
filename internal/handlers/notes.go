package handlers

import (
	"context"
	"fmt"
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
// @Summary Добавление новой заметки
// @Description Обрабатывает запрос на добавление новой заметки.
// @Accept json
// @Produce json
// @Param body body models.NoteInput true "Данные новой заметки"
// @Router /notes [post]
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

// EditNoteHandler обрабатывает запрос на редактирование заметки.
// @Summary Редактирование заметки
// @Description Обрабатывает запрос на редактирование заметки.
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор заметки"
// @Param body body models.NoteInput true "Новые данные заметки"
// @Router /notes/{id} [put]
func EditNoteHandler(ns services.NoteService, us *services.UserService, jwtKey string) gin.HandlerFunc {
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

// DeleteNoteHandler обрабатывает запрос на удаление заметки.
// @Summary Удаление заметки
// @Description Обрабатывает запрос на удаление заметки.
// @Produce json
// @Param id path int true "Идентификатор заметки"
// @Router /notes/{id} [delete]
func DeleteNoteHandler(ns services.NoteService, jwtKey string) gin.HandlerFunc {
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

// GetNotesHandler обрабатывает запрос на получение заметок с возможностью фильтрации.
// @Summary Получение заметок
// @Description Обрабатывает запрос на получение заметок с возможностью фильтрации.
// @Accept json
// @Produce json
// @Param start_date query string false "Дата начала в формате 'ГГГГ-ММ-ДД'"
// @Param end_date query string false "Дата окончания в формате 'ГГГГ-ММ-ДД'"
// @Param username query string false "Имя пользователя"
// @Param date query string false "Дата в формате 'ГГГГ-ММ-ДД'"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество записей на странице"
// @Router /notes [get]
func GetNotesHandler(ns services.NoteService, us services.UserService, jwtKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлечение токена из куки запроса
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			return
		}

		// Проверка и расшифровка токена
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
		currentUserID := claims.UserID

		// Извлечение параметров фильтрации из URL-запроса
		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")
		username := c.Query("username")
		dateStr := c.Query("date")

		fmt.Println("Start Date:", startDateStr)
		fmt.Println("End Date:", endDateStr)
		fmt.Println("Username:", username)
		fmt.Println("Date:", dateStr)

		// Преобразование параметров фильтрации
		var startDate, endDate, date time.Time
		if startDateStr != "" {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты начала. Используйте 'ГГГГ-ММ-ДД'"})
				return
			}
		}
		if endDateStr != "" {
			endDate, err = time.Parse("2006-01-02", endDateStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты окончания. Используйте 'ГГГГ-ММ-ДД'"})
				return
			}
		}
		if dateStr != "" {
			date, err = time.Parse("2006-01-02", dateStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты. Используйте 'ГГГГ-ММ-ДД'"})
				return
			}
		}

		// Получение идентификатора пользователя по его имени, если указан параметр username
		var filterUserID int
		if username != "" {
			user, err := us.GetUserByUsername(username)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь не найден"})
				return
			}
			filterUserID = user.ID
		}

		// Извлечение параметров страницы из URL-запроса
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page <= 0 {
			page = 1
		}
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil || limit <= 0 {
			limit = 10
		}

		offset := (page - 1) * limit

		fmt.Println("Offset:", offset)
		fmt.Println("Limit:", limit)

		// Фильтрация заметок по дате создания, пользователю и дню
		var notes []models.Note
		var errorGetNotes error
		switch {
		case !startDate.IsZero() && !endDate.IsZero() && username != "":
			notes, errorGetNotes = ns.GetNotesByUserIDAndDateRange(context.Background(), filterUserID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), offset, limit)
		case !startDate.IsZero() && !endDate.IsZero():
			notes, errorGetNotes = ns.GetNotesByDateRange(context.Background(), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), offset, limit)
		case !date.IsZero() && username != "":
			notes, errorGetNotes = ns.GetNotesByUserIDAndDay(context.Background(), filterUserID, date.Format("2006-01-02"), offset, limit)
		case !date.IsZero():
			notes, errorGetNotes = ns.GetNotesByDay(context.Background(), date.Format("2006-01-02"), offset, limit)
		case username != "":
			notes, errorGetNotes = ns.GetNotesByUserID(context.Background(), filterUserID, offset, limit)
		default:
			notes, errorGetNotes = ns.GetNotes(context.Background(), offset, limit)
		}

		if errorGetNotes != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении заметок"})
			return
		}

		fmt.Println("Number of Notes:", len(notes))

		// Создание списка для ответа
		response := make([]gin.H, 0, len(notes))
		for _, note := range notes {
			author, err := us.GetUserByID(note.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении информации об авторе заметки"})
				return
			}

			noteData := gin.H{
				"title":  note.Title,
				"text":   note.Text,
				"author": author.Username,
			}

			// Добавление признака belongsToCurrentUser только если он равен true
			if note.UserID == currentUserID {
				noteData["belongsToCurrentUser"] = true
			}

			response = append(response, noteData)
		}

		c.JSON(http.StatusOK, response)
	}
}
