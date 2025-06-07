package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func Membership(userId, groupId int) (string, error) {
	approval, adminId, err := repository.GroupMembership(userId, groupId)
	if err != nil {
		fmt.Println("error at Membership:", err)
		return "", err
	}

	if userId == adminId {
		return "admin", nil
	}

	return approval, nil
}

func CreateGroup(group model.Group, userId int) (int, error) {
	var err error
	group.ID, err = repository.CreateGroup(group, userId)
	if err != nil {
		return 0, err
	}
	err = repository.AddGroupMember(userId, group.ID)
	return group.ID, err
}

// get invite info: group id and user id
func ApproveGroupInvitation(userID int, data []string) int {

	action := data[1]
	inviteID, err := strconv.Atoi(data[0])
	if err != nil {
		return http.StatusBadRequest
	}

	// get invite info: group id and user id
	groupID, invitedID, err := repository.GetGroupInviteInfo(inviteID)
	if err != nil {
		fmt.Println("Failed to get group invitation info", err)
		return http.StatusInternalServerError
	}

	if userID != invitedID {
		fmt.Println("Only own invitations could be approved")
		return http.StatusUnauthorized
	}

	membership, err := Membership(userID, groupID)
	if err != nil {
		fmt.Println("Failed to determine group membership status", err)
		return http.StatusInternalServerError
	}

	var statusCode int
	switch membership {
	case "accepted":
		// User is already a member, just update invite status
		statusCode = repository.UpdateGroupInviteStatus(inviteID, userID, action)
	case "pending", "declined", "":
		// Update both invite and member status
		statusCode = repository.UpdateGroupInviteStatus(inviteID, userID, action)
		if action == "accepted" {
			err = repository.AddGroupMember(userID, groupID)
			if err != nil {
				fmt.Println("Failed to add group member", err)
				return http.StatusInternalServerError
			}
		}

	default:
		fmt.Println("Invalid membership status", err)
		return http.StatusInternalServerError
	}

	return statusCode
}

func GetUsersMembership(invitables []model.InvitableUser, users []model.User, groupId int) ([]model.InvitableUser, error) {
	for _, user := range users {
		var inv model.InvitableUser
		inv.User = user
		membership, err := repository.GetMembershipStatus(user.ID, groupId)
		if err != nil {
			fmt.Println("error getting membership at HandleGroupInvitationSearch:", err)
			return nil, err
		}
		inv.Membership = membership
		invitables = append(invitables, inv)
	}
	return invitables, nil
}

