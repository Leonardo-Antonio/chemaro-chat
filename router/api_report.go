package router

import (
	"net/http"

	"github.com/Leonardo-Antonio/chemaro/handler"
	"github.com/gorilla/mux"
)

func APIReports(r *mux.Router) {
	r.HandleFunc("/api/v1/reports/groups", handler.ReportGroupsHandler).Methods(http.MethodGet)
}
