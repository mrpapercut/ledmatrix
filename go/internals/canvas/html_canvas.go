package canvas

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type HTMLCanvas struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mutex     sync.Mutex
}

func (h *HTMLCanvas) init() {
	h.clients = make(map[*websocket.Conn]bool)
	h.broadcast = make(chan []byte)
	go h.startServer()
}

func (h *HTMLCanvas) startServer() {
	http.HandleFunc("/ws", h.handleConnections)
	go h.handleMessages()

	err := http.ListenAndServe("localhost:3000", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func (h *HTMLCanvas) handleConnections(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading GET request: %v", err)
		return
	}
	defer conn.Close()

	h.mutex.Lock()
	h.clients[conn] = true
	h.mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			h.mutex.Lock()
			delete(h.clients, conn)
			h.mutex.Unlock()
			break
		}

		h.broadcast <- message
	}
}

func (h *HTMLCanvas) handleMessages() {
	for {
		message := <-h.broadcast

		h.mutex.Lock()
		for client := range h.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mutex.Unlock()
	}
}

func (h *HTMLCanvas) sendMessage(command string, args ...[]string) {
	h.mutex.Lock()
	for client := range h.clients {
		parsedArgs, _ := json.Marshal(args)
		message := fmt.Sprintf("{\"command\": \"%s\", \"args\": %s}", command, string(parsedArgs))

		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Error writing message: %v", err)
			client.Close()
			delete(h.clients, client)
		}
	}
	h.mutex.Unlock()
}

func (h *HTMLCanvas) Clear() error {
	h.sendMessage("clear_screen")
	return nil
}

func (h *HTMLCanvas) Close() error {
	h.sendMessage("shutdown")
	return nil
}

func (h *HTMLCanvas) DrawScreen(pixeldata [][]int, colors []int, offsetX int, offsetY int) error {
	var pixelstrings []string

	for y := 0; y < len(pixeldata); y++ {
		for x := 0; x < len(pixeldata[y]); x++ {
			colorIndex := pixeldata[y][x]
			if colorIndex == -1 {
				continue
			}

			color := colors[colorIndex]
			if color == 0 {
				continue
			}

			pixelstrings = append(pixelstrings, fmt.Sprintf("%d:%d:%d", x+offsetX, y+offsetY, color))
		}
	}

	h.sendMessage("draw_screen", pixelstrings)
	return nil
}
