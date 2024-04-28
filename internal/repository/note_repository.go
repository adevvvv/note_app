package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"note_app/internal/models"
)

// NoteRepository интерфейс для работы с заметками в базе данных.
type NoteRepository interface {
	AddNote(ctx context.Context, note *models.Note) (int, error)
	GetNoteByID(ctx context.Context, noteID int) (*models.Note, error)
	UpdateNote(ctx context.Context, noteID int, note *models.Note) error
	DeleteNote(ctx context.Context, noteID int) error
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
