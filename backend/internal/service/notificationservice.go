package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// SendNotificationDeletedWS sends a WebSocket message to notify a user that a notification has been deleted.
func SendNotificationDeletedWS(notificationID int64, recipientUserID int) error {
	content := map[string]string{"id": strconv.FormatInt(notificationID, 10)}
	jsonContent, err := json.Marshal(content)
	if err != nil {
		log.Printf("Error marshalling notification deletion content: %v", err)
		return err
	}

	wsMsg := model.WSMessage{
		Type:    "notification_deleted",
		To:      strconv.Itoa(recipientUserID),
		From:    "system_notification_deletion",
		Content: string(jsonContent),
	}

	select {
	case model.Broadcast <- wsMsg:
		log.Printf("Notification deletion message for notification ID %d queued for user %s.", notificationID, wsMsg.To)
	default:
		log.Printf("Broadcast channel full or unavailable. Notification deletion message for notification ID %d for user %s not sent via WebSocket immediately.", notificationID, wsMsg.To)
		// Depending on requirements, you might return an error here:
		// return fmt.Errorf("broadcast channel full, message for notification %d not sent", notificationID)
	}

	return nil
}

func GetNotifications(userID int) ([]model.Notification, int) {
	notifications, err := repository.GetAllNotificatons(userID)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	for i := range notifications {
		var pending bool
		var err error

		switch notifications[i].Type {
		case "follow_request":
			pending, err = repository.CheckFollowRequestStatus(*notifications[i].FollowReqID)
		case "group_invitation":
			pending, err = repository.CheckInvitationStatus(*notifications[i].GroupInviteID)
		case "group_join_request":
			pending, err = repository.CheckJoinRequestStatus(*notifications[i].SenderID, *notifications[i].GroupID)
		case "event_creation":
			pending, err = repository.CheckEventInvitationStatus(notifications[i].UserID, *notifications[i].EventID)
		}

		if err != nil {
			return nil, http.StatusInternalServerError
		}
		notifications[i].Pending = pending
	}

	return notifications, http.StatusOK
}

func JoinReqsByGroupId(userId int, groupIdStr string) ([]model.Notification, int) {
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		fmt.Println("Error converting id at HandleJoinReqsByGroupId", err)
		return nil, http.StatusBadRequest
	}

	membership, err := Membership(userId, groupId)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	var notifications []model.Notification
	if membership == "admin" {
		// get notifications for people wanting to join this group
		notifications, err = repository.GetJoinRequests(groupId)

		if err != nil {
			return nil, http.StatusInternalServerError
		}
	}

	return notifications, http.StatusOK
}
