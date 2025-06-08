package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"database/sql"

	"fmt"
)

func ReadAllCommentsForPost(postID int, userID int) ([]model.Comment, error) {
	var comments []model.Comment

	selectQuery := `
SELECT 
    u.id AS user_id,
    u.first_name,
    u.last_name,
    u.created_at AS user_created_at,
    u.updated_at AS user_updated_at,
    u.updated_by AS user_updated_by,
	u.avatar_path,
    c.id AS comment_id,
    c.post_id AS comment_post_id,
    c.user_id AS comment_user_id,
    c.content AS comment_description, 
    c.status AS comment_status,
    c.created_at AS comment_created_at, 
    c.updated_by AS comment_updated_by,
	c.image_path
FROM comments c
INNER JOIN users u ON c.user_id = u.id
WHERE 
    c.status != 'delete' 
    AND u.status != 'delete'
    AND c.post_id = ?
ORDER BY c.created_at DESC;
`

	rows, err := database.DB.Query(selectQuery, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment model.Comment
		var user model.User
		var avatarUrl sql.NullString

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.UpdatedBy,
			//&user.AvatarPath,
			&avatarUrl,

			&comment.ID,
			&comment.PostId,
			&comment.UserId,
			&comment.Content,
			&comment.Status,
			&comment.CreatedAt,
			&comment.UpdatedBy,
			&comment.ImagePath,
		)
		if err != nil {
			return nil, err
		}

		if avatarUrl.Valid {
			user.AvatarPath = avatarUrl.String
		} else {
			user.AvatarPath = ""
		}

		comment.User = user
		comment.ISCreatedByMe = (user.ID == userID)

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func InsertComment(content string, userID int, postID int, imagePath *string) error {

	var query string
	var args []any

	if imagePath != nil {
		query = "INSERT INTO comments (user_id, content, post_id, image_path) VALUES (?, ?, ?, ?)"
		args = []any{userID, content, postID, *imagePath}
	} else {
		query = "INSERT INTO comments (user_id, content,  post_id) VALUES (?, ?, ?)"
		args = []any{userID, content, postID}
	}

	_, err := database.DB.Exec(query, args...)
	if err != nil {
		fmt.Println("error 1 at insert comment", err)
		return err
	}

	return nil
}

func ReadAllCommentsForGroupPost(postID int, userID int) ([]model.Comment, error) {
	var comments []model.Comment

	selectQuery := `
SELECT 
    u.id AS user_id,
    u.first_name,
    u.last_name,
    u.created_at AS user_created_at,
    u.updated_at AS user_updated_at,
    u.updated_by AS user_updated_by,
	u.avatar_path,
    c.id AS comment_id,
    c.group_post_id AS comment_post_id,
    c.user_id AS comment_user_id,
    c.content AS comment_description, 
    c.status AS comment_status,
    c.created_at AS comment_created_at, 
    c.updated_by AS comment_updated_by,
	c.image_path
FROM group_comments c
INNER JOIN users u ON c.user_id = u.id
WHERE 
    c.status != 'delete' 
    AND u.status != 'delete'
    AND c.group_post_id = ?
ORDER BY c.created_at DESC;
`

	rows, err := database.DB.Query(selectQuery, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment model.Comment
		var user model.User

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.UpdatedBy,
			&user.AvatarPath,

			&comment.ID,
			&comment.PostId,
			&comment.UserId,
			&comment.Content,
			&comment.Status,
			&comment.CreatedAt,
			&comment.UpdatedBy,
			&comment.ImagePath,
		)
		if err != nil {
			return nil, err
		}

		comment.User = user
		comment.ISCreatedByMe = (user.ID == userID)

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	fmt.Println("image path: ", comments[0].User.AvatarPath)
	return comments, nil
}

func InsertGroupComment(content string, userID int, postID int, imagePath *string) error {

	var query string
	var args []any

	if imagePath != nil {
		query = "INSERT INTO group_comments (user_id, content, group_post_id, image_path) VALUES (?, ?, ?, ?)"
		args = []any{userID, content, postID, *imagePath}
	} else {
		query = "INSERT INTO group_comments (user_id, content, group_post_id) VALUES (?, ?, ?)"
		args = []any{userID, content, postID}
	}

	_, err := database.DB.Exec(query, args...)
	if err != nil {
		fmt.Println("error 1 at insert group_comment", err)
		return err
	}

	return nil
}
