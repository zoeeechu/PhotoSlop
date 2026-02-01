package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	broadcast = make(chan string)
	messages  = make(map[string]json.RawMessage)
)

func handleGame(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	log.Println("Client connected")

	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", msg)
	}
}

func broadcastMessages() {
	for msg := range broadcast {
		clientsMu.Lock()
		for conn := range clients {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				log.Println("Write error:", err)
				conn.Close()
				delete(clients, conn)
			}
		}
		clientsMu.Unlock()
	}
}

func readConsole() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msgType := strings.TrimSpace(scanner.Text())
		if msg, ok := messages[msgType]; ok {
			broadcast <- string(msg)
			log.Printf("Sending %s", msgType)
		} else {
			log.Printf("Unknown message type: %s", msgType)
		}
	}
}

func loadMessages() {
	data, err := os.ReadFile("test-protocol.json")
	if err != nil {
		log.Fatal("Failed to read test-protocol.json:", err)
	}

	decoder := json.NewDecoder(strings.NewReader(string(data)))
	for decoder.More() {
		var msg map[string]interface{}
		if err := decoder.Decode(&msg); err != nil {
			continue
		}
		if from, ok := msg["from"].(string); ok && from == "server" {
			if msgType, ok := msg["type"].(string); ok {
				raw, _ := json.Marshal(msg)
				messages[msgType] = raw
			}
		}
	}
	log.Printf("Loaded %d server messages", len(messages))
}

func main() {
	loadMessages()
	go broadcastMessages()
	go readConsole()

	http.HandleFunc("/game", handleGame)
	log.Println("Mock server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
