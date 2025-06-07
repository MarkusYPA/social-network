package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"backend/internal/utils"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func GetUserById(id int, viewFull bool) (model.User, error) {
	var user model.User
	var nickname sql.NullString
	var about sql.NullString
	var avatarUrl sql.NullString

	err := database.DB.QueryRow(`
		SELECT id, nickname, email, first_name, last_name, date_of_birth, about_me, avatar_path, is_public
		FROM users WHERE id = ?`, id).
		Scan(&user.ID, &nickname, &user.Email, &user.FirstName, &user.LastName, &user.Birthday, &about, &avatarUrl, &user.IsPublic)

	if nickname.Valid {
		user.Username = nickname.String
	} else {
		user.Username = ""
	}

	if about.Valid {
		user.About = about.String
	} else {
		user.About = ""
	}

	if avatarUrl.Valid {
		user.AvatarPath = avatarUrl.String
	} else {
		user.AvatarPath = ""
	}

	if !viewFull {
		user.Email = ""
		user.Birthday = ""
		user.About = ""
	}

	return user, err
}

func GetUserByEmail(req model.LoginRequest) (model.User, error) {
	var user model.User
	var nickname sql.NullString
	var about sql.NullString
	var avatarUrl sql.NullString

	err := database.DB.QueryRow(`
	SELECT id, nickname, email, first_name, last_name, date_of_birth, password_hash, about_me, avatar_path, is_public
	FROM users WHERE email = ?`, req.Email).
		Scan(&user.ID, &nickname, &user.Email, &user.FirstName, &user.LastName, &user.Birthday, &user.Password, &about, &avatarUrl, &user.IsPublic)

	if err != nil {
		fmt.Println("error getting user by email:", err)
		return user, err
	}

	if nickname.Valid {
		user.Username = nickname.String
	} else {
		user.Username = ""
	}

	if about.Valid {
		user.About = about.String
	} else {
		user.About = ""
	}

	if avatarUrl.Valid {
		user.AvatarPath = avatarUrl.String
	} else {
		user.AvatarPath = ""
	}

	return user, nil
}

func GetUserIdAndExpirationBySessionId(cookie *http.Cookie) (int, time.Time, error) {
	var userID int
	var expiresAt time.Time

	query := `SELECT user_id, expires_at FROM sessions WHERE session_token = ?`
	err := database.DB.QueryRow(query, cookie.Value).Scan(&userID, &expiresAt)

	return userID, expiresAt, err
}

