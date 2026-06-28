package handlers

import (
	"encoding/json"
	"net/http"
)

type healthResponsible struct {
	Status string `json:"status"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "http error:", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responsible := healthResponsible{Status: "OK"}
	json.NewEncoder(w).Encode(responsible)
}