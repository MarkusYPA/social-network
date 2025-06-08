package model

import (
	"database/sql"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	ID         int        `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Birthday   string     `json:"birthday"` // Birthday can be used to check if returned profile is partly hidden or not
	Password   string     `json:"password"`
	About      string     `json:"about_me"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	UpdatedBy  *int       `json:"updated_by"`
	Status     string     `json:"status"`
	AvatarPath string     `json:"avatar_url"`
	IsPublic   bool       `json:"is_public"`
}

type Post struct {
	ID               int     `json:"id"`
	UserID           int     `json:"user_id"`
	Username         string  `json:"username"`
	AvatarPath       string  `json:"avatar_url"`
	Content          string  `json:"content"`
	ImagePath        *string `json:"image_path,omitempty"`
	GroupID          *int    `json:"group_id,omitempty"` // nil for regular posts
	GroupName        *string `json:"group_name,omitempty"`
	CreatedAt        string  `json:"created_at"`
	NumberOfComments int     `json:"numberOfComments"`
	PostType         string  `json:"postType"`
	Privacy          *string `json:"privacy,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileData struct {
	FirstName    string
	LastName     string
	DOB          string
	Nickname     string
	About        string
	AvatarPath   sql.NullString
	DeleteAvatar bool
	IsPublic     bool
}

type Group struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Comment struct {
	ID               int        `json:"id"`
	PostId           int        `json:"post_id"`
	CommentId        int        `json:"comment_id"`
	Content          string     `json:"content"`
	UserId           int        `json:"user_id"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	UpdatedBy        *int       `json:"updated_by"`
	IsLikedByUser    bool       `json:"liked"`
	IsDislikedByUser bool       `json:"disliked"`
	NumberOfLikes    int        `json:"number_of_likes"`
	NumberOfDislikes int        `json:"number_of_dislikes"`
	ImagePath        *string    `json:"image_path,omitempty"`
	User             User       `json:"user"`
	RepliesCount     int        `json:"repliesCount"`
	ISCreatedByMe    bool       `json:"isCreatedByMe"`
}

type FollowRequest struct {
	TargetID int    `json:"target_id"`
	Action   string `json:"action"` // "request", "follow", "unfollow"
}

type Event struct {
	ID          *int   `json:"id,omitempty"`
	Group       string `json:"group"`
	GroupID     int    `json:"group_id"`
	CreatorID   int    `json:"creator_id"`
	Creator     User   `json:"creator"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EventDate   string `json:"event_datetime"`
	//EventDate   time.Time `json:"event_datetime"`	// trying if time.Time type works
	Going      []User `json:"going"`
	NotGoing   []User `json:"not_going"`
	NoResponse []User `json:"no_response"`
}

type EventResponse struct {
	EventID  int    `json:"event_id"`
	UserID   int    `json:"user_id"`
	Response string `json:"response"` // going / not_going / pending
}

type GroupRequest struct {
	TargetID int    `json:"target_id"`
	Action   string `json:"action"` // "request", "leave", "delete"
}

type GroupRequestApproval struct {
	GroupID     int `json:"group_id"`
	RequesterID int `json:"requester_id"`
}

type GroupInvitation struct {
	ID      int    `json:"id"`
	UserId  int    `json:"user_id"`
	GroupID string `json:"group_id"`
	Inviter int    `json:"inviter_id"`
}

type InvitableUser struct {
	User       User   `json:"user"`
	Membership string `json:"membership"`
}

type Notification struct {
	ID            int     `json:"id"`
	Type          string  `json:"type"` // 'follow_request', 'group_invitation', 'group_join_request', 'event_creation'
	UserID        int     `json:"user_id"`
	SenderID      *int    `json:"sender_id,omitempty"`
	SenderName    *string `json:"sender_name,omitempty"`
	FollowReqID   *int    `json:"follow_req_id,omitempty"`
	GroupInviteID *int    `json:"group_invite_id,omitempty"`
	GroupID       *int    `json:"group_id,omitempty"`
	GroupTitle    *string `json:"group_title,omitempty"`
	EventID       *int    `json:"event_id,omitempty"`
	EventTitle    *string `json:"event_title,omitempty"`
	Content       *string `json:"content,omitempty"`
	IsRead        *bool   `json:"is_read,omitempty"`
	Pending       bool    `json:"pending"`
	CreatedAt     string  `json:"created_at"`
}

type WSMessage struct {
	Type     string `json:"type"`
	From     string `json:"from"`
	FromName string `json:"from_name,omitempty"`
	To       string `json:"receiver_id,omitempty"`
	Content  string `json:"content,omitempty"`
}

type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan WSMessage // individual send channel
}

type ChatMessage struct {
	ID         int    `json:"id"`
	SenderName string `json:"sender_name"`
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

type Chat struct {
	IsActive bool          `json:"is_active"` // is there a follow relation between the users
	Name     string        `json:"name"`
	UserID   string        `json:"user_id"` // in case there are no messages
	Messages []ChatMessage `json:"messages"`
}

var (
	Broadcast = make(chan WSMessage)     // for global broadcast
	Clients   = make(map[string]*Client) // userID -> client
	Mu        sync.Mutex
)
