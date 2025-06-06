module backend

go 1.20

require (
	github.com/golang-migrate/migrate/v4 v4.18.3
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/joho/godotenv v1.5.1
	github.com/mattn/go-sqlite3 v1.14.28
	golang.org/x/crypto v0.38.0
	golang.org/x/oauth2 v0.18.0
// Remove golang.org/x/oauth2 and golang.org/x/oauth2/github for now
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
