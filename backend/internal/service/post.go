package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/utils"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func CreatePost(content, privacyLvl string, imagePath *string, userID int, viewerIDs []int) (model.Post, int) {

	var post model.Post
	id, createdAt, err := repository.InsertPost(userID, content, privacyLvl, imagePath)
	if err != nil {
		return post, http.StatusInternalServerError
	}

	if privacyLvl == "private" {
		err := repository.AddViewersToPrivatePost(id, viewerIDs)
		if err != nil {
			return post, http.StatusInternalServerError
		}
	}

	usr, err := repository.GetUserById(userID, true)
	if err != nil {
		return post, http.StatusNotFound
	}

	post.ID = id
	post.UserID = userID
	post.Username = usr.FirstName + " " + usr.LastName
	post.AvatarPath = usr.AvatarPath
	post.ImagePath = imagePath
	post.Content = content
	post.CreatedAt = createdAt
	post.PostType = "regular"
	post.Privacy = &privacyLvl
	return post, http.StatusOK
}

func SaveUploadedFile(file multipart.File, header *multipart.FileHeader, filePath string) (string, error) {

	ext := filepath.Ext(header.Filename)

	if ext == "" || !utils.IsAllowedImageExtension(ext) {
		fmt.Println("bad extension:", ext)
		return "", fmt.Errorf("illegal extension")
	}

	filename := uuid.New().String() + filepath.Ext(header.Filename)
	dst := filepath.Join("data/uploads/", filePath, filename)

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return "/" + dst, nil
}

func GetFeedPosts(cursor, limitStr, lastPostIdStr string, userId int) ([]model.Post, error) {
	limit := 10
	lastPostId := 0
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if lpid, err := strconv.Atoi(lastPostIdStr); err == nil && lpid > 0 {
		lastPostId = lpid
	}

	cursorTime := time.Now().UTC() // use current time by default

	if cursor != "" {
		t, err := time.Parse(time.RFC3339, cursor)
		if err == nil {
			cursorTime = t
		}
	}

	cursorTime = cursorTime.Truncate(time.Second)

	posts, err := repository.GetFeedPostsBefore(userId, cursorTime, limit, lastPostId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func CommentsForPost(postIDstring, postType string, userID int) ([]model.Comment, error) {
	PostID, err := strconv.Atoi(postIDstring)
	if err != nil {
		fmt.Println("can not convert postId")
		return nil, err
	}

	var comments []model.Comment
	if postType == "regular" {
		comments, err = repository.ReadAllCommentsForPost(PostID, userID)
	} else if postType == "group" {
		comments, err = repository.ReadAllCommentsForGroupPost(PostID, userID)
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return comments, nil
}

func CreateCommentsForPost(UserID int, PostIDstring, postType, content string, imagePath *string) error {
	PostID, err := strconv.Atoi(PostIDstring)
	if err != nil {
		return err
	}

	if postType == "regular" {
		err = repository.InsertComment(content, UserID, PostID, imagePath)
	} else if postType == "group" {
		err = repository.InsertGroupComment(content, UserID, PostID, imagePath)
	}

	return err
}
