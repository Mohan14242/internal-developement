package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"log"
	"src/src/internal/cicd"
	"src/src/internal/db"
)

type DeployRequest struct {
	Environment string `json:"environment"`
}

func DeployServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// /services/{serviceName}/deploy
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[2] != "deploy" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	serviceName := parts[1]
	log.Printf("[[deploying the serving]] %s ",serviceName)

	var req DeployRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if req.Environment == "" {
		http.Error(w, "environment required", http.StatusBadRequest)
		return
	}

	// env → branch mapping
	var branch string
	switch req.Environment {
	case "dev":
		branch = "dev"
	case "test":
		branch = "test"
	case "prod":
		branch = "master"
	default:
		http.Error(w, "invalid environment", http.StatusBadRequest)
		return
	}

	log.Printf("[[deploying to the environment is %s]]",req.Environment)

	log.Printf("getting the DB details")

	// 🔍 Get CICD type & repo info
	var cicdType, repo string
	err := db.DB.QueryRow(`
		SELECT cicd_type, repo_name
		FROM services
		WHERE service_name = ?`,
		serviceName,
	).Scan(&cicdType, &repo)
	if err != nil {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}
	log.Printf("Selecting the pipeline type and triggering the deployment")
	// 🚀 Trigger correct CICD
	switch cicdType {
	case "jenkins":
		err = cicd.TriggerJenkinsDeploy(serviceName, branch)

	case "github":
		err = cicd.TriggerGitHubDeploy(repo, branch)

	default:
		http.Error(w, "unsupported cicd type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message":"deployment triggered"}`))
}


