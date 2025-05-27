package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Namespace représente un namespace simplifié
type Namespace struct {
	Name   string            `json:"name"`
	Status string            `json:"status"`
	Labels map[string]string `json:"labels,omitempty"`
	Age    string            `json:"age"`
}

// Deployment représente un déploiement simplifié
type Deployment struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Replicas  int32             `json:"replicas"`
	Ready     int32             `json:"ready"`
	Available int32             `json:"available"`
	Labels    map[string]string `json:"labels,omitempty"`
	Age       string            `json:"age"`
}

// CronJob représente un cronjob simplifié
type CronJob struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Schedule  string            `json:"schedule"`
	Suspend   bool              `json:"suspend"`
	Active    int               `json:"active"`
	LastRun   string            `json:"lastRun,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
	Age       string            `json:"age"`
}

// StatefulSet représente un statefulset simplifié
type StatefulSet struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Replicas  int32             `json:"replicas"`
	Ready     int32             `json:"ready"`
	Labels    map[string]string `json:"labels,omitempty"`
	Age       string            `json:"age"`
}

// NamespaceResponse représente la réponse de l'API pour les namespaces
type NamespaceResponse struct {
	Namespaces []Namespace `json:"namespaces"`
	Count      int         `json:"count"`
}

// DeploymentResponse représente la réponse de l'API pour les déploiements
type DeploymentResponse struct {
	Deployments []Deployment `json:"deployments"`
	Count       int          `json:"count"`
}

// CronJobResponse représente la réponse de l'API pour les cronjobs
type CronJobResponse struct {
	CronJobs []CronJob `json:"cronjobs"`
	Count    int       `json:"count"`
}

// StatefulSetResponse représente la réponse de l'API pour les statefulsets
type StatefulSetResponse struct {
	StatefulSets []StatefulSet `json:"statefulsets"`
	Count        int           `json:"count"`
}

// ErrorResponse représente une réponse d'erreur
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// KubernetesClient encapsule le client Kubernetes
type KubernetesClient struct {
	clientset *kubernetes.Clientset
}

// NewKubernetesClient crée un nouveau client Kubernetes
func NewKubernetesClient() (*KubernetesClient, error) {
	var config *rest.Config
	var err error

	// Essayer d'utiliser la configuration in-cluster
	config, err = rest.InClusterConfig()
	if err != nil {
		// Si on n'est pas dans un cluster, utiliser kubeconfig
		var kubeconfig string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}

		// Vérifier si KUBECONFIG est défini
		if kubeconfigEnv := os.Getenv("KUBECONFIG"); kubeconfigEnv != "" {
			kubeconfig = kubeconfigEnv
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la création de la config Kubernetes: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du client Kubernetes: %v", err)
	}

	return &KubernetesClient{
		clientset: clientset,
	}, nil
}

// GetNamespaces récupère tous les namespaces du cluster
func (k *KubernetesClient) GetNamespaces(ctx context.Context) ([]Namespace, error) {
	namespaceList, err := k.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des namespaces: %v", err)
	}

	var namespaces []Namespace
	for _, ns := range namespaceList.Items {
		namespace := Namespace{
			Name:   ns.Name,
			Status: string(ns.Status.Phase),
			Labels: ns.Labels,
			Age:    ns.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		namespaces = append(namespaces, namespace)
	}

	return namespaces, nil
}

// GetDeployments récupère tous les déploiements du cluster ou d'un namespace spécifique
func (k *KubernetesClient) GetDeployments(ctx context.Context, namespace string) ([]Deployment, error) {
	deploymentList, err := k.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des déploiements: %v", err)
	}

	var deployments []Deployment
	for _, deploy := range deploymentList.Items {
		deployment := Deployment{
			Name:      deploy.Name,
			Namespace: deploy.Namespace,
			Replicas:  *deploy.Spec.Replicas,
			Ready:     deploy.Status.ReadyReplicas,
			Available: deploy.Status.AvailableReplicas,
			Labels:    deploy.Labels,
			Age:       deploy.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		deployments = append(deployments, deployment)
	}

	return deployments, nil
}

// GetCronJobs récupère tous les cronjobs du cluster ou d'un namespace spécifique
func (k *KubernetesClient) GetCronJobs(ctx context.Context, namespace string) ([]CronJob, error) {
	cronJobList, err := k.clientset.BatchV1().CronJobs(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des cronjobs: %v", err)
	}

	var cronJobs []CronJob
	for _, cj := range cronJobList.Items {
		var lastRun string
		if cj.Status.LastScheduleTime != nil {
			lastRun = cj.Status.LastScheduleTime.Format("2006-01-02 15:04:05")
		}

		cronJob := CronJob{
			Name:      cj.Name,
			Namespace: cj.Namespace,
			Schedule:  cj.Spec.Schedule,
			Suspend:   cj.Spec.Suspend != nil && *cj.Spec.Suspend,
			Active:    len(cj.Status.Active),
			LastRun:   lastRun,
			Labels:    cj.Labels,
			Age:       cj.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		cronJobs = append(cronJobs, cronJob)
	}

	return cronJobs, nil
}

// GetStatefulSets récupère tous les statefulsets du cluster ou d'un namespace spécifique
func (k *KubernetesClient) GetStatefulSets(ctx context.Context, namespace string) ([]StatefulSet, error) {
	statefulSetList, err := k.clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des statefulsets: %v", err)
	}

	var statefulSets []StatefulSet
	for _, sts := range statefulSetList.Items {
		statefulSet := StatefulSet{
			Name:      sts.Name,
			Namespace: sts.Namespace,
			Replicas:  *sts.Spec.Replicas,
			Ready:     sts.Status.ReadyReplicas,
			Labels:    sts.Labels,
			Age:       sts.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		statefulSets = append(statefulSets, statefulSet)
	}

	return statefulSets, nil
}

// APIServer représente le serveur API
type APIServer struct {
	k8sClient *KubernetesClient
	router    *mux.Router
}

// NewAPIServer crée un nouveau serveur API
func NewAPIServer() (*APIServer, error) {
	k8sClient, err := NewKubernetesClient()
	if err != nil {
		return nil, err
	}

	server := &APIServer{
		k8sClient: k8sClient,
		router:    mux.NewRouter(),
	}

	server.setupRoutes()
	return server, nil
}

// setupRoutes configure les routes de l'API
func (s *APIServer) setupRoutes() {
	s.router.HandleFunc("/api/v1/namespaces", s.getNamespaces).Methods("GET")
	s.router.HandleFunc("/api/v1/deployments", s.getDeployments).Methods("GET")
	s.router.HandleFunc("/api/v1/cronjobs", s.getCronJobs).Methods("GET")
	s.router.HandleFunc("/api/v1/statefulsets", s.getStatefulSets).Methods("GET")
	s.router.HandleFunc("/health", s.healthCheck).Methods("GET")
	s.router.HandleFunc("/", s.rootHandler).Methods("GET")
}

// getNamespaces handler pour récupérer les namespaces
func (s *APIServer) getNamespaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	namespaces, err := s.k8sClient.GetNamespaces(ctx)
	if err != nil {
		log.Printf("Erreur lors de la récupération des namespaces: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "internal_server_error",
			Message: "Impossible de récupérer les namespaces",
		})
		return
	}

	response := NamespaceResponse{
		Namespaces: namespaces,
		Count:      len(namespaces),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erreur lors de l'encodage JSON: %v", err)
	}
}

// getDeployments handler pour récupérer les déploiements
func (s *APIServer) getDeployments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Récupérer le paramètre namespace (optionnel)
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "" // Tous les namespaces
	}

	ctx := r.Context()
	deployments, err := s.k8sClient.GetDeployments(ctx, namespace)
	if err != nil {
		log.Printf("Erreur lors de la récupération des déploiements: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "internal_server_error",
			Message: "Impossible de récupérer les déploiements",
		})
		return
	}

	response := DeploymentResponse{
		Deployments: deployments,
		Count:       len(deployments),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erreur lors de l'encodage JSON: %v", err)
	}
}

// getCronJobs handler pour récupérer les cronjobs
func (s *APIServer) getCronJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Récupérer le paramètre namespace (optionnel)
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "" // Tous les namespaces
	}

	ctx := r.Context()
	cronJobs, err := s.k8sClient.GetCronJobs(ctx, namespace)
	if err != nil {
		log.Printf("Erreur lors de la récupération des cronjobs: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "internal_server_error",
			Message: "Impossible de récupérer les cronjobs",
		})
		return
	}

	response := CronJobResponse{
		CronJobs: cronJobs,
		Count:    len(cronJobs),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erreur lors de l'encodage JSON: %v", err)
	}
}

// getStatefulSets handler pour récupérer les statefulsets
func (s *APIServer) getStatefulSets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Récupérer le paramètre namespace (optionnel)
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "" // Tous les namespaces
	}

	ctx := r.Context()
	statefulSets, err := s.k8sClient.GetStatefulSets(ctx, namespace)
	if err != nil {
		log.Printf("Erreur lors de la récupération des statefulsets: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "internal_server_error",
			Message: "Impossible de récupérer les statefulsets",
		})
		return
	}

	response := StatefulSetResponse{
		StatefulSets: statefulSets,
		Count:        len(statefulSets),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erreur lors de l'encodage JSON: %v", err)
	}
}

// healthCheck handler pour vérifier la santé de l'application
func (s *APIServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status":  "healthy",
		"service": "kubernetes-namespaces-api",
	}

	json.NewEncoder(w).Encode(response)
}

// rootHandler handler pour la route racine
func (s *APIServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message": "API Kubernetes Resources",
		"version": "1.0.0",
		"endpoints": map[string]string{
			"namespaces":   "/api/v1/namespaces",
			"deployments":  "/api/v1/deployments",
			"cronjobs":     "/api/v1/cronjobs",
			"statefulsets": "/api/v1/statefulsets",
			"health":       "/health",
		},
		"query_parameters": map[string]string{
			"namespace": "Filtrer par namespace (optionnel pour deployments, cronjobs, statefulsets)",
		},
	}

	json.NewEncoder(w).Encode(response)
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

func main() {
	// Récupérer le port depuis les variables d'environnement
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Créer et démarrer le serveur API
	server, err := NewAPIServer()
	if err != nil {
		log.Fatalf("Erreur lors de la création du serveur API: %v", err)
	}

	log.Fatal(server.Start(port))
}