func DeleteGroup(userID, targetID int) int {
	err := repository.DeleteGroupWithDependencies(userID, targetID)
	if err != nil {
		fmt.Println("Error deleting group and dependencies:", err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func PostsByGroupId(userId, targetId int) ([]model.Post, error) {
	viewGroup, err := repository.ViewFullGroupOrNot(userId, targetId)
	if err != nil {
		return nil, err
	}

	var posts []model.Post

	if viewGroup {
		posts, err = repository.GetGroupPostsByGroupId(targetId)
		if err != nil {
			return nil, err
		}
	}

	return posts, err
}

func MembersByGroupId(userId, targetId int) ([]model.User, error) {
	viewGroup, err := repository.ViewFullGroupOrNot(userId, targetId)
	if err != nil {
		return nil, err
	}

	var users []model.User

	if viewGroup {
		users, err = repository.GetGroupMembersByGroupId(targetId)
		if err != nil {
			return nil, err
		}
	}

	return users, err
}

func EventsByGroupId(userId, targetId int) ([]model.Event, error) {
	viewGroup, err := repository.ViewFullGroupOrNot(userId, targetId)
	if err != nil {
		return nil, err
	}

	var events []model.Event

	if viewGroup {
		events, err = repository.GetGroupEventsByGroupId(targetId)
		if err != nil {
			return nil, err
		}
	}

	return events, err
}

func ValidMembership(userID int, groupIDStr string) (int, error) {
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		fmt.Println("Invalid group_id at CreateGroupPostHandler", err)
		return 0, err
	}

	membership, err := Membership(userID, groupID)
	if err != nil {
		return 0, err
	}

	if membership != "accepted" && membership != "admin" {
		return 0, err
	}

	return groupID, nil
}

func CreateGroupPost(userID, groupID int, content string, imagePath *string) (model.Post, int) {
	var post model.Post

	id, createdAt, err := repository.InsertGroupPost(userID, groupID, content, imagePath)
	if err != nil {
		return post, http.StatusInternalServerError
	}

	user, err := repository.GetUserById(userID, true)
	if err != nil {
		return post, http.StatusNotFound
	}

	group, err := repository.GetGroupById(groupID)
	if err != nil {
		return post, http.StatusNotFound
	}

	post = model.Post{
		ID:         int(id),
		UserID:     userID,
		Username:   user.FirstName + " " + user.LastName,
		AvatarPath: user.AvatarPath,
		Content:    content,
		ImagePath:  imagePath,
		GroupID:    &groupID,
		GroupName:  &group.Title,
		CreatedAt:  createdAt,
		PostType:   "group",
	}

	return post, http.StatusOK
}

func GroupMembership(userID int, req model.GroupRequest) int {
	var statusCode int

	switch req.Action {
	case "request":
		var gmId int
		gmId, statusCode = repository.GroupRequest(userID, req.TargetID) // 'approval_status' to pending, 'status' to enable
		if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) {
			fmt.Println("error code at HandleFollowAction:", statusCode)
			return statusCode
		}

		adminID, err := repository.GetAdminIdByGroupId(req.TargetID)
		if err != nil {
			fmt.Println("error getting admin in HandleGroupMembership:", err)
			return http.StatusBadRequest
		}

		_, err = repository.InsertNotification(userID, adminID, "group_join_request", gmId) // last id needs to be id at group members table
		if err != nil {
			return http.StatusBadRequest
		}
	case "leave":
		statusCode = repository.LeaveGroup(userID, req.TargetID) // 'status' to delete
	case "cancel": // User (userID) cancels their join request to group (req.TargetID)
		notificationID, err := repository.RemoveGroupRequestNotification(userID, req.TargetID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error removing group join request notification: %v", err)
			return http.StatusInternalServerError
		}

		if err == sql.ErrNoRows {
			log.Printf("No active group join request notification found for user %d to group %d.", userID, req.TargetID)
		}

		// Proceed to remove the group join request (leave group)
		statusCode = repository.LeaveGroup(userID, req.TargetID) // req.TargetID is groupID
		if statusCode != http.StatusOK {
			log.Printf("Error leaving group (cancelling request): status code %d for user %d, group %d", statusCode, userID, req.TargetID)
			return statusCode
		}

		// If notification was successfully marked as deleted, send WS message to group admin
		if err == nil { // sql.ErrNoRows would mean err != nil
			adminID, adminErr := repository.GetAdminIdByGroupId(req.TargetID)
			if adminErr != nil {
				log.Printf("Error getting admin ID for group %d: %v", req.TargetID, adminErr)
				// Not sending HTTP error here as the main operation (LeaveGroup) succeeded.
				// But important to log for observability.
			} else {
				wsErr := SendNotificationDeletedWS(notificationID, adminID)
				if wsErr != nil {
					log.Printf("Error sending group join request deleted WebSocket message for notification %d to admin %d of group %d: %v", notificationID, adminID, req.TargetID, wsErr)
					// Do not typically send HTTP error for WebSocket issues
				}
			}
		}
	case "delete":
		statusCode = DeleteGroup(userID, req.TargetID)
	default:
		statusCode = http.StatusBadRequest
	}

	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		fmt.Println("error code at HandleFollowAction:", statusCode)
		return statusCode
	}

	return http.StatusOK
}

func GroupById(userId, targetId int) (model.Group, string, error) {
	group, err := repository.GetGroupById(targetId)
	if err != nil {
		return group, "", err
	}

	membership, err := Membership(userId, targetId)
	if err != nil {
		return group, "", err
	}

	return group, membership, nil
}

func GroupRequestApprove(userID int, req model.GroupRequestApproval, action string) int {
	groupID := req.GroupID
	requesterID := req.RequesterID

	membership, err := Membership(userID, groupID)
	if err != nil {
		return http.StatusInternalServerError
	}

	if membership != "admin" {
		return http.StatusUnauthorized
	}

	statusCode := repository.ApproveGroupRequest(requesterID, groupID, userID, action)

	return statusCode
}

func GroupInvitation(userID int, groupInvite model.GroupInvitation) int {
	groupInvite.Inviter = userID

	var err error
	groupInvite.ID, err = repository.InviteToGroup(groupInvite)
	if err != nil {
		return http.StatusUnauthorized
	}

	_, err = repository.InsertNotification(groupInvite.Inviter, groupInvite.UserId, "group_invitation", groupInvite.ID)
	if err != nil {
		return http.StatusBadRequest
	}

	return http.StatusOK
}

func GroupInvitationSearch(query, groupIdStr string, invitables []model.InvitableUser) ([]model.InvitableUser, int) {
	users, err := repository.SearchUsers(query)
	if err != nil {
		return invitables, http.StatusInternalServerError
	}

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		fmt.Println("strconv.Atoi error at HandleGroupInvitationSearch:", err)
		return invitables, http.StatusBadRequest
	}

	for _, user := range users {
		var inv model.InvitableUser
		inv.User = user
		membership, err := repository.GetMembershipStatus(user.ID, groupId)
		if err != nil {
			fmt.Println("error getting membership at HandleGroupInvitationSearch:", err)
			return invitables, http.StatusInternalServerError
		}
		inv.Membership = membership
		invitables = append(invitables, inv)
	}

	return invitables, http.StatusOK
}
