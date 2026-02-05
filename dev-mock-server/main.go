package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients          = make(map[*websocket.Conn]bool)
	clientsMu        sync.Mutex
	broadcast        = make(chan string)
	messages         = make(map[string]json.RawMessage)
	clientIcons      = make(map[*websocket.Conn]string)
	validatedClients = make(map[*websocket.Conn]bool)
	clientGuids      = make(map[*websocket.Conn]string)
	microgameResults = make(map[string]float64)
	availableIcons   = []string{"ava", "cubert", "potat", "azra", "eepy", "frog", "iconokeeb", "void", "zoe"}
	gameStarted      = false
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

	go func() {
		time.Sleep(100 * time.Millisecond)
		sendAvailableIcons(conn)
	}()

	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		delete(clientIcons, conn)
		clientsMu.Unlock()
		broadcastAvailableIcons()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", msg)
		
		var packet map[string]interface{}
		if err := json.Unmarshal(msg, &packet); err == nil {
			if packet["type"] == "icon-selection" {
				if icon, ok := packet["icon"].(string); ok {
					clientsMu.Lock()
					
					// Check if icon is already taken by another client
					iconTaken := false
					for client, selectedIcon := range clientIcons {
						if client != conn && selectedIcon == icon {
							iconTaken = true
							break
						}
					}
					
					if !iconTaken {
						clientIcons[conn] = icon
						clientsMu.Unlock()
						broadcastAvailableIcons()
					} else {
						clientsMu.Unlock()
						// Send updated available list to this client
						sendAvailableIcons(conn)
						log.Printf("Icon %s already taken, rejected selection", icon)
					}
				}
			} else if packet["type"] == "name-entry" {
				guid := uuid.New().String()
				response := map[string]interface{}{
					"type": "name-validation",
					"from": "server",
					"guid": guid,
				}
				data, _ := json.Marshal(response)
				conn.WriteMessage(websocket.TextMessage, data)
				log.Println("Sent name-validation")
				
				clientsMu.Lock()
				validatedClients[conn] = true
				clientGuids[conn] = guid
				validatedCount := len(validatedClients)
				clientsMu.Unlock()
				
				if validatedCount == 1 && !gameStarted {
					gameStarted = true
					go microgameLoop()
				}
			} else if packet["type"] == "microgame-result" {
				if id, ok := packet["id"].(string); ok {
					if percentage, ok := packet["percentage"].(float64); ok {
						clientsMu.Lock()
						microgameResults[id] = percentage
						resultCount := len(microgameResults)
						clientsMu.Unlock()
						
						log.Printf("Received result from %s: %.1f%%", id, percentage)
						
						if resultCount == 3 {
							determineWinner()
						}
					}
				}
			}
		}
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

func getAvailableIcons() []string {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	
	selectedIcons := make(map[string]bool)
	for _, icon := range clientIcons {
		selectedIcons[icon] = true
	}
	
	available := []string{}
	for _, icon := range availableIcons {
		if !selectedIcons[icon] {
			available = append(available, icon)
		}
	}
	return available
}

func sendAvailableIcons(conn *websocket.Conn) {
	available := getAvailableIcons()
	msg := map[string]interface{}{
		"type":  "icons-available",
		"from":  "server",
		"icons": available,
	}
	data, _ := json.Marshal(msg)
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println("Failed to send icons-available:", err)
	} else {
		log.Println("Sent icons-available to client")
	}
}

func broadcastAvailableIcons() {
	available := getAvailableIcons()
	msg := map[string]interface{}{
		"type":  "icons-available",
		"from":  "server",
		"icons": available,
	}
	data, _ := json.Marshal(msg)
	
	clientsMu.Lock()
	for conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("Write error:", err)
		}
	}
	clientsMu.Unlock()
	log.Printf("Broadcasted available icons: %v", available)
}

func startMicrogame() {
	data := messages["microgame-start"]
	
	clientsMu.Lock()
	for conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("Write error:", err)
		}
	}
	clientsMu.Unlock()
	log.Println("Started microgame")
}

func microgameLoop() {
	toggle := false
	for {
		if toggle {
			msg := map[string]interface{}{
				"type": "display-text",
				"from": "server",
				"text": "Test message",
			}
			data, _ := json.Marshal(msg)
			clientsMu.Lock()
			for conn := range clients {
				conn.WriteMessage(websocket.TextMessage, data)
			}
			clientsMu.Unlock()
			log.Println("Sent display-text")
		} else {
			startMicrogame()
		}
		toggle = !toggle
		time.Sleep(7 * time.Second)
	}
}

func determineWinner() {
	var winnerGuid string
	var highestPercentage float64
	
	clientsMu.Lock()
	for guid, percentage := range microgameResults {
		if percentage > highestPercentage {
			highestPercentage = percentage
			winnerGuid = guid
		}
	}
	clientsMu.Unlock()
	
	// Send microgame-end to all clients
	endMsg := map[string]interface{}{
		"type": "microgame-end",
		"from": "server",
	}
	endData, _ := json.Marshal(endMsg)
	
	clientsMu.Lock()
	for conn := range clients {
		conn.WriteMessage(websocket.TextMessage, endData)
	}
	clientsMu.Unlock()
	
	time.Sleep(500 * time.Millisecond)
	
	// Send winner message
	winMsg := map[string]interface{}{
		"type": "display-text",
		"from": "server",
		"text": "You won!",
	}
	winData, _ := json.Marshal(winMsg)
	
	clientsMu.Lock()
	for conn, guid := range clientGuids {
		if guid == winnerGuid {
			if err := conn.WriteMessage(websocket.TextMessage, winData); err != nil {
				log.Println("Write error:", err)
			}
		}
	}
	clientsMu.Unlock()
	log.Printf("Winner: %s with %.1f%%", winnerGuid, highestPercentage)
}

func main() {
	loadMessages()
	go broadcastMessages()
	go readConsole()

	http.HandleFunc("/game", handleGame)
	log.Println("Mock server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
