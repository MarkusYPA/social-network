package config

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv" // Import the godotenv library
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/github"
)

const DBPath = "data/db/app.db"
const MigrationsPath = "migrations"

var (
    Port              string
    FrontendURL       string
    GithubOAuthConfig *oauth2.Config
)

func InitConfig() {
    // Load environment variables from .env file FIRST
    // It's good practice to check for an error, but don't fatal if the file simply doesn't exist (e.g., in production)
    err := godotenv.Load()
    if err != nil {
        // Log a warning if .env is not found, but allow
        // the program to continue if environment variables are set another way.
        // If you *require* a .env file, you could log.Fatal here.
        log.Printf("WARNING: Error loading .env file (it might not exist or be unreadable): %v. Relying on system environment variables.\n", err)
    }

    Port = os.Getenv("PORT")
    if Port == "" {
        fmt.Println("INFO: PORT environment variable not set, falling back to default port 8080")
        Port = "8080" // Default port
    }

    FrontendURL = os.Getenv("FRONTEND_URL")
    if FrontendURL == "" {
        fmt.Println("INFO: FRONTEND_URL environment variable not set, falling back to default frontend URL http://localhost:5173")
        FrontendURL = "http://localhost:5173" // Default frontend URL
    }

    GithubOAuthConfig = &oauth2.Config{
        ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
        ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
        RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
        Scopes:       []string{"user:email"},
        Endpoint:     github.Endpoint,
    }

    // Validate critical OAuth configuration
    // These checks will now occur AFTER the .env file has been loaded
    if GithubOAuthConfig.ClientID == "" {
        log.Fatal("ERROR: GITHUB_CLIENT_ID environment variable is not set. GitHub OAuth cannot proceed.")
    }
    if GithubOAuthConfig.ClientSecret == "" {
        log.Fatal("ERROR: GITHUB_CLIENT_SECRET environment variable is not set. GitHub OAuth cannot proceed.")
    }
    if GithubOAuthConfig.RedirectURL == "" {
        log.Fatal("ERROR: GITHUB_REDIRECT_URL environment variable is not set. GitHub OAuth cannot proceed.")
    }
}