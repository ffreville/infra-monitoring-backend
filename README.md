# Kubernetes Resources API

Une API REST écrite en Go pour interagir avec les ressources d'un cluster Kubernetes. Cette application permet de lister facilement les namespaces, déploiements, cronjobs et statefulsets via des endpoints HTTP.

## 🚀 Fonctionnalités

- **Namespaces** : Liste tous les namespaces du cluster
- **Déploiements** : Récupère les déploiements avec leurs statuts
- **CronJobs** : Affiche les tâches planifiées et leur état
- **StatefulSets** : Liste les ensembles avec état et leurs répliques
- **Health Check** : Endpoint de vérification de santé
- **Filtrage** : Support du filtrage par namespace

## 📋 Prérequis

- Go 1.21 ou plus récent
- Accès à un cluster Kubernetes
- Kubeconfig configuré ou déploiement dans un cluster

## 🛠️ Installation

### Développement local

1. **Cloner le projet**
```bash
git clone <repository-url>
cd infra-monitoring-backend
```

2. **Initialiser le module Go**
```bash
go mod init infra-monitoring-backend
```

3. **Installer les dépendances**
```bash
go get github.com/gorilla/mux
go get k8s.io/client-go@latest
go mod tidy
```

4. **Lancer l'application**
```bash
go run main.go
```

### Avec Docker

1. **Construire l'image**
```bash
docker build -t infra-monitoring-backend:latest .
```

2. **Lancer le conteneur**
```bash
# Avec kubeconfig local
docker run -d \
  --name infra-monitoring-backend \
  -p 8080:8080 \
  -v ~/.kube/config:/home/appuser/.kube/config:ro \
  infra-monitoring-backend:latest
```

## 🔧 Configuration

### Variables d'environnement

| Variable | Description | Défaut |
|----------|-------------|---------|
| `PORT` | Port d'écoute du serveur | `8080` |
| `KUBECONFIG` | Chemin vers le fichier kubeconfig | `~/.kube/config` |

### Authentification Kubernetes

L'application supporte deux modes d'authentification :

1. **In-Cluster** : Utilise automatiquement les credentials du ServiceAccount quand déployée dans Kubernetes
2. **Externe** : Utilise le fichier kubeconfig local

## 📡 API Endpoints

### Informations générales
```http
GET /
```
Retourne les informations sur l'API et les endpoints disponibles.

### Health Check
```http
GET /health
```
Vérifie la santé de l'application.

### Namespaces
```http
GET /api/v1/namespaces
```
Liste tous les namespaces du cluster.

**Réponse exemple :**
```json
{
  "namespaces": [
    {
      "name": "default",
      "status": "Active",
      "labels": {
        "kubernetes.io/metadata.name": "default"
      },
      "age": "2024-01-15 10:30:45"
    }
  ],
  "count": 1
}
```

### Déploiements
```http
GET /api/v1/deployments[?namespace=<namespace>]
```
Liste les déploiements du cluster ou d'un namespace spécifique.

**Réponse exemple :**
```json
{
  "deployments": [
    {
      "name": "nginx-deployment",
      "namespace": "default",
      "replicas": 3,
      "ready": 3,
      "available": 3,
      "labels": {
        "app": "nginx"
      },
      "age": "2024-01-15 14:30:45"
    }
  ],
  "count": 1
}
```

### CronJobs
```http
GET /api/v1/cronjobs[?namespace=<namespace>]
```
Liste les cronjobs du cluster ou d'un namespace spécifique.

**Réponse exemple :**
```json
{
  "cronjobs": [
    {
      "name": "backup-job",
      "namespace": "default",
      "schedule": "0 2 * * *",
      "suspend": false,
      "active": 0,
      "lastRun": "2024-01-15 02:00:00",
      "labels": {
        "app": "backup"
      },
      "age": "2024-01-10 10:15:30"  
    }
  ],
  "count": 1
}
```

### StatefulSets
```http
GET /api/v1/statefulsets[?namespace=<namespace>]
```
Liste les statefulsets du cluster ou d'un namespace spécifique.

**Réponse exemple :**
```json
{
  "statefulsets": [
    {
      "name": "postgres-db",
      "namespace": "database",
      "replicas": 3,
      "ready": 3,
      "labels": {
        "app": "postgres"
      },
      "age": "2024-01-12 09:20:15"
    }
  ],
  "count": 1
}
```

## 💡 Exemples d'utilisation

### Curl

```bash
# Lister tous les namespaces
curl http://localhost:8080/api/v1/namespaces

# Lister tous les déploiements
curl http://localhost:8080/api/v1/deployments

# Lister les déploiements du namespace "production"
curl "http://localhost:8080/api/v1/deployments?namespace=production"

# Lister les cronjobs du namespace "batch"
curl "http://localhost:8080/api/v1/cronjobs?namespace=batch"

# Vérifier la santé de l'application
curl http://localhost:8080/health
```

### JavaScript/Fetch

```javascript
// Récupérer tous les déploiements
const response = await fetch('http://localhost:8080/api/v1/deployments');
const data = await response.json();
console.log(`${data.count} déploiements trouvés`);

// Filtrer par namespace
const prodDeployments = await fetch(
  'http://localhost:8080/api/v1/deployments?namespace=production'
);
```


### Métriques
L'application expose des informations de santé via l'endpoint `/health`.

## 📄 Licence

Ce projet est distribué sous licence MIT. Voir le fichier `LICENSE` pour plus de détails.

## 🤝 Contribution

Les contributions sont les bienvenues ! N'hésitez pas à :

1. Fork le projet
2. Créer une branche pour votre fonctionnalité
3. Commiter vos changements
4. Pousser vers la branche
5. Ouvrir une Pull Request

## 📞 Support

Pour toute question ou problème :
- Ouvrez une issue sur GitHub
- Consultez la documentation Kubernetes
- Vérifiez les logs de l'application
