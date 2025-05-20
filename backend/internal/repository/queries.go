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

	return user, err
}

func InsertSession(sessionID string, user model.User, expiresAt time.Time) error {
	_, err := database.DB.Exec("INSERT INTO sessions (session_token, user_id, expires_at) VALUES (?, ?, ?)", sessionID, user.ID, expiresAt)
	return err
}

func GetUserIdAndExpirationBySessionId(cookie *http.Cookie) (int, time.Time, error) {
	var userID int
	var expiresAt time.Time

	query := `SELECT user_id, expires_at FROM sessions WHERE session_token = ?`
	err := database.DB.QueryRow(query, cookie.Value).Scan(&userID, &expiresAt)

	return userID, expiresAt, err
}

func DeleteSessionById(id string) error {
	_, err := database.DB.Exec("DELETE FROM sessions WHERE session_token = ?", id)
	return err
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

/* func GetAllUsers() ([]model.User, error) {
	rows, _ := database.DB.Query("SELECT id, nickname, email, first_name, last_name, date_of_birth, about_me, avatar_path, is_public FROM users")
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var nickname sql.NullString
		var about sql.NullString
		var avatarUrl sql.NullString

		err := rows.Scan(&u.ID, &nickname, &u.Email, &u.FirstName, &u.LastName, &u.Birthday, &about, &avatarUrl, &u.IsPublic)
		if err != nil {
			return users, err
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
} */

func GetAllGroups() ([]model.Group, error) {
	rows, err := database.DB.Query(`
		SELECT id, title, description 
		FROM groups 
		WHERE status = 'enable'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []model.Group
	for rows.Next() {
		var g model.Group
		err := rows.Scan(&g.ID, &g.Title, &g.Description)
		if err != nil {
			return groups, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

func GetGroupsByUserId(userId int) ([]model.Group, error) {
	// group info from groups where that id can be found on same row as userId on group_members
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

	//fmt.Println("user groups:", groups)

	return groups, nil
}

/* func GetAllPosts() ([]model.Post, error) {
	rows, err := database.DB.Query(`
	SELECT posts.id, posts.user_id, users.first_name, users.last_name, users.avatar_path, posts.content, posts.created_at
	FROM posts
	JOIN users ON posts.user_id = users.id
	ORDER BY posts.id DESC;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		var firstname, lastname string
		var avatarUrl sql.NullString

		err := rows.Scan(&p.ID, &p.UserID, &firstname, &lastname, &avatarUrl, &p.Content, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		if avatarUrl.Valid {
			p.AvatarPath = avatarUrl.String
		} else {
			p.AvatarPath = ""
		}

		p.Username = firstname + " " + lastname
		posts = append(posts, p)
	}

	return posts, nil
} */

// GetFeedPostsBefore gets posts from userID user's follows and groups
// using cursor-based pagination: anything before the previous post (cursorTime)
// up to limit (default 10) items.
func GetFeedPostsBefore(userID int, cursorTime time.Time, limit, lastPostId int) ([]model.Post, error) {
	query := `
    -- Regular posts
    SELECT DISTINCT
        p.id,
        p.user_id,
        u.first_name,
        u.last_name,
        u.avatar_path,
        p.content,
        p.image_path,
        NULL AS group_id,
        NULL AS group_name,
        p.created_at AS created_at_sort
    FROM posts p
    JOIN users u ON p.user_id = u.id
    WHERE p.status = 'enable'
      AND (
          p.user_id = ?
          OR p.user_id IN (
              SELECT followed_id FROM follow_requests
              WHERE follower_id = ? AND approval_status = 'accepted'
          )
          )
      AND p.created_at < ?
	  AND p.id != ?
        
    UNION ALL
    
    -- Group posts
    SELECT DISTINCT
        gp.id,
        gp.user_id,
        u.first_name,
        u.last_name,
        u.avatar_path,
        gp.content,
        gp.image_path,
        gp.group_id,
        g.title AS group_name,
        gp.created_at AS created_at_sort
    FROM group_posts gp
    JOIN group_members gm ON gp.group_id = gm.group_id
        AND gm.user_id = ? AND gm.approval_status = 'accepted'
    JOIN groups g ON gp.group_id = g.id
    JOIN users u ON gp.user_id = u.id
    WHERE gp.status = 'enable'
      AND gp.created_at < ?
	  AND gp.id != ?
    
    ORDER BY created_at_sort DESC
    LIMIT ?;`

	rows, err := database.DB.Query(query, userID, userID, cursorTime, lastPostId, userID, cursorTime, lastPostId, limit)
	if err != nil {
		fmt.Println("query err at GetFeedPostsBefore:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var firstname, lastname string
		var avatarUrl sql.NullString

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&firstname,
			&lastname,
			&avatarUrl,
			&post.Content,
			&post.ImagePath,
			&post.GroupID,
			&post.GroupName,
			&post.CreatedAt,
		)
		if err != nil {
			fmt.Println("scan rows err at GetFeedPostsBefore:", err)
			return nil, err
		}
		if avatarUrl.Valid {
			post.AvatarPath = avatarUrl.String
		} else {
			post.AvatarPath = ""
		}
		post.Username = firstname + " " + lastname
		posts = append(posts, post)
	}

	return posts, nil
}

func GetPostsByUserId(targetId int) ([]model.Post, error) {
	rows, err := database.DB.Query(`
	SELECT posts.id, posts.user_id, users.first_name, users.last_name, users.avatar_path, posts.content, posts.created_at
	FROM posts
	JOIN users ON posts.user_id = users.id
	WHERE posts.user_id = ?
	ORDER BY posts.id DESC;`, targetId)

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

		err := rows.Scan(&p.ID, &p.UserID, &firstname, &lastname, &avatarUrl, &p.Content, &p.CreatedAt)
		if err != nil {
			fmt.Println("scan error at GetPostsByUserId", err)
			return nil, err
		}

		if avatarUrl.Valid {
			p.AvatarPath = avatarUrl.String
		} else {
			p.AvatarPath = ""
		}

		p.Username = firstname + " " + lastname
		posts = append(posts, p)
	}

	return posts, nil
}

func InsertPost(userID int, content string, privacy string) (int, string, error) {
	result, err := database.DB.Exec(
		"INSERT INTO posts (user_id, content, privacy_level) VALUES (?, ?, ?)",
		userID, content, privacy,
	)

	if err != nil {
		fmt.Println("error 1 at insert post", err)
		return 0, "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("error 2 at insert post", err)
		return 0, "", err
	}

	var createdAt string
	_ = database.DB.QueryRow("SELECT created_at FROM posts WHERE id = ?", id).Scan(&createdAt)

	return int(id), createdAt, nil
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

	fmt.Println("data at UpdateUser:", data)

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
