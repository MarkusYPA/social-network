package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

func FollowRequest(followerID, followedID int) (int, int) {
	if followerID == followedID {
		return 0, http.StatusBadRequest
	}

	isPublic, statusCode := repository.ProfilePrivacyByUserId(followedID)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		return 0, statusCode
	}
	var frID int
	var err error

	if !isPublic {
		frID, err = repository.SendFollowRequest(followerID, followedID)
		if err != nil {
			log.Println("Error inserting/updating follow request:", err)
			return 0, http.StatusInternalServerError
		}
	} else {
		return 0, http.StatusBadRequest
	}

	return frID, http.StatusOK
}

func Follow(followerID, followedID int) int {
	if followerID == followedID {
		return http.StatusBadRequest
	}

	isPublic, statusCode := repository.ProfilePrivacyByUserId(followedID)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		return statusCode
	}

	if isPublic {
		err := repository.StartToFollow(followerID, followedID)
		if err != nil {
			log.Println("Error inserting/updating follow:", err)
			return http.StatusInternalServerError
		}
	} else {
		return http.StatusBadRequest
	}

	return http.StatusOK
}

func FollowAction(userID int, req model.FollowRequest) int {
	var err error
	var statusCode, frID int

	switch req.Action {
	case "request":
		frID, statusCode = FollowRequest(userID, req.TargetID)
		if statusCode == http.StatusOK {
			_, err = repository.InsertNotification(userID, req.TargetID, "follow_request", frID)
			if err != nil {
				return http.StatusBadRequest
			}
		}
	case "follow":
		statusCode = Follow(userID, req.TargetID)
	case "unfollow":
		statusCode = repository.RemoveFollow(userID, req.TargetID)
	case "cancel":
		notificationID, err := repository.RemoveFollowRequestNotification(userID, req.TargetID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error removing follow request notification: %v", err)
			return http.StatusInternalServerError
		}

		if err == sql.ErrNoRows {
			log.Printf("No active follow request notification found for user %d to target %d. Proceeding to remove follow request.", userID, req.TargetID)
		}

		// Proceed to remove the follow request itself
		statusCode = repository.RemoveFollow(userID, req.TargetID)
		if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) {
			log.Printf("Error removing follow request: status code %d", statusCode)
			return statusCode
		}

		// If notification was successfully marked as deleted, send WS message
		if err == nil { // sql.ErrNoRows would mean err != nil
			wsErr := SendNotificationDeletedWS(notificationID, req.TargetID)
			if wsErr != nil {
				log.Printf("Error sending notification deleted WebSocket message for notification %d to user %d: %v", notificationID, req.TargetID, wsErr)
				// Do not typically send HTTP error for WebSocket issues
			}
		}
	default:
		return http.StatusBadRequest
	}

	return statusCode
}

func FollowRequestApprove(userID int, data []string) int {

	action := data[1]
	requestID, err := strconv.Atoi(data[0])
	if err != nil {
		return http.StatusBadRequest
	}

	var statusCode int
	switch action {
	case "accept":
		statusCode = repository.AcceptFollowRequest(userID, requestID)
	case "decline":
		statusCode = repository.DeclineFollowRequest(userID, requestID)
	default:
		statusCode = http.StatusBadRequest
	}

	return statusCode
}
