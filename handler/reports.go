package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Leonardo-Antonio/chemaro/db/memory"
)

func ReportGroupsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "ok",
		"data": map[string]any{
			"groups": memory.GetAll(),
		},
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
