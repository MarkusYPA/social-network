package handlers

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	user, statusCode := service.Login(w, r)

	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // Error code
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"user": user,
	})
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	err := service.RemoveSession(w, r)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// "/api/me" is a common RESTful convention that returns the currently authenticated user's info
func HandleMe(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	user, err := repository.GetUserById(userID, true)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed"})
		fmt.Println("00")
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to parse form"})
		return
	}

	// Required fields
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	firstName := strings.TrimSpace(r.FormValue("firstName"))
	lastName := strings.TrimSpace(r.FormValue("lastName"))
	dob := r.FormValue("dob")

	// Optional fields
	nickname := strings.TrimSpace(r.FormValue("nickname"))
	about := strings.TrimSpace(r.FormValue("about"))

	if email == "" || password == "" || firstName == "" || lastName == "" || dob == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Missing required fields"})
		return
	}

	// Handle optional avatar upload
	var avatarPath sql.NullString
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		avatarPath, err = service.UploadAvatar(file, header)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Failed to save image"})
			return
		}
	} else {
		fmt.Println("Image reading error at registering:", err)
	}

	userInfo := struct {
		Email      string
		Password   string
		FirstName  string
		LastName   string
		DOB        string
		Nickname   string
		About      string
		AvatarPath sql.NullString
	}{
		Email:      email,
		Password:   password,
		FirstName:  firstName,
		LastName:   lastName,
		DOB:        dob,
		Nickname:   nickname,
		About:      about,
		AvatarPath: avatarPath,
	}

	errMsg, statusCode := service.RegisterUser(userInfo)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"message": errMsg})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func HandleUpdateMe(w http.ResponseWriter, r *http.Request) {

	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error parsing form"})
		return
	}

	updateData := model.UpdateProfileData{
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
		DOB:       r.FormValue("dob"),
		Nickname:  r.FormValue("nickname"),
		About:     r.FormValue("about"),
	}

	if r.FormValue("is_public") == "true" {
		updateData.IsPublic = true
	}

	// Handle optional avatar
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		updateData.AvatarPath, err = service.UploadAvatar(file, header)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save image"})
			return
		}
	} else {
		//fmt.Println("No avatar file found at updating profile:", err)
	}

	// Handle delete_avatar flag
	if r.FormValue("delete_avatar") == "true" {
		updateData.DeleteAvatar = true
	}

	usr, errMsg, errStatus := service.UpdateUserProfile(userId, updateData)
	if errMsg != "" {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errStatus)
		json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
		fmt.Println("Error 3", errMsg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usr)
}
