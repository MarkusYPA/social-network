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

func HandleSuggestGroups(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groups, err := repository.GetRecommendedGroups(userId)
	if err != nil {
		http.Error(w, "Could not fetch recommended groups", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(groups)
}

func SearchGroups(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		json.NewEncoder(w).Encode([]model.Group{}) // Return empty array for empty query
		return
	}

	groups, err := repository.SearchGroups(query)
	if err != nil {
		http.Error(w, "Error searching groups", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(groups)
}

func HandlePostsByGroupId(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/group/posts/")
	targetId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	posts, err := service.PostsByGroupId(userId, targetId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func HandleMembersByGroupId(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/group/members/")
	targetId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	users, err := service.MembersByGroupId(userId, targetId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func HandleEventsByGroupId(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/group/events/")
	targetId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	events, err := service.EventsByGroupId(userId, targetId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

func CreateGroupPostHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ValidateSession(r) // replace with your actual session logic
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		fmt.Println("error reading data at CreateGroupPostHandler", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	groupIDStr := r.FormValue("group_id")
	if content == "" || groupIDStr == "" {
		fmt.Println("Missing fields at CreateGroupPostHandler", err)
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	groupID, err := service.ValidMembership(userID, groupIDStr)
	if err != nil {
		http.Error(w, "No valid membership found", http.StatusBadRequest)
		return
	}

	var imagePath *string
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		savedPath, saveErr := service.SaveUploadedFile(file, header, "posts")
		if saveErr != nil {
			http.Error(w, "Error saving image", http.StatusInternalServerError)
			return
		}
		imagePath = &savedPath
	} else if err != http.ErrMissingFile {
		fmt.Println("Error reading file at CreateGroupPostHandler", err)
		http.Error(w, "Error reading image", http.StatusInternalServerError)
		return
	}

	post, statusCode := service.CreateGroupPost(userID, groupID, content, imagePath)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) {
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func HandleGroupMembership(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed at HandleGroupMembership")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("ValidateSession error at HandleGroupMembership:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req model.GroupRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("json error at HandleGroupMembership:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	statusCode := service.GroupMembership(userID, req)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) {
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleGroupsByUserId(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("validate error in HandleGroupsByUserId:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groups, err := repository.GetGroupsByUserId(userId)
	if err != nil {
		http.Error(w, "Failed to fetch groups", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(groups)
}

func HandleGroupRequests(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("validate error in HandleGroupRequestsByUserId:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groups, err := repository.GetGroupRequestsByUserId(userId)
	if err != nil {
		http.Error(w, "Failed to fetch group requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(groups)
}

func HandleGroupInvitations(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("validate error in HandleGroupInvitationsByUserId:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groups, err := repository.GetGroupInvitationsByUserId(userId)
	if err != nil {
		http.Error(w, "Failed to fetch group invitations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(groups)
}

func HandleGroupsAdministered(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("validate error in HandleGroupInvitationsByUserId:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groups, err := repository.GetGroupsAdministeredByUserId(userId)
	if err != nil {
		http.Error(w, "Failed to fetch group invitations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(groups)
}

func HandleGroupById(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("validate error in HandleGroupsByUserId:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/group/")
	targetId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	group, membership, err := service.GroupById(userId, targetId)
	if err != nil {
		http.Error(w, "Failed to fetch group", http.StatusInternalServerError)
		return
	}

	resp := struct {
		Group      model.Group `json:"group"`
		Membership string      `json:"membership"`
	}{
		Group:      group,
		Membership: membership,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func HandleCreateGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("validate error in HandleGroupsByUserId:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var grp model.Group
	err = json.NewDecoder(r.Body).Decode(&grp)
	if err != nil {
		fmt.Println("json error at HandleCreateGroup:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	grp.ID, err = service.CreateGroup(grp, userId)
	if err != nil {
		fmt.Println("CreateGroup error in HandleCreateGroup:", err)
		http.Error(w, "failed to create group", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(grp)
}

func HandleGroupRequestApprove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed at HandleFollowRequestApprove")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("ValidateSession error at HandleFollowRequestApprove:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req model.GroupRequestApproval
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body at HandleFollowRequestApprove", http.StatusBadRequest)
		return
	}

	action := strings.TrimPrefix(r.URL.Path, "/api/group/requests/")
	if action != "accepted" && action != "declined" {
		http.Error(w, "Invalid request action syntax", http.StatusBadRequest)
		return
	}

	statusCode := service.GroupRequestApprove(userID, req, action)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		fmt.Println("error code at HandleGroupRequestApprove:", statusCode)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleGroupInvitation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed at HandleFollowRequestApprove")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var groupInvite model.GroupInvitation
	err = json.NewDecoder(r.Body).Decode(&groupInvite)
	if err != nil {
		fmt.Println("json error at HandleGroupInvitation:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	statusCode := service.GroupInvitation(userID, groupInvite)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		fmt.Println("error code at HandleGroupRequestApprove:", statusCode)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleGroupInvitationSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed at HandleGroupInvitationSearch")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var invitables []model.InvitableUser
	query := r.URL.Query().Get("query")
	if query == "" {
		json.NewEncoder(w).Encode(invitables) // Return empty array for empty query
		return
	}

	type req struct {
		GroupId string `json:"group_id"`
	}
	var groupInfo req
	err = json.NewDecoder(r.Body).Decode(&groupInfo)
	if err != nil {
		fmt.Println("json error at HandleGroupInvitationSearch:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	invitables, statusCode := service.GroupInvitationSearch(query, groupInfo.GroupId, invitables)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		fmt.Println("error code at HandleGroupRequestApprove:", statusCode)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invitables)
}

func HandleApproveGroupInvitation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed at HandleFollowRequestApprove")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/group/invite/"), "/")
	if len(data) != 2 {
		http.Error(w, "Invalid request action syntax", http.StatusBadRequest)
		return
	}

	statusCode := service.ApproveGroupInvitation(userID, data)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		fmt.Println("error code at ApproveGroupInvitation:", statusCode)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}
