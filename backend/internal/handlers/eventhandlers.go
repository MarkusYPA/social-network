package handlers

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"encoding/json"
	"fmt"

	// "fmt"
	"net/http"
	"strconv"
	"strings"
)

func HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var event model.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		fmt.Println("decoding error at HandleCreateEvent:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	event.CreatorID = userID
	event.Creator, err = repository.GetUserById(userID, false)
	if err != nil {
		fmt.Println("error getting creator at HandleCreateEvent:", err)
		http.Error(w, "Failed to get user data", http.StatusInternalServerError)
		return
	}
	group, err := repository.GetGroupById(event.GroupID)
	if err != nil {
		fmt.Println("error getting group at HandleCreateEvent:", err)
		http.Error(w, "Failed to get group data", http.StatusInternalServerError)
		return
	}
	event.Group = group.Title

	fmt.Println("Event before creation:", event)

	id, err := repository.CreateEvent(event)
	if err != nil {
		fmt.Println("error creating event at HandleCreateEvent:", err)
		http.Error(w, "Could not create event", http.StatusInternalServerError)
		return
	}

	event.ID = &id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func HandleEventResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var resp model.EventResponse
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	resp.UserID = userID

	if resp.Response != "going" && resp.Response != "not_going" && resp.Response != "pending" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	oldResponse, err := repository.GetEventResponse(resp.EventID, resp.UserID)
	if err != nil {
		http.Error(w, "Could not get previous response", http.StatusInternalServerError)
		return
	}

	if oldResponse == resp.Response { // remove old response when clicking same button
		resp.Response = "pending"
	}

	err = repository.SaveEventResponse(resp.EventID, resp.UserID, resp.Response)
	if err != nil {
		fmt.Println("Saving event response failed:", err)
		http.Error(w, "Could not save response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetHandleEventByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract the event ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/events/")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	// Check if the user is a member of the group
	event, err := repository.GetEventByID(eventID)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	isMember, err := repository.CheckUserGroupMembership(userID, event.GroupID)
	if err != nil {
		http.Error(w, "Failed to check group membership", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func GetEventsByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// idStr := strings.TrimPrefix(r.URL.Path, "/api/events/user/")
	// userID, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
	// 	return
	// }

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	events, err := repository.GetEventsByUser(userID)
	if err != nil {
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func GetEventsByGroupID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupStr := strings.TrimPrefix(r.URL.Path, "/api/events/group/")
	groupID, err := strconv.Atoi(groupStr)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Check if the user is a member of the group
	isMember, err := repository.CheckUserGroupMembership(userID, groupID)
	if err != nil {
		http.Error(w, "Failed to check group membership", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	events, err := repository.GetEventsByGroup(groupID, userID)
	if err != nil {
		http.Error(w, "Failed to fetch events for group", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
