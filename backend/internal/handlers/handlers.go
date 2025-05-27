package handlers

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	user, statusCode := service.Login(w, r)

	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // Error code
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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
	json.NewEncoder(w).Encode(user)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		fmt.Println("00")
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		fmt.Println("01", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
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
		fmt.Println("02 missing fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Handle optional avatar upload
	var avatarPath sql.NullString
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		avatarPath, err = service.UploadAvatar(file, header)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
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
		http.Error(w, errMsg, statusCode)
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
		http.Error(w, "Error parsing form", http.StatusBadRequest)
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
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
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
		http.Error(w, errMsg, errStatus)
		fmt.Println("Error 3", errMsg)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usr)
}

func HandleUserByID(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	targetId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	usr, statusCode := service.UserById(userId, targetId)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // Error code
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usr)
}

// SearchUsers handles the user search request.
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		json.NewEncoder(w).Encode([]model.User{}) // Return empty array for empty query
		return
	}

	users, err := repository.SearchUsers(query)
	if err != nil {
		http.Error(w, "Error searching users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func GetFeedPosts(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")
	lastPostIdStr := r.URL.Query().Get("last_post_id")
	limit := 10
	lastPostId := 0
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if lpid, err := strconv.Atoi(lastPostIdStr); err == nil && lpid > 0 {
		lastPostId = lpid
	}

	cursorTime := time.Now().UTC() // use current time by default

	if cursor != "" {
		t, err := time.Parse(time.RFC3339, cursor)
		if err == nil {
			cursorTime = t
		}
	}

	cursorTime = cursorTime.Truncate(time.Second)

	posts, err := repository.GetFeedPostsBefore(userId, cursorTime, limit, lastPostId)
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func HandlePostsByUserId(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	targetId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var posts []model.Post
	posts, err = repository.GetPostsByUserId(userId, targetId)
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// handleCreatePost adds a post to the database and returns the new one
func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		fmt.Println("error reading data at HandleCreatePost", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	privacyLvl := r.FormValue("privacy_level")
	if content == "" || privacyLvl == "" {
		fmt.Println("Missing fields at HandleCreatePost", err)
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	var viewers []int
	if privacyLvl == "private" {
		err := json.Unmarshal([]byte(r.FormValue("selected_viewers")), &viewers)
		if err != nil {
			http.Error(w, "Invalid viewers list", http.StatusBadRequest)
			return
		}
		//fmt.Println("viewers:", viewers)
	}

	var imagePath *string
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		savedPath, saveErr := service.SaveUploadedFile(file, header)
		if saveErr != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
		imagePath = &savedPath
	} else if err != http.ErrMissingFile {
		fmt.Println("Error reading file at CreateGroupPostHandler", err)
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}

	post, statusCode := service.CreatePost(content, privacyLvl, imagePath, userID, viewers)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func GetSuggestedUsers(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	suggestions, err := repository.GetSuggestedUsers(userId)
	if err != nil {
		http.Error(w, "Failed to fetch suggestions", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(suggestions)
}
