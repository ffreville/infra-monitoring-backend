package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ffreville/infra-monitoring-backend/models"
)

// writeJSONResponse écrit une réponse JSON
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Erreur lors de l'encodage JSON: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "encoding_error", "Erreur lors de l'encodage de la réponse")
	}
}

// writeErrorResponse écrit une réponse d'erreur JSON
func writeErrorResponse(w http.ResponseWriter, statusCode int, errorType, message string) {
	w.WriteHeader(statusCode)
	errorResponse := models.ErrorResponse{
		Error:   errorType,
		Message: message,
	}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Erreur lors de l'encodage de la réponse d'erreur: %v", err)
	}
}
