package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"database/sql"
	"fmt"
	"time"
)

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
		p.privacy_level AS privacy,
        NULL AS group_id,
        NULL AS group_name,
        p.created_at AS created_at_sort,
    	COUNT(c.id) AS comment_count,
    	'regular' AS post_type
    FROM posts p
    JOIN users u ON p.user_id = u.id
	LEFT JOIN comments c ON c.post_id = p.id AND c.status != 'delete'
    WHERE p.status = 'enable'
      AND (
	  	  -- own posts
          p.user_id = ?

		  -- non-private posts from followed users
          OR (
		      p.privacy_level != 'private'
		      AND p.user_id IN (
                SELECT followed_id FROM follow_requests
                WHERE follower_id = ? AND approval_status = 'accepted'
              )
            )
		  
		  -- private posts
          OR (
              p.privacy_level = 'private'
			  AND p.user_id IN (
                SELECT followed_id FROM follow_requests
                WHERE follower_id = ? AND approval_status = 'accepted'
                )
              AND EXISTS (
                  SELECT 1 FROM post_privacy pp
                  WHERE pp.post_id = p.id
                    AND pp.user_id = ?
                    AND pp.status = 'enable'
                )
            )
        )

      AND p.created_at < ?
	  AND p.id != ?
    GROUP BY p.id, u.id    
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
		NULL AS privacy,
        gp.group_id,
        g.title AS group_name,
        gp.created_at AS created_at_sort,
    	COUNT(gc.id) AS comment_count,
    	'group' AS post_type
    FROM group_posts gp
    JOIN group_members gm ON gp.group_id = gm.group_id
        AND gm.user_id = ? AND gm.approval_status = 'accepted'
    JOIN groups g ON gp.group_id = g.id
    JOIN users u ON gp.user_id = u.id
	LEFT JOIN group_comments gc ON gc.group_post_id = gp.id AND gc.status != 'delete'
    WHERE gp.status = 'enable'
      AND gp.created_at < ?
	  AND gp.id != ?
    GROUP BY gp.id, u.id, g.id
    ORDER BY created_at_sort DESC
    LIMIT ?;`

	rows, err := database.DB.Query(query, userID, userID, userID, userID, cursorTime, lastPostId, userID, cursorTime, lastPostId, limit)
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
			&post.Privacy,
			&post.GroupID,
			&post.GroupName,
			&post.CreatedAt,
			&post.NumberOfComments,
			&post.PostType,
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

func GetPostsByUserId(userId, targetId int) ([]model.Post, error) {
	rows, err := database.DB.Query(`
	SELECT
		p.id,
		p.user_id,
		u.first_name,
		u.last_name,
		u.avatar_path,
		p.content,
		p.created_at AS created_at_sort,
		p.image_path,
		COUNT(c.id) AS comment_count,
		'regular' AS post_type,
		p.privacy_level AS privacy,
		NULL AS group_id,
        NULL AS group_name
	FROM posts p
	JOIN users u ON p.user_id = u.id
	LEFT JOIN comments c ON c.post_id = p.id AND c.status != 'delete'
	WHERE p.status = 'enable' AND p.user_id = ?
	      AND (
		  -- posts on own profile
		    p.user_id = ?

		  -- public posts
		  OR
		    p.privacy_level = 'public'	

		  -- almost private posts from followed users
          OR (
		      p.privacy_level = 'almost_private'
		      AND p.user_id IN (
                SELECT followed_id FROM follow_requests
                WHERE follower_id = ? AND approval_status = 'accepted'
              )
            )
		  
		  -- private posts
          OR (
              p.privacy_level = 'private'
			  AND p.user_id IN (
                SELECT followed_id FROM follow_requests
                WHERE follower_id = ? AND approval_status = 'accepted'
                )
              AND EXISTS (
                  SELECT 1 FROM post_privacy pp
                  WHERE pp.post_id = p.id
                    AND pp.user_id = ?
                    AND pp.status = 'enable'
                )
            )
        )
	GROUP BY p.id, u.id		
	UNION ALL
    
    -- Group posts
    SELECT DISTINCT
        gp.id,
        gp.user_id,
        u.first_name,
        u.last_name,
        u.avatar_path,
        gp.content,
        gp.created_at AS created_at_sort,
        gp.image_path,
    	COUNT(gc.id) AS comment_count,
    	'group' AS post_type,
		NULL AS privacy,
        gp.group_id,
        g.title AS group_name
    FROM group_posts gp
    JOIN group_members gm ON gp.group_id = gm.group_id
        AND gm.user_id = ? AND gm.approval_status = 'accepted'		-- from groups where active user is member
    JOIN groups g ON gp.group_id = g.id
    JOIN users u ON gp.user_id = u.id
	LEFT JOIN group_comments gc ON gc.group_post_id = gp.id AND gc.status != 'delete'
    WHERE gp.status = 'enable'  AND gp.user_id = ?      			-- posts made by target user
    GROUP BY gp.id, u.id
    ORDER BY created_at_sort DESC`, targetId, userId, userId, userId, userId, userId, targetId)

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
		err := rows.Scan(&p.ID, &p.UserID, &firstname, &lastname, &avatarUrl, &p.Content, &p.CreatedAt, &p.ImagePath, &p.NumberOfComments, &p.PostType, &p.Privacy, &p.GroupID, &p.GroupName)
		if err != nil {
			fmt.Println("scan error at GetPostsByUserId", err)
			return nil, err
		}
		//p.PostType = "regular"
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

func InsertPost(userID int, content string, privacy string, imagePath *string) (int, string, error) {
	var query string
	var args []any

	if imagePath != nil {
		query = "INSERT INTO posts (user_id, content, privacy_level, image_path) VALUES (?, ?, ?, ?)"
		args = []any{userID, content, privacy, *imagePath}
	} else {
		query = "INSERT INTO posts (user_id, content, privacy_level) VALUES (?, ?, ?)"
		args = []any{userID, content, privacy}
	}

	result, err := database.DB.Exec(query, args...)
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
	err = database.DB.QueryRow("SELECT created_at FROM posts WHERE id = ?", id).Scan(&createdAt)
	if err != nil {
		fmt.Println("error 3 at insert post (fetching created_at)", err)
		return 0, "", err
	}

	return int(id), createdAt, nil
}

func AddViewersToPrivatePost(postId int, viewerIDs []int) error {
	query := `INSERT INTO post_privacy (post_id, user_id) VALUES (?, ?)`
	for _, viewer := range viewerIDs {
		_, err := database.DB.Exec(query, postId, viewer)
		if err != nil {
			fmt.Println("error inserting viewers for private post:", err)
			return err
		}
	}
	return nil
}
