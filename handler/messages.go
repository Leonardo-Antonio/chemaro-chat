package handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/Leonardo-Antonio/chemaro/db"
	"github.com/Leonardo-Antonio/chemaro/dto"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func NewChatHandler(w http.ResponseWriter, r *http.Request) {
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
}

func DeleteMessagesHandler(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	messages := db.DB.Get(code)
	for _, message := range messages {
		if strings.Contains(message.Type, "image") || strings.Contains(message.Type, "video") {
			os.Remove(message.Message)
		}
	}

	db.DB.Delete(code)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "ok",
		"action":  "reload",
	})
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data := db.DB.Get(code)
	if data == nil {
		data = []dto.Message{}
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "ok",
		"data": map[string]any{
			"messages": data,
		},
	})
}
