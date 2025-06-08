package main

import (
	"backend/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/service"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func setHandlers() {
	// CORS to allow developement on same address
	http.HandleFunc("/api/users/", middleware.WithCORS(handlers.HandleUserByID))
	http.HandleFunc("/api/users/search", middleware.WithCORS(handlers.SearchUsers))
	http.HandleFunc("/api/posts/", middleware.WithCORS(handlers.HandlePostsByUserId))
	http.HandleFunc("/api/posts/create", middleware.WithCORS(handlers.HandleCreatePost))
	http.HandleFunc("/api/group/posts/", middleware.WithCORS(handlers.HandlePostsByGroupId))
	http.HandleFunc("/api/group/members/", middleware.WithCORS(handlers.HandleMembersByGroupId))
	http.HandleFunc("/api/group/events/", middleware.WithCORS(handlers.HandleEventsByGroupId))
	http.HandleFunc("/api/homefeed", middleware.WithCORS(handlers.HandleGetFeedPosts))

	http.HandleFunc("/api/events/create", middleware.WithCORS(handlers.HandleCreateEvent))
	http.HandleFunc("/api/events/respond", middleware.WithCORS(handlers.HandleEventResponse))
	http.HandleFunc("/api/events/", middleware.WithCORS(handlers.GetEventByID))
	http.HandleFunc("/api/events/user", middleware.WithCORS(handlers.GetEventsByUserID))
	http.HandleFunc("/api/events/group/", middleware.WithCORS(handlers.GetEventsByGroupID))

	http.HandleFunc("/api/suggestgroups", middleware.WithCORS(handlers.HandleSuggestGroups))
	http.HandleFunc("/api/groups/search", middleware.WithCORS(handlers.SearchGroups))
	http.HandleFunc("/api/groups/user/", middleware.WithCORS(handlers.HandleGroupsByUserId))            // groups with user id
	http.HandleFunc("/api/groups/requested", middleware.WithCORS(handlers.HandleGroupRequests))         // active user group requests
	http.HandleFunc("/api/groups/invitations", middleware.WithCORS(handlers.HandleGroupInvitations))    // active user group invitations
	http.HandleFunc("/api/groups/administered", middleware.WithCORS(handlers.HandleGroupsAdministered)) // active user group invitations
	http.HandleFunc("/api/groups/create", middleware.WithCORS(handlers.HandleCreateGroup))
	http.HandleFunc("/api/group/", middleware.WithCORS(handlers.HandleGroupById)) // group with group id
	http.HandleFunc("/api/group/join", middleware.WithCORS(handlers.HandleGroupMembership))
	http.HandleFunc("/api/group/requests/{approval_status}", middleware.WithCORS(handlers.HandleGroupRequestApprove))
	http.HandleFunc("/api/group/invite", middleware.WithCORS(handlers.HandleGroupInvitation))
	http.HandleFunc("/api/group/invite/{id}/{approval_status}", middleware.WithCORS(handlers.HandleApproveGroupInvitation))
	http.HandleFunc("/api/group/invite/search", middleware.WithCORS(handlers.HandleGroupInvitationSearch))
	http.HandleFunc("/api/group/chat/messages/", middleware.WithCORS(handlers.HandleGetGroupMessagesByGroupId))
	http.HandleFunc("/api/group-posts/create", middleware.WithCORS(handlers.CreateGroupPostHandler))

	http.HandleFunc("/api/login", middleware.WithCORS(handlers.HandleLogin))

	http.HandleFunc("/api/auth/github", middleware.WithCORS(handlers.HandleGitHubLogin))
	http.HandleFunc("/api/auth/github/callback", middleware.WithCORS(handlers.HandleGitHubCallback))
	http.HandleFunc("/api/register", middleware.WithCORS(handlers.HandleRegister))
	http.HandleFunc("/api/logout", middleware.WithCORS(handlers.HandleLogout))
	http.HandleFunc("/api/me", middleware.WithCORS(handlers.HandleMe))
	http.HandleFunc("/api/me/update", middleware.WithCORS(handlers.HandleUpdateMe))
	http.HandleFunc("/api/following/", middleware.WithCORS(handlers.HandleFollowing))
	http.HandleFunc("/api/follow", middleware.WithCORS(handlers.HandleFollowAction))
	http.HandleFunc("/api/followers/", middleware.WithCORS(handlers.GetFollowers))
	http.HandleFunc("/api/followed/", middleware.WithCORS(handlers.GetFollowedUsers))
	http.HandleFunc("/api/follow/requests/sent", middleware.WithCORS(handlers.GetSentFollowRequests))
	http.HandleFunc("/api/follow/requests/received", middleware.WithCORS(handlers.GetReceivedFollowRequests))
	http.HandleFunc("/api/follow/requests/{id}/accept", middleware.WithCORS(handlers.HandleFollowRequestApprove))
	http.HandleFunc("/api/follow/requests/{id}/decline", middleware.WithCORS(handlers.HandleFollowRequestApprove))
	http.HandleFunc("/api/suggest/users", middleware.WithCORS(handlers.GetSuggestedUsers))

	http.HandleFunc("/api/notifications", middleware.WithCORS(handlers.HandleGetNotifications))
	http.HandleFunc("/api/notifications/{id}", middleware.WithCORS(handlers.GetNotificationByID))
	http.HandleFunc("/api/notifications/{id}/read", middleware.WithCORS(handlers.ReadNotification))
	http.HandleFunc("/api/notifications/new", middleware.WithCORS(handlers.GetNewNotifications))
	http.HandleFunc("/api/notifications/{id}/joingroup", middleware.WithCORS(handlers.HandleJoinReqsByGroupId))

	http.HandleFunc("/api/chat/messages", middleware.WithCORS(handlers.HandleGetUserMessages))

	http.HandleFunc("/ws", middleware.WithCORS(handlers.HandleWSConnections)) // Is CORS needed for websockets?
	//http.HandleFunc("/ws", handlers.HandleWSConnections)
	http.HandleFunc("/api/comments/show", middleware.WithCORS(handlers.HandleCommentsForPost))
	http.HandleFunc("/api/comments/create", middleware.WithCORS(handlers.HandleCreateCommentsForPost))

	// Serve the image directories as static content with CORS

	avatarsFS := http.FileServer(http.Dir("./data/uploads/avatars"))
	avatarHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware.WithCORS(avatarsFS.ServeHTTP)(w, r)
	})
	http.Handle("/data/uploads/avatars/", http.StripPrefix("/data/uploads/avatars/", avatarHandler))

	postsFS := http.FileServer(http.Dir("./data/uploads/posts"))
	postImageHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware.WithCORS(postsFS.ServeHTTP)(w, r)
	})
	http.Handle("/data/uploads/posts/", http.StripPrefix("/data/uploads/posts/", postImageHandler))

	commentsFS := http.FileServer(http.Dir("./data/uploads/comments"))
	commentImageHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware.WithCORS(commentsFS.ServeHTTP)(w, r)
	})
	http.Handle("/data/uploads/comments/", http.StripPrefix("/data/uploads/comments/", commentImageHandler))

	defaultFS := http.FileServer(http.Dir("./data/default"))
	defaultImageHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware.WithCORS(defaultFS.ServeHTTP)(w, r)
	})
	http.Handle("/data/default/", http.StripPrefix("/data/default/", defaultImageHandler))
}

func main() {

	config.InitConfig()
	fmt.Println(config.GithubOAuthConfig.ClientID)
	handlers.SetGitHubOAuthConfig(config.GithubOAuthConfig)
	err := database.NewDatabase(config.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	//deleteUnusedImages()	// delete database and run backend once before commenting this back in

	go service.StartBroadcastListener()

	setHandlers()
	fmt.Printf("Backend running on port %s, allowing requests from %s\n", config.Port, config.FrontendURL)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
