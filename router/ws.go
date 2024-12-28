package router

import (
	"github.com/Leonardo-Antonio/chemaro/handler"
	"github.com/gorilla/mux"
)

func APIWs(r *mux.Router) {
	r.HandleFunc("/api/v1/chat/{code}", handler.NewWebSocketHandler().HandleWebSocket)
}
