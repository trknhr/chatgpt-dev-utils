package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"github.com/trknhr/chatgpt-dev-utils/internal/ui"
)

// WebSocket related variables
var (
	upgrader  = websocket.Upgrader{}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan string)
)

func main() {
	// Always start WebSocket server
	go startWebSocketServer()

	// Create model with WebSocket integration
	model := ui.InitialModel(broadcast, func() int {
		return len(clients)
	})

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func startWebSocketServer() {
	http.HandleFunc("/ws", handleConnections)

	// Simple HTTP health check
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	go handleMessages()

	go func() {
		err := http.ListenAndServe(":32123", nil)
		if err != nil {
			log.Printf("WebSocket server error: %v", err)
		}
	}()
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer ws.Close()
	clients[ws] = true

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			// log.Println("Extension disconnected")
			delete(clients, ws)
			break
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
