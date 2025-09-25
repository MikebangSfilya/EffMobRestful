package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	datatransfer "subscription/internal/api/dto"
)

func writeJSON(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode subscription: %v", err)
		datatransfer.WriteError(w, "failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}

func readJSON(r *http.Request, data any) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}
	return nil
}
