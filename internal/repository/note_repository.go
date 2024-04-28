package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"note_app/internal/models"
	"note_app/pkg/utils"
)

// NoteRepository интерфейс для работы с заметками в базе данных.
type NoteRepository interface {
	AddNote(ctx context.Context, note *models.Note) (int, error)
	GetNoteByID(ctx context.Context, noteID int) (*models.Note, error)
	UpdateNote(ctx context.Context, noteID int, note *models.Note) error
	DeleteNote(ctx context.Context, noteID int) error
	GetNotesByUserID(ctx context.Context, userID, offset, limit int) ([]models.Note, error)
	GetNotesByDateRange(ctx context.Context, startDate, endDate string, offset, limit int) ([]models.Note, error)
	GetNotesByDay(ctx context.Context, day string, offset, limit int) ([]models.Note, error)
	GetNotes(ctx context.Context, offset, limit int) ([]models.Note, error)
	GetNotesByUserIDAndDateRange(ctx context.Context, userID int, startDate, endDate string, offset, limit int) ([]models.Note, error)
	GetNotesByUserIDAndDate(ctx context.Context, userID int, date string, offset, limit int) ([]models.Note, error)
	GetNotesByUserIDAndDay(ctx context.Context, userID int, day string, offset, limit int) ([]models.Note, error)
}

// noteRepository реализация интерфейса NoteRepository.
type noteRepository struct {
	db *sql.DB
}

// NewNoteRepository создает новый экземпляр NoteRepository.
func NewNoteRepository(db *sql.DB) NoteRepository {
	return &noteRepository{db: db}
}

