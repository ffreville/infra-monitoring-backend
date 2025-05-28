package models

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

// HealthResponse représente la réponse du health check
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// RootResponse représente la réponse de la route racine
type RootResponse struct {
	Message         string            `json:"message"`
	Version         string            `json:"version"`
	Endpoints       map[string]string `json:"endpoints"`
	QueryParameters map[string]string `json:"query_parameters"`
}
