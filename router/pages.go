package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Pages(r *mux.Router) {
	assetsDir := http.Dir("./assets/")
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(assetsDir)))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	r.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/chat.html")
	})
}
