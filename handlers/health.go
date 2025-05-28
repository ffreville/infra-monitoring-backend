package handlers

import (
	"net/http"

	"github.com/ffreville/infra-monitoring-backend/models"
)

// HealthHandler gère les requêtes de santé
type HealthHandler struct{}

// NewHealthHandler crée un nouveau handler de santé
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck handler pour vérifier la santé de l'application
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := models.HealthResponse{
		Status:  "healthy",
		Service: "kubernetes-namespaces-api",
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// RootHandler handler pour la route racine
func (h *HealthHandler) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := models.RootResponse{
		Message: "API Kubernetes Resources",
		Version: "1.0.0",
		Endpoints: map[string]string{
			"namespaces":   "/api/v1/namespaces",
			"deployments":  "/api/v1/deployments",
			"cronjobs":     "/api/v1/cronjobs",
			"statefulsets": "/api/v1/statefulsets",
			"health":       "/health",
		},
		QueryParameters: map[string]string{
			"namespace": "Filtrer par namespace (optionnel pour deployments, cronjobs, statefulsets)",
		},
	}

	writeJSONResponse(w, http.StatusOK, response)
}
