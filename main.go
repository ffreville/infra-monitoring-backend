package main

import (
	"log"
	"os"

	"github.com/ffreville/infra-monitoring-backend/server"
)

func main() {
	// Récupérer le port depuis les variables d'environnement
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Créer et démarrer le serveur API
	apiServer, err := server.NewAPIServer()
	if err != nil {
		log.Fatalf("Erreur lors de la création du serveur API: %v", err)
	}

	log.Fatal(apiServer.Start(port))
}
