package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ValidateSession extracts the session cookie, validates it, and returns the user_id
func ValidateSession(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, errors.New("missing session cookie")
	}

	userID, expiresAt, err := repository.GetUserIdAndExpirationBySessionId(cookie)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("invalid session")
		}
		return 0, err
	}

	if time.Now().After(expiresAt) {
		return 0, errors.New("session expired")
	}

	return userID, nil
}

func Login(w http.ResponseWriter, r *http.Request) (model.User, int) {
	var user model.User
	var emptyUser model.User

	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return user, http.StatusBadRequest
	}

	user, _, err := repository.GetUserByEmail(req)
	if err != nil {
		return emptyUser, http.StatusUnauthorized
	}

	// Compare the entered password with the stored hashed password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// Password is incorrect
		return emptyUser, http.StatusUnauthorized
	}

	err = CreateSession(user, w)
	if err != nil {
		return user, http.StatusInternalServerError
	}

	user.Password = ""
	return user, http.StatusOK
}

func CreateSession(user model.User, w http.ResponseWriter) error {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	err := repository.InsertSession(sessionID, user, expiresAt)
	if err != nil {
		fmt.Printf("ERROR: Failed to insert session for user : %v\n", err)
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

func RemoveSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		err = repository.DeleteSessionById(cookie.Value) // proceed to wipe cookie regardless
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // expired
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return err
}
