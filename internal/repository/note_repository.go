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
	err := nr.db.QueryRowContext(
		ctx,
		"INSERT INTO notes (user_id, title, text, created_at, author) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		note.UserID,
		note.Title,
		note.Text,
		note.CreatedAt,
		note.Author,
	).Scan(&id)
	if err != nil {
		log.Printf("Ошибка при добавлении заметки: %v", err)
		return 0, fmt.Errorf("не удалось добавить заметку: %v", err)
	}
	return id, nil
}
