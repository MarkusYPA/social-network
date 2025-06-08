package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"database/sql"
	"fmt"
	// "database/sql"
	// "fmt"
	// "net/http"
)

func CreateEvent(event model.Event) (int, error) {
	query := `INSERT INTO events (group_id, creator_id, title, description, event_datetime)
	          VALUES (?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, event.GroupID, event.CreatorID, event.Title, event.Description, event.EventDate)
	if err != nil {
		return 0, err
	}
	eventID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Loop through the group members, crete event responses and send notifications
	members, err := GetGroupMembersByGroupId(int(event.GroupID))
	if err != nil {
		return 0, err
	}
	for _, member := range members {
		err = SaveEventResponse(int(eventID), member.ID, "pending")
		if err != nil {
			return 0, err
		}
		// send notifications
		_, err = InsertNotification(event.CreatorID, member.ID, "event_creation", int(eventID))
		if err != nil {
			return 0, err
		}
	}

	return int(eventID), err
}

func SaveEventResponse(eventID, userID int, response string) error {
	query := `
	INSERT INTO event_responses (event_id, user_id, response)
	VALUES (?, ?, ?)
	ON CONFLICT(event_id, user_id) DO UPDATE SET response = ?`

	_, err := database.DB.Exec(query, eventID, userID, response, response)
	return err
}

func GetEventResponse(eventID, userID int) (string, error) {
	query := `
	SELECT er.response
	FROM event_responses er
	WHERE er.event_id = ? AND er.user_id = ?`

	var response string
	err := database.DB.QueryRow(query, eventID, userID).Scan(&response)

	if err == sql.ErrNoRows {
		return "", nil
	}
	return response, err
}

func GetEventByID(eventID int) (model.Event, error) {
	var e model.Event

	query := `
        SELECT 
            e.id, e.group_id, g.title as group_title, 
            e.creator_id, u.first_name, u.last_name,
            e.title, e.description, e.event_datetime
        FROM events e
        LEFT JOIN groups g ON e.group_id = g.id
        LEFT JOIN users u ON e.creator_id = u.id
        WHERE e.id = ? AND e.status = 'enable'
    `
	err := database.DB.QueryRow(query, eventID).Scan(
		&e.ID, &e.GroupID, &e.Group,
		&e.CreatorID, &e.Creator.FirstName, &e.Creator.LastName,
		&e.Title, &e.Description, &e.EventDate)
	if err != nil {
		return e, err
	}

	responses, err := database.DB.Query(`
        SELECT 
            u.id, u.first_name, u.last_name, er.response
        FROM event_responses er
        JOIN users u ON er.user_id = u.id
        WHERE er.event_id = ?`, eventID)
	if err != nil {
		return e, err
	}
	defer responses.Close()

	for responses.Next() {
		var user model.User
		var response string
		err := responses.Scan(&user.ID, &user.FirstName, &user.LastName, &response)
		if err != nil {
			return e, err
		}

		switch response {
		case "going":
			e.Going = append(e.Going, user)
		case "not_going":
			e.NotGoing = append(e.NotGoing, user)
		case "pending":
			e.NoResponse = append(e.NoResponse, user)
		}
	}

	return e, nil
}

func CheckUserGroupMembership(userID, groupID int) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM group_members
			WHERE user_id = ? AND group_id = ? AND status = 'enable' AND approval_status = 'accepted'
		)
	`
	err := database.DB.QueryRow(query, userID, groupID).Scan(&exists)
	return exists, err
}

func getUsersByResponse(eventID int, response string) ([]model.User, error) {
	query := `
	SELECT u.id, u.first_name, u.last_name, u.avatar_path
	FROM event_responses er
	JOIN users u ON er.user_id = u.id
	WHERE er.event_id = ? AND er.response = ? AND er.status = 'enable' AND u.status = 'enable'
	`

	rows, err := database.DB.Query(query, eventID, response)
	if err != nil {
		fmt.Println("query err at getUsersByResponse:", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var avatarUrl sql.NullString

		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &avatarUrl)
		if err != nil {
			return nil, err
		}

		if avatarUrl.Valid {
			u.AvatarPath = avatarUrl.String
		} else {
			u.AvatarPath = ""
		}
		users = append(users, u)
	}
	//fmt.Println(len(users), "users at event", eventID, "that are", response)
	return users, nil
}

