# Build stage
FROM golang:1.24-alpine AS builder

# Installer les certificats SSL et git
RUN apk add --no-cache ca-certificates git

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers go.mod et go.sum
COPY go.mod go.sum ./

# Télécharger les dépendances
RUN go mod download

# Copier le code source
COPY . .

# Compiler l'application
# CGO_ENABLED=0 pour un binaire statique
# GOOS=linux pour s'assurer que le binaire fonctionne sur Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o github.com/ffreville/infra-monitoring-backend

# Production stage
FROM alpine:latest

# Installer les certificats SSL pour les appels HTTPS
RUN apk --no-cache add ca-certificates

# Créer un utilisateur non-root pour la sécurité
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Définir le répertoire de travail
WORKDIR /app

# Copier le binaire depuis le stage de build
COPY --from=builder /app/github.com/ffreville/infra-monitoring-backend .

# Changer le propriétaire du fichier
RUN chown appuser:appgroup infra-monitoring-backend

# Utiliser l'utilisateur non-root
USER appuser

# Exposer le port 8080
EXPOSE 8080

# Définir les variables d'environnement par défaut
ENV PORT=8080

# Commande par défaut
CMD ["./infra-monitoring-backend"]
