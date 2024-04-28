package services

import (
	"context"
	"note_app/internal/models"
	"note_app/internal/repository"
)

// NoteService предоставляет методы для работы с заметками.
type NoteService interface {
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

// noteService реализация интерфейса NoteService.
type noteService struct {
	repo repository.NoteRepository
}

// NewNoteService создает новый экземпляр NoteService.
func NewNoteService(repo repository.NoteRepository) NoteService {
	return &noteService{repo: repo}
}

// AddNote добавляет новую заметку.
func (ns *noteService) AddNote(ctx context.Context, note *models.Note) (int, error) {
	return ns.repo.AddNote(ctx, note)
}

// GetNoteByID возвращает заметку по её ID.
func (ns *noteService) GetNoteByID(ctx context.Context, noteID int) (*models.Note, error) {
	return ns.repo.GetNoteByID(ctx, noteID)
}

// UpdateNote обновляет заметку.
func (ns *noteService) UpdateNote(ctx context.Context, noteID int, note *models.Note) error {
	return ns.repo.UpdateNote(ctx, noteID, note)
}

// DeleteNote удаляет заметку.
func (ns *noteService) DeleteNote(ctx context.Context, noteID int) error {
	return ns.repo.DeleteNote(ctx, noteID)
}

// GetNotesByUserID возвращает заметки пользователя.
func (ns *noteService) GetNotesByUserID(ctx context.Context, userID, offset, limit int) ([]models.Note, error) {
	return ns.repo.GetNotesByUserID(ctx, userID, offset, limit)
}

// GetNotesByDateRange возвращает заметки за определенный период времени.
func (ns *noteService) GetNotesByDateRange(ctx context.Context, startDate, endDate string, offset, limit int) ([]models.Note, error) {
	return ns.repo.GetNotesByDateRange(ctx, startDate, endDate, offset, limit)
}

// GetNotesByDay возвращает заметки за определенный день.
func (ns *noteService) GetNotesByDay(ctx context.Context, day string, offset, limit int) ([]models.Note, error) {
	return ns.repo.GetNotesByDay(ctx, day, offset, limit)
}

// GetNotes возвращает все заметки.
func (ns *noteService) GetNotes(ctx context.Context, offset, limit int) ([]models.Note, error) {
	return ns.repo.GetNotes(ctx, offset, limit)
}

// GetNotesByUserIDAndDateRange возвращает заметки пользователя в определенном диапазоне дат.
func (ns *noteService) GetNotesByUserIDAndDateRange(ctx context.Context, userID int, startDate, endDate string, offset, limit int) ([]models.Note, error) {
	return ns.repo.GetNotesByUserIDAndDateRange(ctx, userID, startDate, endDate, offset, limit)
}

// GetNotesByUserIDAndDate возвращает заметки пользователя за определенную дату.
func (ns *noteService) GetNotesByUserIDAndDate(ctx context.Context, userID int, date string, offset, limit int) ([]models.Note, error) {
	return ns.repo.GetNotesByUserIDAndDate(ctx, userID, date, offset, limit)
}

// GetNotesByUserIDAndDay возвращает заметки, созданные указанным пользователем за указанный день.
func (ns *noteService) GetNotesByUserIDAndDay(ctx context.Context, userID int, day string, offset, limit int) ([]models.Note, error) {
	return ns.repo.GetNotesByUserIDAndDay(ctx, userID, day, offset, limit)
}
