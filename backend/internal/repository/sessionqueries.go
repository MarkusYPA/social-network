package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"time"
)

func InsertSession(sessionID string, user model.User, expiresAt time.Time) error {
	_, err := database.DB.Exec("INSERT INTO sessions (session_token, user_id, expires_at) VALUES (?, ?, ?)", sessionID, user.ID, expiresAt)
	return err
}

func DeleteSessionById(id string) error {
	_, err := database.DB.Exec("DELETE FROM sessions WHERE session_token = ?", id)
	return err
}
