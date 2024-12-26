package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Leonardo-Antonio/chemaro/db/memory"
	"github.com/Leonardo-Antonio/loadconfig"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	err := loadconfig.New(loadconfig.TYPE_FILE_DOT_ENV, loadconfig.ENV_DEFAULT).Load()
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()

	assetsDir := http.Dir("./assets/")
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(assetsDir)))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	r.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/chat.html")
	})

	r.HandleFunc("/api/v1/chat/open/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := mux.Vars(r)["code"]
		password := r.URL.Query().Get("psw")
		uuid := uuid.New().String()

		groupId := base64.StdEncoding.EncodeToString([]byte(code + "__###__" + password))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"message": "ok",
			"action":  "redirect",
			"data": map[string]any{
				"code":     code,
				"redirect": "/chat?code=" + groupId + "&uuid=" + uuid,
			},
		})
	}).Methods(http.MethodPost)

	r.HandleFunc("/api/v1/chat/{code}", NewWebSocketHandler().HandleWebSocket)
	r.HandleFunc("/api/v1/chat/{code}/messages", func(w http.ResponseWriter, r *http.Request) {
		code := mux.Vars(r)["code"]
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data := memory.Get(code)
		if data == nil {
			data = []memory.Message{}
		}

		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"message": "ok",
			"data": map[string]any{
				"messages": data,
			},
		})
	}).Methods(http.MethodGet)

	r.HandleFunc("/api/v1/chat/{code}/messages", func(w http.ResponseWriter, r *http.Request) {
		code := mux.Vars(r)["code"]
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		messages := memory.Get(code)
		for _, message := range messages {
			if strings.Contains(message.Type, "image") || strings.Contains(message.Type, "video") {
				os.Remove(message.Message)
			}
		}

		memory.Delete(code)
		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"message": "ok",
			"action":  "reload",
		})
	}).Methods(http.MethodDelete)

	r.HandleFunc("/api/v1/upload/files", func(w http.ResponseWriter, r *http.Request) {
		// obtener imagen que llega en formulario "file"
		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		name := handler.Filename
		contentType := handler.Header.Get("Content-Type")
		size := handler.Size
		ext := filepath.Ext(name)
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

		pathFiles := path.Join("assets", "temp")
		fullpath := path.Join(pathFiles, filename)
		tempFile, err := os.Create(fullpath)
		if err != nil {
			log.Println(err)
			return
		}
		defer tempFile.Close()

		_, err = io.Copy(tempFile, file)
		if err != nil {
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"message": "ok",
			"data": map[string]any{
				"url":  fullpath,
				"name": name,
				"type": contentType,
				"size": size,
			},
		})
	}).Methods(http.MethodPost)

	port := loadconfig.GetEnv("APP_PORT").String()
	log.Printf("Listening on port http://0.0.0.0:%s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

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
