package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"database/sql"
	"fmt"
	"net/http"
)

func ViewFullGroupOrNot(userId, targetId int) (bool, error) {
	db := database.DB
	var grId int
	err := db.QueryRow(`
	SELECT id 
	FROM group_members 
	WHERE group_id = ? AND user_id = ? AND approval_status = 'accepted' AND status = 'enable'`, targetId, userId).Scan(&grId)
	if err != nil {
		if err == sql.ErrNoRows { // No membership found so no permission to view full group
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GroupMembership(userId, targetId int) (string, int, error) {
	db := database.DB

	var approval string
	var adminId int
	err := db.QueryRow(`
	SELECT gm.approval_status, g.creator_id 
	FROM group_members gm
	JOIN groups g ON gm.group_id = g.id
	WHERE gm.group_id = ? AND gm.user_id = ? AND gm.status = 'enable' AND g.status = 'enable'`, targetId, userId).Scan(&approval, &adminId)
	if err != nil {
		if err == sql.ErrNoRows { // No membership found
			//fmt.Println("no membership found for user", userId, "in group", targetId)
			return approval, adminId, nil
		}
		return approval, adminId, err
	}
	return approval, adminId, nil
}

func GetRecommendedGroups(userID int) ([]model.Group, error) {
	query := `
		SELECT DISTINCT g.id, g.title, g.description
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.approval_status = 'accepted'
		  AND gm.user_id IN (
		      SELECT followed_id FROM follow_requests
		      WHERE follower_id = ? AND approval_status = 'accepted'
		      UNION
		      SELECT follower_id FROM follow_requests
		      WHERE followed_id = ? AND approval_status = 'accepted'
		      UNION
		      SELECT gm2.user_id
		      FROM group_members gm1
		      JOIN group_members gm2 ON gm1.group_id = gm2.group_id
		      WHERE gm1.user_id = ? AND gm1.approval_status = 'accepted'
		        AND gm2.approval_status = 'accepted'
		        AND gm2.user_id != ?
		  )
		  AND g.status = 'enable'
		  AND g.id NOT IN (
		      SELECT group_id FROM group_members
		      WHERE user_id = ? AND approval_status = 'accepted'
		  );
	`

	rows, err := database.DB.Query(query, userID, userID, userID, userID, userID)
	if err != nil {
		fmt.Println("query error at GetRecommendedGroups", err)
		return nil, err
	}
	defer rows.Close()

	var groups []model.Group
	for rows.Next() {
		var g model.Group
		if err := rows.Scan(&g.ID, &g.Title, &g.Description); err != nil {
			fmt.Println("scan error at GetRecommendedGroups", err)
			return groups, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

func SearchGroups(query string) ([]model.Group, error) {
	searchTerm := "%" + query + "%"
	sqlQuery := `
		SELECT id, title, description
		FROM groups
		WHERE status = 'enable'
		  AND (title LIKE ? OR description LIKE ?)
		ORDER BY
		  CASE
		    WHEN title LIKE ? THEN 0
		    ELSE 1
		  END,
		  title ASC;
	`

	rows, err := database.DB.Query(sqlQuery, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []model.Group
	for rows.Next() {
		var g model.Group
		if err := rows.Scan(&g.ID, &g.Title, &g.Description); err != nil {
			return groups, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

func GetGroupById(groupId int) (model.Group, error) {
	var group model.Group

	err := database.DB.QueryRow(`
	SELECT id, title, description
	FROM groups
	WHERE id = ? AND status = 'enable'`, groupId).Scan(&group.ID, &group.Title, &group.Description)

	if err != nil {
		fmt.Println("error getting user by email:", err)
		return group, err
	}

	return group, nil
}

func GetGroupPostsByGroupId(groupId int) ([]model.Post, error) {
	rows, err := database.DB.Query(`
	SELECT gp.id, gp.user_id, gp.image_path, gp.content, gp.created_at, u.first_name, u.last_name, u.avatar_path, COUNT(gc.id) AS comment_count
	FROM group_posts gp
	JOIN users u ON gp.user_id = u.id
	LEFT JOIN group_comments gc ON gc.group_post_id = gp.id AND gc.status != 'delete'
    WHERE gp.status = 'enable'
	AND gp.group_id = ?
	GROUP BY gp.id
	ORDER BY gp.id DESC;`, groupId)

	if err != nil {
		fmt.Println("rows error at GetPostsByUserId", err)
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		var firstname, lastname string
		var avatarUrl sql.NullString

		err := rows.Scan(&p.ID, &p.UserID, &p.ImagePath, &p.Content, &p.CreatedAt, &firstname, &lastname, &avatarUrl, &p.NumberOfComments)
		if err != nil {
			fmt.Println("scan error at GetPostsByUserId", err)
			return nil, err
		}

		if avatarUrl.Valid {
			p.AvatarPath = avatarUrl.String
		} else {
			p.AvatarPath = ""
		}
		p.PostType = "group"
		p.Username = firstname + " " + lastname
		posts = append(posts, p)
	}

	return posts, nil
}

func GetGroupMembersByGroupId(groupId int) ([]model.User, error) {
	rows, err := database.DB.Query(`
	SELECT u.id, u.first_name, u.last_name, u.avatar_path
	FROM users u
	JOIN group_members gm ON u.id = gm.user_id
	WHERE gm.group_id = ? AND gm.status = 'enable' AND u.status = 'enable' AND approval_status = 'accepted';`, groupId)

	if err != nil {
		fmt.Println("rows error at GetGroupMembersByGroupId", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var firstname, lastname string
		var avatarUrl sql.NullString

		err := rows.Scan(&u.ID, &firstname, &lastname, &avatarUrl)
		if err != nil {
			fmt.Println("scan error at GetPostsByUserId", err)
			return nil, err
		}

		if avatarUrl.Valid {
			u.AvatarPath = avatarUrl.String
		} else {
			u.AvatarPath = ""
		}

		u.Username = firstname + " " + lastname
		users = append(users, u)
	}

	return users, nil
}

func GetGroupEventsByGroupId(groupId int) ([]model.Event, error) {
	rows, err := database.DB.Query(`
	SELECT e.id, g.title, e.group_id, e.title, e.description, e.event_datetime
	FROM events e
	JOIN groups g ON e.group_id = g.id
	WHERE e.group_id = ?;`, groupId)

	if err != nil {
		fmt.Println("rows error at GetGroupMembersByGroupId", err)
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var e model.Event
		err := rows.Scan(&e.ID, &e.Group, &e.GroupID, &e.Title, &e.Description, &e.EventDate)
		if err != nil {
			fmt.Println("scan error at GetPostsByUserId", err)
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func InsertGroupPost(userID, groupID int, content string, imagePath *string) (int64, string, error) {
	query := `
		INSERT INTO group_posts (user_id, group_id, content, image_path)
		VALUES (?, ?, ?, ?)
	`
	// Use Exec for INSERT statements when you don't expect rows to be returned directly.
	result, err := database.DB.Exec(query, userID, groupID, content, imagePath)
	if err != nil {
		return 0, "", fmt.Errorf("failed to execute insert: %w", err)
	}

	// Get the last inserted row ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, "", fmt.Errorf("failed to get last insert ID: %w", err)
	}

	var createdAt string
	row := database.DB.QueryRow("SELECT created_at FROM group_posts WHERE id = ?", id)
	err = row.Scan(&createdAt)
	if err != nil {
		return id, "", fmt.Errorf("failed to get created_at for ID %d: %w", id, err)
	}

	return id, createdAt, nil
}

func GroupRequest(userID, groupID int) (int, int) {
	query := `
	INSERT INTO group_members (group_id, user_id, approval_status, status)
	VALUES (?, ?, 'pending', 'enable')
	ON CONFLICT(group_id, user_id) DO UPDATE SET
		approval_status = 'pending',
		status = 'enable',
		updated_by = ?
	`
	_, err := database.DB.Exec(query, groupID, userID, userID)
	if err != nil {
		return 0, http.StatusInternalServerError
	}

	var id int
	row := database.DB.QueryRow("SELECT id FROM group_members WHERE group_id = ? AND user_id = ?", groupID, userID)
	err = row.Scan(&id)
	if err != nil {
		return 0, http.StatusInternalServerError
	}

	return id, http.StatusOK
}

func LeaveGroup(userID, groupID int) int {
	var creatorID int
	err := database.DB.QueryRow("SELECT creator_id FROM groups WHERE id = ?", groupID).Scan(&creatorID)
	if err != nil {
		fmt.Println("Error checking group creator:", err)
		return http.StatusInternalServerError
	}
	if creatorID == userID {
		// creator/admin should not leave the group
		return http.StatusForbidden
	}

	query := `
		UPDATE group_members
		SET status = 'delete', updated_by = ?
		WHERE group_id = ? AND user_id = ?
	`
	res, err := database.DB.Exec(query, userID, groupID, userID)
	if err != nil {
		fmt.Println("Error at LeaveGroup:", err)
		return http.StatusInternalServerError
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return http.StatusNotFound
	}
	return http.StatusOK
}

func DeleteGroupWithDependencies(userID, groupID int) error {

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Soft delete related group_members
	_, err = tx.Exec(`UPDATE group_members SET status = 'delete', updated_by = ? WHERE group_id = ?`, userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete group_members: %w", err)
	}

	// Soft delete related group_invitations
	_, err = tx.Exec(`UPDATE group_invitations SET status = 'delete', updated_by = ? WHERE group_id = ?`, userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete group_invitations: %w", err)
	}

	// Soft delete related group_posts
	_, err = tx.Exec(`UPDATE group_posts SET status = 'delete', updated_by = ? WHERE group_id = ?`, userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete group_posts: %w", err)
	}

	// Soft delete related group_comments (via posts)
	_, err = tx.Exec(`
		UPDATE group_comments 
		SET status = 'delete', updated_by = ? 
		WHERE group_post_id IN (SELECT id FROM group_posts WHERE group_id = ?)`, userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete group_comments: %w", err)
	}

	// Soft delete related events
	_, err = tx.Exec(`UPDATE events SET status = 'delete', updated_by = ? WHERE group_id = ?`, userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete events: %w", err)
	}

	// Soft delete event_responses for those events
	_, err = tx.Exec(`
		UPDATE event_responses 
		SET status = 'delete', updated_by = ? 
		WHERE event_id IN (SELECT id FROM events WHERE group_id = ?)`, userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete event_responses: %w", err)
	}

	// Soft delete group_messages
	_, err = tx.Exec(`UPDATE group_messages SET status = 'delete', updated_by = ? WHERE group_id = ?`, userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete group_messages: %w", err)
	}

	// Soft delete relevant notifications
	_, err = tx.Exec(`
		UPDATE notifications 
		SET status = 'delete', updated_by = ? 
		WHERE 
			(type = 'group_invitation' AND group_invite_id IN (
				SELECT id FROM group_invitations WHERE group_id = ?
			))
			OR 
			(type = 'group_join_request' AND group_members_id IN (
				SELECT id FROM group_members WHERE group_id = ?
			))`, userID, groupID, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete notifications: %w", err)
	}

	// Soft delete the group itself
	_, err = tx.Exec(`UPDATE groups SET status = 'delete', updated_by = ? WHERE id = ? AND creator_id = ?`, userID, groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func GetAdminIdByGroupId(groupId int) (int, error) {
	var adminId int
	err := database.DB.QueryRow(`
	SELECT creator_id
	FROM groups
	WHERE id =?`, groupId).Scan(&adminId)
	if err != nil {
		return 0, err
	}
	return adminId, nil
}

func CreateGroup(group model.Group, userId int) (int, error) {
	query := `
	INSERT INTO groups (creator_id, title, description)
	VALUES (?, ?, ?)`

	res, err := database.DB.Exec(query, userId, group.Title, group.Description)
	if err != nil {
		fmt.Println("exec error creating group:", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func AddGroupMember(userId, groupId int) error {
	query := `
	INSERT INTO group_members (group_id, user_id, approval_status)
	VALUES (?, ?, 'accepted')	
	ON CONFLICT (group_id, user_id) DO UPDATE SET
		status = 'enable',
		updated_by = ?
	`
	_, err := database.DB.Exec(query, groupId, userId, userId)

	if err != nil {
		fmt.Println("exec error adding group member:", err)
	}
	return err
}

func ApproveGroupRequest(userID, groupID, adminID int, action string) int {
	query := `
		UPDATE group_members SET
			approval_status = ?,
			updated_by = ?
		WHERE 
			status = 'enable' AND 
			approval_status = 'pending' AND 
			user_id = ? AND 
			group_id = ?
	`
	_, err := database.DB.Exec(query, action, adminID, userID, groupID)
	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func InviteToGroup(groupInvite model.GroupInvitation) (int, error) {
	query := `
INSERT INTO group_invitations (group_id, user_id, inviter_id, approval_status)
VALUES (?, ?, ?, 'pending')
ON CONFLICT(group_id, user_id) DO UPDATE SET
    approval_status = 'pending',
    status = 'enable',
    updated_by = EXCLUDED.inviter_id
WHERE
    approval_status != 'accepted'
    OR 
	status = 'delete';`

	_, err := database.DB.Exec(
		query,
		groupInvite.GroupID,
		groupInvite.UserId,
		groupInvite.Inviter,
	)
	if err != nil {
		return 0, err
	}

	query = `SELECT id FROM group_invitations WHERE group_id = ? AND user_id = ?` // last insert won't work when update at on conflict
	var groupInviteId int
	err = database.DB.QueryRow(query, groupInvite.GroupID, groupInvite.UserId).Scan(&groupInviteId)

	return int(groupInviteId), err
}

func GetMembershipStatus(userID, groupID int) (string, error) {
	// Check group_members table
	var memberStatus string
	memberQuery := `
		SELECT approval_status 
		FROM group_members 
		WHERE user_id = ? AND group_id = ? AND status = 'enable'
	`
	err := database.DB.QueryRow(memberQuery, userID, groupID).Scan(&memberStatus)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to query group_members: %w", err)
	}

	// Check group_invitations table
	var invitationStatus string
	invitationQuery := `
		SELECT approval_status 
		FROM group_invitations 
		WHERE user_id = ? AND group_id = ? AND status = 'enable'
	`
	err = database.DB.QueryRow(invitationQuery, userID, groupID).Scan(&invitationStatus)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to query group_invitations: %w", err)
	}

	// Apply business logic
	if memberStatus == "accepted" {
		return "accepted", nil
	}

	if invitationStatus == "pending" {
		return "invited", nil
	}

	return "", nil
}

func GetGroupInviteInfo(inviteID int) (int, int, error) {
	// Check group_members table
	var groupID, invitedID int
	memberQuery := `
		SELECT group_id, user_id
		FROM group_invitations
		WHERE id = ? AND status = 'enable'
	`
	err := database.DB.QueryRow(memberQuery, inviteID).Scan(&groupID, &invitedID)
	if err != nil && err != sql.ErrNoRows {
		return 0, 0, fmt.Errorf("failed to query group invitation info: %w", err)
	}

	return groupID, invitedID, nil
}

func UpdateGroupInviteStatus(inviteID, userID int, action string) int {
	query := `
		UPDATE group_invitations SET
			approval_status = ?,
			updated_by = ?,
			status = 'delete'  
		WHERE 
		id = ? AND
		status = 'enable'  
	`
	_, err := database.DB.Exec(query, action, userID, inviteID)
	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func RemoveGroupRequestNotification(userID, groupID int) (int64, error) {
	var notificationID int64
	err := database.DB.QueryRow(`
		SELECT id FROM notifications
		WHERE type = 'group_join_request'
		AND ref_id = (
			SELECT id FROM group_members
			WHERE user_id = ? AND group_id = ?
		)
		AND status != 'delete'
	`, userID, groupID).Scan(&notificationID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("notification not found: %w", err)
		}
		fmt.Println("Error retrieving notification ID:", err)
		return 0, err
	}

	query := `
	UPDATE notifications
	SET updated_by = ?, status = 'delete'
	WHERE id = ?
	`
	_, err = database.DB.Exec(query, userID, notificationID)

	if err != nil {
		fmt.Println("error removing notification at RemoveGroupRequestNotification", err)
		return 0, err
	}

	return notificationID, nil
}

func GetGroupsByUserId(userId int) ([]model.Group, error) {
	query := `
	SELECT g.id, g.title, g.description
	FROM groups g
	JOIN group_members gm ON g.id = gm.group_id
	WHERE gm.user_id = ? AND g.status = 'enable' AND gm.status = 'enable' AND gm.approval_status = 'accepted' ;
	`

	rows, err := database.DB.Query(query, userId)
	if err != nil {
		fmt.Println("query error in GetGroupsByUserId", err)
		return nil, err
	}
	defer rows.Close()

	var groups []model.Group
	for rows.Next() {
		var g model.Group
		err := rows.Scan(&g.ID, &g.Title, &g.Description)
		if err != nil {
			fmt.Println("scan error in GetGroupsByUserId", err)
			return groups, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

func GetGroupRequestsByUserId(userId int) ([]model.Group, error) {
	query := `
	SELECT g.id, g.title, g.description
	FROM groups g
	JOIN group_members gm ON g.id = gm.group_id
	WHERE gm.user_id = ? AND g.status = 'enable' AND gm.status = 'enable' AND gm.approval_status = 'pending' ;
	`

	rows, err := database.DB.Query(query, userId)
	if err != nil {
		fmt.Println("query error in GetGroupRequestsByUserId", err)
		return nil, err
	}
	defer rows.Close()

	var groups []model.Group
	for rows.Next() {
		var g model.Group
		err := rows.Scan(&g.ID, &g.Title, &g.Description)
		if err != nil {
			fmt.Println("scan error in GetGroupRequestsByUserId", err)
			return groups, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

func GetGroupInvitationsByUserId(userId int) ([]model.Group, error) {
	query := `
	SELECT g.id, g.title, g.description
	FROM groups g
	JOIN group_invitations gi ON g.id = gi.group_id
	WHERE gi.user_id = ? AND g.status = 'enable' AND gi.status = 'enable' AND gi.approval_status = 'pending' ;
	`

	rows, err := database.DB.Query(query, userId)
	if err != nil {
		fmt.Println("query error in GetGroupInvitationsByUserId", err)
		return nil, err
	}
	defer rows.Close()

	var groups []model.Group
	for rows.Next() {
		var g model.Group
		err := rows.Scan(&g.ID, &g.Title, &g.Description)
		if err != nil {
			fmt.Println("scan error in GetGroupInvitationsByUserId", err)
			return groups, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

func GetGroupsAdministeredByUserId(userId int) ([]model.Group, error) {
	query := `
	SELECT g.id, g.title, g.description
	FROM groups g
	WHERE g.creator_id = ? AND g.status = 'enable';
	`

	rows, err := database.DB.Query(query, userId)
	if err != nil {
		fmt.Println("query error in GetGroupsAdministeredByUserId", err)
		return nil, err
	}
	defer rows.Close()

	var groups []model.Group
	for rows.Next() {
		var g model.Group
		err := rows.Scan(&g.ID, &g.Title, &g.Description)
		if err != nil {
			fmt.Println("scan error in GetGroupsAdministeredByUserId", err)
			return groups, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}
