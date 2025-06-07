package handlers

import (
	"backend/internal/repository"
	"backend/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func HandleGetUserMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chats, err := repository.GetUserChats(userID)
	if err != nil {
		//fmt.Println("error getting chats at HandleGetUserMessages", err)
		http.Error(w, "Failed to fetch chats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}

func HandleGetGroupMessagesByGroupId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/group/chat/messages/")
	groupId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//fmt.Println("The user and group Ids:", userID, groupId)

	chat, err := repository.GetGroupChat(groupId)
	if err != nil {
		//fmt.Println("error getting chats at HandleGetGroupMessagesByGroupId", err)
		http.Error(w, "Failed to fetch chat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)

	//w.WriteHeader(http.StatusOK)
}
