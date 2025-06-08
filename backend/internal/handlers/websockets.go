package handlers

import (
	"backend/config"
	"backend/internal/model"
	"backend/internal/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == config.FrontendURL // explicitly allow frontend origin
	},
}

func HandleWSConnections(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &model.Client{
		UserID: fmt.Sprint(userID),
		Conn:   conn,
		Send:   make(chan model.WSMessage, 10), // buffered send
	}

	model.Mu.Lock()
	model.Clients[client.UserID] = client
	model.Mu.Unlock()

	go service.ReadPump(client)
	go service.WritePump(client)

	// send pong to test connection
	msg := model.WSMessage{
		Type:    "pong",
		From:    "system_pong",
		To:      client.UserID,
		Content: "ensuring connection works",
	}

	err = client.Conn.WriteJSON(msg)
	if err != nil {
		fmt.Println("error at pong:", err)
	}
}
