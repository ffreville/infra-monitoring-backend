package server

import (
	"log"
	"net/http"

	"github.com/ffreville/infra-monitoring-backend/client"
	"github.com/ffreville/infra-monitoring-backend/handlers"
	"github.com/gorilla/mux"
)

// APIServer représente le serveur API
type APIServer struct {
	k8sHandler    *handlers.KubernetesHandler
	healthHandler *handlers.HealthHandler
	router        *mux.Router
}

// NewAPIServer crée un nouveau serveur API
func NewAPIServer() (*APIServer, error) {
	k8sClient, err := client.NewKubernetesClient()
	if err != nil {
		return nil, err
	}

	server := &APIServer{
		k8sHandler:    handlers.NewKubernetesHandler(k8sClient),
		healthHandler: handlers.NewHealthHandler(),
		router:        mux.NewRouter(),
	}

	server.setupRoutes()
	return server, nil
}

// setupRoutes configure les routes de l'API
func (s *APIServer) setupRoutes() {
	// Routes API v1
	apiV1 := s.router.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/namespaces", s.k8sHandler.GetNamespaces).Methods("GET")
	apiV1.HandleFunc("/deployments", s.k8sHandler.GetDeployments).Methods("GET")
	apiV1.HandleFunc("/cronjobs", s.k8sHandler.GetCronJobs).Methods("GET")
	apiV1.HandleFunc("/statefulsets", s.k8sHandler.GetStatefulSets).Methods("GET")

	// Routes de santé et informations
	s.router.HandleFunc("/health", s.healthHandler.HealthCheck).Methods("GET")
	s.router.HandleFunc("/", s.healthHandler.RootHandler).Methods("GET")
}

// Start démarre le serveur API
func (s *APIServer) Start(port string) error {
	log.Printf("Démarrage du serveur API sur le port %s", port)
	log.Printf("Endpoints disponibles:")
	log.Printf("  GET /api/v1/namespaces - Lister les namespaces")
	log.Printf("  GET /api/v1/deployments - Lister les déploiements")
	log.Printf("  GET /api/v1/cronjobs - Lister les cronjobs")
	log.Printf("  GET /api/v1/statefulsets - Lister les statefulsets")
	log.Printf("  GET /health - Vérification de santé")
	log.Printf("  GET / - Informations sur l'API")
	log.Printf("Paramètres de requête:")
	log.Printf("  ?namespace=<nom> - Filtrer par namespace (pour deployments, cronjobs, statefulsets)")

	return http.ListenAndServe(":"+port, s.router)
}
