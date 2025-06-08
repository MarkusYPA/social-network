package handlers

import (
	"backend/config"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/internal/utils"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log" // Import log for error handling
	"net/http"

	"golang.org/x/oauth2"
	_ "golang.org/x/oauth2/github"
)

var githubOAuthConfig *oauth2.Config

// SetGitHubOAuthConfig is a function to set the config from main
func SetGitHubOAuthConfig(config *oauth2.Config) {
	githubOAuthConfig = config
}

func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	// Generate a random state string (should be stored in session for validation)
	oauthState := generateStateOauthCookie(w)

	// Ensure the 'user:email' scope is included to get access to the /user/emails endpoint
	// You already have this in your config, but it's good to confirm.
	url := githubOAuthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state matches
	fmt.Println("DEBUG: Handling GitHub callback with state:", r.FormValue("state"))
	oauthStateCookie, err := r.Cookie("oauthstate")
	if err != nil || r.FormValue("state") != oauthStateCookie.Value {
		log.Printf("ERROR: Invalid OAuth state. Cookie: %v, FormValue: %v", oauthStateCookie, r.FormValue("state"))
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	token, err := githubOAuthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Printf("ERROR: Failed to exchange token: %v", err)
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user info from GitHub (basic profile)
	client := githubOAuthConfig.Client(context.Background(), token)
	respUser, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Printf("ERROR: Failed to get user info from /user endpoint: %v", err)
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer respUser.Body.Close()

	var githubUser struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
		Email string `json:"email"` // This will be the public email or empty
	}

	if err := json.NewDecoder(respUser.Body).Decode(&githubUser); err != nil {
		log.Printf("ERROR: Failed to decode user info from /user endpoint: %v", err)
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// --- NEW PART: Fetching all emails ---
	respEmails, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		log.Printf("ERROR: Failed to get user emails from /user/emails endpoint: %v", err)
		http.Error(w, "Failed to get user emails: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer respEmails.Body.Close()

	var githubEmails []struct {
		Email      string `json:"email"`
		Primary    bool   `json:"primary"`
		Verified   bool   `json:"verified"`
		Visibility string `json:"visibility"`
	}

	if err := json.NewDecoder(respEmails.Body).Decode(&githubEmails); err != nil {
		log.Printf("ERROR: Failed to decode user emails from /user/emails endpoint: %v", err)
		http.Error(w, "Failed to decode user emails: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Now you have all the emails in `githubEmails`.
	// You can iterate through them to find the primary, verified, or any specific email you need.
	var primaryEmail string
	for _, email := range githubEmails {
		if email.Primary && email.Verified {
			primaryEmail = email.Email
			break
		}
	}

	// You can now use `githubUser.Login` (username), `githubUser.Name`, and `primaryEmail`
	// (or any other email from `githubEmails`) for your user management.

	// Combine the data for the response (or use them for DB operations)
	response := struct {
		ID           int      `json:"id"`
		Login        string   `json:"login"`
		Name         string   `json:"name"`
		PublicEmail  string   `json:"public_email"`  // Email from /user endpoint
		PrimaryEmail string   `json:"primary_email"` // Primary verified email from /user/emails
		AllEmails    []string `json:"all_emails"`    // All emails
	}{
		ID:           githubUser.ID,
		Login:        githubUser.Login,
		Name:         githubUser.Name,
		PublicEmail:  githubUser.Email,
		PrimaryEmail: primaryEmail,
		AllEmails:    make([]string, len(githubEmails)),
	}

	for i, email := range githubEmails {
		response.AllEmails[i] = email.Email
	}
	var req model.LoginRequest
	req.Email = primaryEmail // Use primary email for user lookup
	req.Password = ""        // Password is not used for OAuth login
	user, found, err := repository.GetUserByEmail(req)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("ERROR: Database error while checking user by email: %v", err)
		http.Error(w, "Internal server error during user lookup", http.StatusInternalServerError)
		return
	}

	if !found {
		// User does not exist, create a new one
		log.Printf("INFO: User with email %s not found. Creating new user.", primaryEmail)

		// Generate random password (you should hash this in a real app!)
		randomPassword := "pass"
		//randomPassword := utils.GenerateRandomString(16) // Define this helper
		randomDOB := utils.GenerateRandomDateOfBirth() // Define this helper

		// newUser := model.User{
		// 	Username:  githubUser.Login, // GitHub login as username
		// 	Email:     primaryEmail,     // Primary verified GitHub email
		// 	Password:  randomPassword,   // A temporary, random password (should be hashed)
		// 	Birthday:  randomDOB,        // Random Date of Birth
		// 	FirstName: githubUser.Login,
		// 	LastName:  "", // GitHub doesn't provide last name, so you might want to set it to empty or a default value
		// 	// Set other fields as necessary, e.g., created_at, updated_at
		// }
		userInfo := struct {
			Email      string
			Password   string
			FirstName  string
			LastName   string
			DOB        string
			Nickname   string
			About      string
			AvatarPath sql.NullString
		}{
			Email:      primaryEmail,
			Password:   randomPassword,
			FirstName:  githubUser.Login, // Use GitHub name as first name
			LastName:   "",
			DOB:        randomDOB,
			Nickname:   githubUser.Login,
			About:      "",
			AvatarPath: sql.NullString{String: DefaultAvatarPath, Valid: true}, // No avatar path initially
		}

		errMsg, statusCode := service.RegisterUser(userInfo)
		if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
			http.Error(w, errMsg, statusCode)
			return
		}
		user, _, _ = repository.GetUserByEmail(req)

		log.Printf("INFO: New user created successfully: %s (%s)", user.Username, user.Email)

	} else {
		log.Printf("INFO: User with email %s found. Logging in.", primaryEmail)
		// User exists, proceed with login (no need to update their basic profile on every login unless required)
	}

	// Here you would:
	// 1. Check if user exists in your DB by GitHub ID (githubUser.ID)
	// 2. If not, create a new user account, likely using `primaryEmail` as their main contact.
	// 3. Generate a session/JWT token for the user
	// 4. Redirect back to frontend with the token
	err = service.CreateSession(user, w)
	if err != nil {
		log.Printf("ERROR: Failed to create session for user %s: %v", user.Username, err)
		http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}
	redirectURL := fmt.Sprintf("%s/", config.FrontendURL) // Add a query param for frontend to detect
	http.Redirect(w, r, redirectURL, http.StatusFound)
	// For now, just return the combined user info as JSON
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(user)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 32) // Generate a 32-byte random string
	if _, err := rand.Read(b); err != nil {
		log.Println("ERROR: Failed to generate random state:", err)
		// Handle error appropriately, maybe return a generic error to the user
		// or a default state if you have a robust fallback
		return "fallback-random-state" // Fallback for demonstration; in real app, might want to error out
	}
	state := base64.URLEncoding.EncodeToString(b)

	cookie := http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,                 // Set to false in development if using HTTP
		SameSite: http.SameSiteLaxMode, // Recommended for CSRF protection
	}
	http.SetCookie(w, &cookie)
	return state
}
