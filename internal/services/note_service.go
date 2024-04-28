package services

import (
	"context"
	"note_app/internal/models"
	"note_app/internal/repository"
)

// NoteService предоставляет методы для работы с заметками.
type NoteService interface {
	AddNote(ctx context.Context, note *models.Note) (int, error)
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
