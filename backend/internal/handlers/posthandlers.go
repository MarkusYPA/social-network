package handlers

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func HandleGetFeedPosts(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")
	lastPostIdStr := r.URL.Query().Get("last_post_id")

	posts, err := service.GetFeedPosts(cursor, limitStr, lastPostIdStr, userId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
	}

	var imagePath *string
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		savedPath, saveErr := service.SaveUploadedFile(file, header, "posts")
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

func HandleCommentsForPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postIDstring := r.URL.Query().Get("post_id")
	if postIDstring == "" {
		fmt.Println("empty postId")
		json.NewEncoder(w).Encode([]model.User{}) // Return empty array for empty query
		return
	}

	postType := r.URL.Query().Get("type")
	if postType == "" {
		fmt.Println("empty postId")
		json.NewEncoder(w).Encode([]model.User{}) // Return empty array for empty query
		return
	}

	comments, err := service.CommentsForPost(postIDstring, postType, userID)
	if err != nil {
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}

func HandleCreateCommentsForPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postIDstring := r.URL.Query().Get("post_id")
	if postIDstring == "" {
		fmt.Println("empty postId")
		json.NewEncoder(w).Encode([]model.User{}) // Return empty array for empty query
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	var payload struct {
		Content string `json:"content"`
		Type    string `json:"type"`
	}
	// Get form values
	payload.Content = r.FormValue("content")
	payload.Type = r.FormValue("type")

	var imagePath *string
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		savedPath, saveErr := service.SaveUploadedFile(file, header, "comments")
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

	err = service.CreateCommentsForPost(userID, postIDstring, payload.Type, payload.Content, imagePath)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}