// AddNote добавляет новую заметку в базу данных.
func (nr *noteRepository) AddNote(ctx context.Context, note *models.Note) (int, error) {
	var id int
	query := `
		INSERT INTO notes (user_id, title, text, created_at, author)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := nr.db.QueryRowContext(
		ctx, query,
		note.UserID, note.Title, note.Text, note.CreatedAt, note.Author,
	).Scan(&id)
	if err != nil {
		log.Printf("Ошибка при добавлении заметки: %v", err)
		return 0, fmt.Errorf("не удалось добавить заметку: %v", err)
	}
	return id, nil
}

// GetNoteByID возвращает заметку по её ID из базы данных.
func (nr *noteRepository) GetNoteByID(ctx context.Context, noteID int) (*models.Note, error) {
	var note models.Note
	query := `
		SELECT id, user_id, title, text, created_at 
		FROM notes 
		WHERE id = $1
	`
	err := nr.db.QueryRowContext(ctx, query, noteID).
		Scan(&note.ID, &note.UserID, &note.Title, &note.Text, &note.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("заметка не найдена по ID: %d", noteID)
		}
		log.Printf("Ошибка при получении заметки по ID: %v", err)
		return nil, fmt.Errorf("не удалось получить заметку по ID: %v", err)
	}
	return &note, nil
}

// UpdateNote обновляет заметку в базе данных.
func (nr *noteRepository) UpdateNote(ctx context.Context, noteID int, note *models.Note) error {
	query := `
        UPDATE notes 
        SET title = $1, text = $2 
        WHERE id = $3
    `
	result, err := nr.db.ExecContext(ctx, query, note.Title, note.Text, noteID)
	if err != nil {
		log.Printf("Ошибка при обновлении заметки: %v", err)
		return fmt.Errorf("не удалось обновить заметку: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("нет затронутых строк, ID заметки: %d", noteID)
	}
	return nil
}

// DeleteNote удаляет заметку из базы данных.
func (nr *noteRepository) DeleteNote(ctx context.Context, noteID int) error {
	const deleteQuery = "DELETE FROM notes WHERE id = $1"

	result, err := nr.db.ExecContext(ctx, deleteQuery, noteID)
	if err != nil {
		log.Printf("Ошибка при удалении заметки: %v", err)
		return fmt.Errorf("не удалось удалить заметку: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("нет затронутых строк, ID заметки: %d", noteID)
	}

	return nil
}

// GetNotesByUserID возвращает заметки пользователя из базы данных.
func (nr *noteRepository) GetNotesByUserID(ctx context.Context, userID, offset, limit int) ([]models.Note, error) {
	query := `
		SELECT notes.id, notes.user_id, notes.title, notes.text, notes.created_at, users.username 
		FROM notes 
		INNER JOIN users ON notes.user_id = users.id 
		WHERE notes.user_id = $1 
		ORDER BY notes.created_at DESC 
		LIMIT $2 OFFSET $3
	`
	return utils.GetNotes(ctx, nr.db, query, userID, limit, offset)
}

// GetNotesByDateRange возвращает заметки за определенный период времени из базы данных.
func (nr *noteRepository) GetNotesByDateRange(ctx context.Context, startDate, endDate string, offset, limit int) ([]models.Note, error) {
	query := `
		SELECT notes.id, notes.user_id, notes.title, notes.text, notes.created_at, users.username 
		FROM notes 
		INNER JOIN users ON notes.user_id = users.id 
		WHERE notes.created_at >= $1 AND notes.created_at <= $2 
		ORDER BY notes.created_at DESC 
		LIMIT $3 OFFSET $4
	`
	return utils.GetNotes(ctx, nr.db, query, startDate, endDate, limit, offset)
}

// GetNotesByDay возвращает заметки за определенный день из базы данных.
func (nr *noteRepository) GetNotesByDay(ctx context.Context, day string, offset, limit int) ([]models.Note, error) {
	query := `
		SELECT notes.id, notes.user_id, notes.title, notes.text, notes.created_at, users.username 
		FROM notes 
		INNER JOIN users ON notes.user_id = users.id 
		WHERE DATE(notes.created_at) = $1 
		ORDER BY notes.created_at DESC 
		LIMIT $2 OFFSET $3
	`
	return utils.GetNotes(ctx, nr.db, query, day, limit, offset)
}

// GetNotes возвращает все заметки из базы данных.
func (nr *noteRepository) GetNotes(ctx context.Context, offset, limit int) ([]models.Note, error) {
	query := `
		SELECT notes.id, notes.user_id, notes.title, notes.text, notes.created_at, users.username 
		FROM notes 
		INNER JOIN users ON notes.user_id = users.id 
		ORDER BY notes.created_at DESC 
		LIMIT $1 OFFSET $2
	`
	return utils.GetNotes(ctx, nr.db, query, limit, offset)
}

// GetNotesByUserIDAndDateRange возвращает заметки пользователя в определенном диапазоне дат.
func (nr *noteRepository) GetNotesByUserIDAndDateRange(ctx context.Context, userID int, startDate, endDate string, offset, limit int) ([]models.Note, error) {
	query := `
        SELECT notes.id, notes.user_id, notes.title, notes.text, notes.created_at, users.username 
        FROM notes 
        INNER JOIN users ON notes.user_id = users.id 
        WHERE notes.user_id = $1 AND notes.created_at >= $2 AND notes.created_at <= $3
        ORDER BY notes.created_at DESC 
        LIMIT $4 OFFSET $5
    `
	return utils.GetNotes(ctx, nr.db, query, userID, startDate, endDate, limit, offset)
}

// GetNotesByUserIDAndDate возвращает заметки пользователя за определенную дату.
func (nr *noteRepository) GetNotesByUserIDAndDate(ctx context.Context, userID int, date string, offset, limit int) ([]models.Note, error) {
	query := `
        SELECT notes.id, notes.user_id, notes.title, notes.text, notes.created_at, users.username 
        FROM notes 
        INNER JOIN users ON notes.user_id = users.id 
        WHERE notes.user_id = $1 AND DATE(notes.created_at) = $2
        ORDER BY notes.created_at DESC 
        LIMIT $3 OFFSET $4
    `
	return utils.GetNotes(ctx, nr.db, query, userID, date, limit, offset)
}

// GetNotesByUserIDAndDay возвращает заметки, созданные указанным пользователем за указанный день
func (nr *noteRepository) GetNotesByUserIDAndDay(ctx context.Context, userID int, day string, offset, limit int) ([]models.Note, error) {
	query := `
		SELECT notes.id, notes.user_id, notes.title, notes.text, notes.created_at, users.username 
		FROM notes 
		INNER JOIN users ON notes.user_id = users.id 
		WHERE notes.user_id = $1 AND DATE(notes.created_at) = $2 
		ORDER BY notes.created_at DESC 
		LIMIT $3 OFFSET $4
	`
	return utils.GetNotes(ctx, nr.db, query, userID, day, limit, offset)
}