func GetEventsByUser(userID int) ([]model.Event, error) {
	query := `
	SELECT DISTINCT e.id, g.title, e.group_id, e.creator_id, u.first_name, u.last_name, u.avatar_path, e.title, e.description, e.event_datetime
	FROM events e
	JOIN groups g ON e.group_id = g.id
	JOIN users u ON e.creator_id = u.id
	LEFT JOIN event_responses er ON er.event_id = e.id
	WHERE e.creator_id = ? OR er.user_id = ?
	ORDER BY e.event_datetime DESC
	`
	rows, err := database.DB.Query(query, userID, userID)
	if err != nil {
		fmt.Println("query error at GetEventsByUser:", err)
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var e model.Event
		var avatarUrl sql.NullString
		err := rows.Scan(&e.ID, &e.Group, &e.GroupID, &e.CreatorID, &e.Creator.FirstName, &e.Creator.LastName, &avatarUrl, &e.Title, &e.Description, &e.EventDate)
		if err != nil {
			fmt.Println("scan error at GetEventsByUser:", err)
			return nil, err
		}

		if avatarUrl.Valid {
			e.Creator.AvatarPath = avatarUrl.String
		} else {
			e.Creator.AvatarPath = ""
		}

		if e.ID == nil {
			continue
		}

		// Query responses and attach users
		goings, err := getUsersByResponse(*e.ID, "going")
		if err != nil {
			fmt.Println("error getting going users at GetEventsByUser:", err)
			return nil, err
		}
		notGoings, err := getUsersByResponse(*e.ID, "not_going")
		if err != nil {
			fmt.Println("error getting not going users at GetEventsByUser:", err)
			return nil, err
		}
		noResponses, err := getUsersByResponse(*e.ID, "pending")
		if err != nil {
			fmt.Println("error getting not responded users at GetEventsByUser:", err)
			return nil, err
		}

		e.Going = goings
		e.NotGoing = notGoings
		e.NoResponse = noResponses

		events = append(events, e)
	}
	return events, nil
}

func GetEventsByGroup(groupID, userID int) ([]model.Event, error) {
	query := `
	SELECT DISTINCT e.id, g.title, e.group_id, e.creator_id, u.first_name, u.last_name, u.avatar_path, e.title, e.description, e.event_datetime
	FROM events e
	JOIN groups g ON e.group_id = g.id
	JOIN users u ON e.creator_id = u.id
	LEFT JOIN event_responses er ON er.event_id = e.id
	WHERE e.group_id = ?
	  AND e.status = 'enable'
	  AND (e.creator_id = ? OR er.user_id = ?)
	ORDER BY e.event_datetime DESC
	`

	rows, err := database.DB.Query(query, groupID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var e model.Event
		var avatarUrl sql.NullString
		err := rows.Scan(&e.ID, &e.Group, &e.GroupID, &e.CreatorID, &e.Creator.FirstName, &e.Creator.LastName, &avatarUrl, &e.Title, &e.Description, &e.EventDate)
		if err != nil {
			return nil, err
		}

		if avatarUrl.Valid {
			e.Creator.AvatarPath = avatarUrl.String
		} else {
			e.Creator.AvatarPath = ""
		}

		if e.ID == nil {
			continue
		}

		// Query responses and attach users
		goings, err := getUsersByResponse(*e.ID, "going")
		if err != nil {
			return nil, err
		}
		notGoings, err := getUsersByResponse(*e.ID, "not_going")
		if err != nil {
			return nil, err
		}
		noResponses, err := getUsersByResponse(*e.ID, "pending")
		if err != nil {
			return nil, err
		}

		e.Going = goings
		e.NotGoing = notGoings
		e.NoResponse = noResponses

		events = append(events, e)
	}
	return events, nil
}
