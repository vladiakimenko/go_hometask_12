package http

import (
	"encoding/json"
	"log"
	"net/http"
)

func sendJSONResponse(w http.ResponseWriter, data any, statusCode int) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}

func getUserIDFromContext(r *http.Request) (int, bool) {
	value := r.Context().Value(userIDKey)
	userID, ok := value.(int)
	if !ok {
		log.Printf("Failed to assert type int to %v", value)
		return 0, false
	}
	return userID, true
}
