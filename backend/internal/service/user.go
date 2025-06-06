package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/utils"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func UserById(userId, targetId int) (model.User, int) {
	var usr model.User

	getFull, err := repository.ViewFullProfileOrNot(userId, targetId)
	if err != nil {
		return usr, http.StatusBadRequest
	}

	usr, err = repository.GetUserById(targetId, getFull)
	if err != nil {
		return usr, http.StatusNotFound
	}

	return usr, http.StatusOK
}

func RegisterUser(userInfo struct {
	Email      string
	Password   string
	FirstName  string
	LastName   string
	DOB        string
	Nickname   string
	About      string
	AvatarPath sql.NullString
}) (string, int) {

	// Parse and validate date of birth
	parsedDOB, err := time.Parse("2006-01-02", userInfo.DOB)
	if err != nil {
		fmt.Println("03", err)
		return "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("04", err)
		return "Failed to hash password", http.StatusInternalServerError
	}

	errMsg, statusCode := repository.InsertUser(passwordHash, userInfo.Email, userInfo.FirstName, userInfo.LastName, parsedDOB, userInfo.AvatarPath, utils.NullableString(userInfo.Nickname), utils.NullableString(userInfo.About))
	if errMsg != "" {
		return errMsg, statusCode
	}

	return "", http.StatusOK
}

func UpdateUserProfile(userID int, updateData model.UpdateProfileData) (model.User, string, int) {

	if updateData.Nickname == "null" {
		updateData.Nickname = ""
	}

	if updateData.About == "null" {
		updateData.About = ""
	}

	var usr model.User

	if strings.TrimSpace(updateData.FirstName) == "" || strings.TrimSpace(updateData.LastName) == "" || strings.TrimSpace(updateData.DOB) == "" {
		return usr, "required fields missing", http.StatusBadRequest
	}

	usr, errMsg, statusCode := repository.UpdateUser(userID, updateData)
	if errMsg != "" {
		return usr, errMsg, statusCode
	}

	return usr, "", http.StatusOK
}

func UploadAvatar(file multipart.File, header *multipart.FileHeader) (sql.NullString, error) {
	var avatarPath sql.NullString

	ext := filepath.Ext(header.Filename)
	if ext == "" || !utils.IsAllowedImageExtension(ext) {
		fmt.Println("bad extension:", ext)
		return avatarPath, fmt.Errorf("illegal extension")
	}
	filename := fmt.Sprintf("data/uploads/avatars/%d%s", time.Now().UnixNano(), ext)

	dst, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating path at uploadAvatar:", err)
		return avatarPath, err
	}
	defer dst.Close()

	io.Copy(dst, file)

	avatarPath.Valid = true
	avatarPath.String = filename
	return avatarPath, nil
}

func GetFollowSatatus(userId, targetId int) (string, int) {

	if userId == targetId {
		return "me", http.StatusOK
	}

	isPublic, statusCode := repository.ProfilePrivacyByUserId(targetId)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		return "", statusCode
	}

	approval, statusCode := repository.FollowApproval(userId, targetId)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		return "", statusCode
	}

	if approval == "accepted" { // no difference between private and public profile
		return "accepted", http.StatusOK
	}

	if approval == "pending" { // always private profile
		return "pending", http.StatusOK
	}

	// no difference between declined and no status
	if isPublic {
		return "not following public", http.StatusOK
	}
	return "not following private", http.StatusOK
}
