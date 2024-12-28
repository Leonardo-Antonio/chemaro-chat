package router

import (
	"net/http"

	"github.com/Leonardo-Antonio/chemaro/handler"
	"github.com/gorilla/mux"
)

func API(r *mux.Router) {
	r.HandleFunc("/api/v1/chat/open/{code}", handler.NewChatHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/chat/{code}/messages", handler.DeleteMessagesHandler).Methods(http.MethodDelete)
	r.HandleFunc("/api/v1/chat/{code}/messages", handler.GetMessagesHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/upload/files", handler.UploadFileHandler).Methods(http.MethodPost)
}
