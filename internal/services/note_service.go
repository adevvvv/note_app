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
