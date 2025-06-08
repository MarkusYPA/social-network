package handlers

import (
	"backend/internal/repository"
	"backend/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// GetNotifications handles GET /api/notifications and returns all notifications for the authenticated user.
func HandleGetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("not authorized at GetNotifications:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notifications, statusCode := service.GetNotifications(userID)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) {
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(notifications); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetNewNotifications(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	count, err := repository.GetNewNotificatonsCount(userID)
	if err != nil {
		http.Error(w, "Failed to count new notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]int{"count": count}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

// GetNotificationByID handles GET /api/notifications/{id} and returns a single notification for the authenticated user.
func GetNotificationByID(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/notifications/")
	notificationID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	notification, err := repository.GetNotification(userID, notificationID)
	if err != nil {
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notification)
}

// ReadNotification handles POST /api/notifications/{id}/read to mark a notification as read.
func ReadNotification(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/notifications/")
	idStr = strings.TrimSuffix(idStr, "/read")
	notificationID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	err = repository.MarkNotificationAsRead(userID, notificationID)
	if err != nil {
		http.Error(w, "Failed to mark notification as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleJoinReqsByGroupId(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupIdStr := strings.TrimPrefix(r.URL.Path, "/api/notifications/") //   /api/notifications/joingroup/  old version
	groupIdStr = strings.TrimSuffix(groupIdStr, "/joingroup")

	notifications, statusCode := service.JoinReqsByGroupId(userId, groupIdStr)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) {
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}
