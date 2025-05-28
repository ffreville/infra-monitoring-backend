package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ffreville/infra-monitoring-backend/client"
	"github.com/ffreville/infra-monitoring-backend/models"
)

// KubernetesHandler gère les requêtes liées à Kubernetes
type KubernetesHandler struct {
	k8sClient *client.KubernetesClient
}

// NewKubernetesHandler crée un nouveau handler Kubernetes
func NewKubernetesHandler(k8sClient *client.KubernetesClient) *KubernetesHandler {
	return &KubernetesHandler{
		k8sClient: k8sClient,
	}
}

// GetNamespaces handler pour récupérer les namespaces
func (h *KubernetesHandler) GetNamespaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	namespaces, err := h.k8sClient.GetNamespaces(ctx)
	if err != nil {
		log.Printf("Erreur lors de la récupération des namespaces: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "internal_server_error", "Impossible de récupérer les namespaces")
		return
	}

	response := models.NamespaceResponse{
		Namespaces: namespaces,
		Count:      len(namespaces),
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// GetDeployments handler pour récupérer les déploiements
func (h *KubernetesHandler) GetDeployments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Récupérer le paramètre namespace (optionnel)
	namespace := r.URL.Query().Get("namespace")

	ctx := r.Context()
	deployments, err := h.k8sClient.GetDeployments(ctx, namespace)
	if err != nil {
		log.Printf("Erreur lors de la récupération des déploiements: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "internal_server_error", "Impossible de récupérer les déploiements")
		return
	}

	response := models.DeploymentResponse{
		Deployments: deployments,
		Count:       len(deployments),
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// GetCronJobs handler pour récupérer les cronjobs
func (h *KubernetesHandler) GetCronJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Récupérer le paramètre namespace (optionnel)
	namespace := r.URL.Query().Get("namespace")

	ctx := r.Context()
	cronJobs, err := h.k8sClient.GetCronJobs(ctx, namespace)
	if err != nil {
		log.Printf("Erreur lors de la récupération des cronjobs: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "internal_server_error", "Impossible de récupérer les cronjobs")
		return
	}

	response := models.CronJobResponse{
		CronJobs: cronJobs,
		Count:    len(cronJobs),
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// GetStatefulSets handler pour récupérer les statefulsets
func (h *KubernetesHandler) GetStatefulSets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Récupérer le paramètre namespace (optionnel)
	namespace := r.URL.Query().Get("namespace")

	ctx := r.Context()
	statefulSets, err := h.k8sClient.GetStatefulSets(ctx, namespace)
	if err != nil {
		log.Printf("Erreur lors de la récupération des statefulsets: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "internal_server_error", "Impossible de récupérer les statefulsets")
		return
	}

	response := models.StatefulSetResponse{
		StatefulSets: statefulSets,
		Count:        len(statefulSets),
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// writeErrorResponse écrit une réponse d'erreur JSON
func (h *KubernetesHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, errorType, message string) {
	w.WriteHeader(statusCode)
	errorResponse := models.ErrorResponse{
		Error:   errorType,
		Message: message,
	}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Erreur lors de l'encodage de la réponse d'erreur: %v", err)
	}
}

// writeJSONResponse écrit une réponse JSON
func (h *KubernetesHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Erreur lors de l'encodage JSON: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "encoding_error", "Erreur lors de l'encodage de la réponse")
	}
}