func ViewFullProfileOrNot(userId, targetId int) (bool, error) {
	if userId == targetId {
		return true, nil
	}

	db := database.DB

	var isPublic bool
	err := db.QueryRow("SELECT is_public FROM users WHERE id = ?", targetId).Scan(&isPublic)
	if err != nil {
		if err == sql.ErrNoRows {
			// target user not found
			return false, fmt.Errorf("target user not found")
		}
		return false, err
	}

	if isPublic {
		return true, nil
	}

	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM follow_requests
			WHERE follower_id = ? AND followed_id = ? AND approval_status = 'accepted'
		)`, userId, targetId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// SearchUsers retrieves users whose username, first_name, or last_name match the query.
func SearchUsers(query string) ([]model.User, error) {
	q := "%" + query + "%" // Add wildcards for LIKE clause
	rows, err := database.DB.Query(`
		SELECT id, nickname, email, first_name, last_name, date_of_birth, about_me, avatar_path, is_public
		FROM users 
		WHERE nickname LIKE ? OR first_name LIKE ? OR last_name LIKE ?
		`,
		q, q, q,
	)
	if err != nil {
		fmt.Println("query error at SearchUsers:", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var nickname sql.NullString
		var about sql.NullString
		var avatarUrl sql.NullString

		err := rows.Scan(&u.ID, &nickname, &u.Email, &u.FirstName, &u.LastName, &u.Birthday, &about, &avatarUrl, &u.IsPublic)
		if err != nil {
			fmt.Println("scan error at SearchUsers:", err)
			return nil, err
		}
		if nickname.Valid {
			u.Username = nickname.String
		} else {
			u.Username = ""
		}

		if about.Valid {
			u.About = about.String
		} else {
			u.About = ""
		}

		if avatarUrl.Valid {
			u.AvatarPath = avatarUrl.String
		} else {
			u.AvatarPath = ""
		}

		users = append(users, u)
	}

	return users, nil
}

func InsertUser(passwordHash []byte, email, firstName, lastName string, parsedDOB time.Time, avatarPath, nickname, about sql.NullString) (string, int) {

	// Insert user into database
	query := `
		INSERT INTO users (
			email, password_hash, first_name, last_name,
			date_of_birth, avatar_path, nickname, about_me
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`

	_, err := database.DB.Exec(
		query,
		email, passwordHash, firstName, lastName,
		parsedDOB.Format("2006-01-02"),
		avatarPath, nickname, about,
	)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			fmt.Println("08", err)
			return "Email already in use", http.StatusConflict
		} else {
			fmt.Println("09", err)
			return "Failed to register user", http.StatusInternalServerError
		}
	}

	return "", http.StatusOK
}

func UpdateUser(userID int, data model.UpdateProfileData) (model.User, string, int) {
	query := `
		UPDATE users SET
			first_name = ?,
			last_name = ?,
			date_of_birth = ?,
			nickname = ?,
			about_me = ?,
			is_public = ?,
			updated_at = CURRENT_TIMESTAMP
	`
	args := []any{data.FirstName, data.LastName, data.DOB, utils.NullableString(data.Nickname), utils.NullableString(data.About), data.IsPublic}

	if data.DeleteAvatar {
		query += `, avatar_path = NULL`
	} else if data.AvatarPath.Valid {
		query += `, avatar_path = ?`
		args = append(args, data.AvatarPath)
	}

	query += ` WHERE id = ?`
	args = append(args, userID)

	var usr model.User

	_, err := database.DB.Exec(query, args...)
	if err != nil {
		return usr, "Failed to update user", http.StatusInternalServerError
	}

	usr, err = GetUserById(userID, true)
	if err != nil {
		return usr, "Failed to get updated user", http.StatusInternalServerError
	}

	return usr, "", http.StatusOK
}

func ProfilePrivacyByUserId(targetId int) (bool, int) {
	var isPublic bool
	err := database.DB.QueryRow("SELECT is_public FROM users WHERE id = ?", targetId).Scan(&isPublic)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, http.StatusBadRequest
		}
		return false, http.StatusInternalServerError
	}
	return isPublic, http.StatusOK
}

// GetSuggestedUsers finds non-followed non-self users who either:
// - follow or are followed by a user the active user follows or is followed by
// - are in the same group as the active user
// - follow the active user
func GetSuggestedUsers(userID int) ([]model.User, error) {

	// 1. Not userids active follows

	query := `
SELECT DISTINCT u.id, u.first_name, u.last_name, u.avatar_path, u.nickname
FROM users u
WHERE u.id != ?
AND u.id NOT IN (

	-- 0. Not users already followed
    SELECT followed_id FROM follow_requests
    WHERE follower_id = ? AND approval_status = 'accepted'
)
AND u.id IN (

    -- 1. Friends-of-friends: Followed by or following someone the user follows or is followed by
    SELECT fr2.followed_id FROM follow_requests fr1
    JOIN follow_requests fr2 ON fr1.followed_id = fr2.follower_id
    WHERE fr1.follower_id = ? AND fr1.approval_status = 'accepted' AND fr2.approval_status = 'accepted'
    
    UNION

    SELECT fr2.follower_id FROM follow_requests fr1
    JOIN follow_requests fr2 ON fr1.follower_id = fr2.followed_id
    WHERE fr1.followed_id = ? AND fr1.approval_status = 'accepted' AND fr2.approval_status = 'accepted'

    -- 2. In the same group with accepted membership
    UNION

    SELECT gm2.user_id FROM group_members gm1
    JOIN group_members gm2 ON gm1.group_id = gm2.group_id
    WHERE gm1.user_id = ? AND gm1.approval_status = 'accepted'
    AND gm2.approval_status = 'accepted'

    -- 3. Follow the user
    UNION

    SELECT follower_id FROM follow_requests
    WHERE followed_id = ? AND approval_status = 'accepted'
)
`
	rows, err := database.DB.Query(query, userID, userID, userID, userID, userID, userID)
	if err != nil {
		fmt.Println("query error at GetSuggestedUsers:", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var nickname sql.NullString
		var avatarUrl sql.NullString

		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &avatarUrl, &nickname)
		if err != nil {
			fmt.Println("scan error at GetSuggestedUsers:", err)
			return nil, err
		}

		if avatarUrl.Valid {
			u.AvatarPath = avatarUrl.String
		} else {
			u.AvatarPath = ""
		}

		if nickname.Valid {
			u.Username = nickname.String
		} else {
			u.Username = ""
		}

		users = append(users, u)
	}
	return users, nil
}
