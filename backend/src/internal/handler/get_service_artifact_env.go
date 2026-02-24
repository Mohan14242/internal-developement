package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"src/src/internal/db"
)

type ArtifactResponse struct {
	Version    string    `json:"version"`
	ArtifactID string    `json:"artifactId"`
	CreatedAt  time.Time `json:"createdAt"`
}

func GetServiceArtifacts(w http.ResponseWriter, r *http.Request) {
	// 🔒 Allow GET only
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Expected path: /services/{serviceName}/artifacts
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "artifact-by-env" || parts[2] != "artifacts" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	serviceName := parts[1]

	// Query param: environment
	environment := r.URL.Query().Get("environment")
	if environment == "" {
		http.Error(w, "environment query param is required", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(
		`SELECT version, artifact_id, created_at
		 FROM artifacts
		 WHERE service_name = ? AND environment = ?
		 ORDER BY created_at DESC`,
		serviceName, environment,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var artifacts []ArtifactResponse

	for rows.Next() {
		var a ArtifactResponse
		if err := rows.Scan(&a.Version, &a.ArtifactID, &a.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		artifacts = append(artifacts, a)
	}

	if len(artifacts) == 0 {
		http.Error(w, "no artifacts found for service/environment", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artifacts)
}