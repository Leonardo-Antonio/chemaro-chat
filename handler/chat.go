package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Leonardo-Antonio/chemaro/db/memory"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	connections map[string]map[*websocket.Conn]bool
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		connections: make(map[string]map[*websocket.Conn]bool),
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade la conexión a WebSocket
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Obtiene el ID del grupo al que se conecta el cliente
	groupID := mux.Vars(r)["code"]

	// Crea un mapa de conexiones para el grupo si no existe
	if _, ok := h.connections[groupID]; !ok {
		h.connections[groupID] = make(map[*websocket.Conn]bool)
	}

	// Agrega la conexión al mapa de conexiones del grupo
	h.connections[groupID][conn] = true

	// Lee los mensajes del cliente
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// Deserializa el mensaje JSON
		var msg struct {
			UserId  string `json:"userId"`
			Message string `json:"message"`
			Type    string `json:"type"`
		}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println(err)
			continue
		}

		memory.Set(groupID, memory.Message{
			Id:        uuid.New().String(),
			UserId:    msg.UserId,
			Message:   msg.Message,
			Type:      msg.Type,
			CreatedAt: uint64(time.Now().UTC().UnixMilli()),
		})

		h.forwardMessages(groupID)
	}
}

func (h *WebSocketHandler) forwardMessages(groupID string) {
	msgs := memory.Get(groupID)
	if msgs == nil {
		msgs = []memory.Message{}
	}
	buff, err := json.Marshal(msgs)
	if err != nil {
		log.Println(err)
	}

	h.broadcastMessage(groupID, buff)
}

func (h *WebSocketHandler) broadcastMessage(groupID string, message []byte) {
	// Obtiene la lista de conexiones del grupo
	connections := make(map[*websocket.Conn]bool)
	for conn := range h.connections[groupID] {
		connections[conn] = true
	}

	// Envía el mensaje a cada conexión
	for conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			// Elimina la conexión del mapa si falla
			delete(h.connections[groupID], conn)
		}
	}
}
