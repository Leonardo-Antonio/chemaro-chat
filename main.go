package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Leonardo-Antonio/chemaro/router"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	router.Pages(r)
	router.API(r)
	router.APIReports(r)
	router.APIWs(r)

	port := os.Getenv("PORT")
	log.Printf("Listening on port http://0.0.0.0:%s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
