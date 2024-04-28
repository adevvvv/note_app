package utils

import (
	"context"
	"database/sql"
	"note_app/internal/models"
)

// ScanNotes сканирует ряды и возвращает слайс заметок.
func ScanNotes(rows *sql.Rows) ([]models.Note, error) {
	var notes []models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Text, &note.CreatedAt, &note.Author)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}

// GetNotes возвращает заметки из базы данных с учетом заданных параметров.
func GetNotes(ctx context.Context, db *sql.DB, query string, args ...interface{}) ([]models.Note, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return ScanNotes(rows)
}
