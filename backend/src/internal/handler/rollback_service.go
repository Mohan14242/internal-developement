package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"src/src/internal/cicd"
	"src/src/internal/db"
)

type RollbackRequest struct {
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

func RollbackService(w http.ResponseWriter, r *http.Request) {
	// 🔒 Allow POST only
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// /services/{serviceName}/rollback
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3  || parts[0] != "rollback-services" || parts[2] != "rollback" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	serviceName := parts[1]

	// Decode body
	var req RollbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Environment == "" || req.Version == "" {
		http.Error(w, "environment and version are required", http.StatusBadRequest)
		return
	}

	// 🔍 Validate artifact exists
	var exists bool
	err := db.DB.QueryRow(`
		SELECT EXISTS (
		  SELECT 1 FROM artifacts
		  WHERE service_name = ? AND environment = ? AND version = ?
		)`,
		serviceName, req.Environment, req.Version,
	).Scan(&exists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "invalid version for environment", http.StatusBadRequest)
		return
	}
	// 🔍 Get CICD type & repo info
	var cicdType, owner, repo string
	err = db.DB.QueryRow(`
		SELECT cicd_type, owner_team, repo_name
		FROM services
		WHERE service_name = ?`,
		serviceName,
	).Scan(&cicdType, &owner, &repo)
	if err != nil {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}
	// 🚀 Trigger rollback via CICD
	switch cicdType {

	case "jenkins":
		err = cicd.TriggerJenkinsRollback(serviceName, req.Environment, req.Version)

	case "github":
		err = cicd.TriggerGitHubRollback(owner, repo, req.Environment, req.Version)

	default:
		http.Error(w, "unsupported cicd type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Async response
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message":"rollback triggered"}`))
}


