package models

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
	Images    []string          `json:"images"`
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
	Images    []string          `json:"images"`
}

// StatefulSet représente un statefulset simplifié
type StatefulSet struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Replicas  int32             `json:"replicas"`
	Ready     int32             `json:"ready"`
	Labels    map[string]string `json:"labels,omitempty"`
	Age       string            `json:"age"`
	Images    []string          `json:"images"`
}
