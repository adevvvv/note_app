package models

import "time"

type Note struct {
	ID                   int       `json:"id"`
	UserID               int       `json:"user_id"`
	Title                string    `json:"title"`
	Text                 string    `json:"text"`
	CreatedAt            time.Time `json:"created_at"`
	Author               string    `json:"author"`
	BelongsToCurrentUser bool      `json:"belongs_to_current_user,omitempty"`
}
