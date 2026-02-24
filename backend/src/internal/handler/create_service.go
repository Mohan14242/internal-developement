package handler

import (
	"io"
	"log"
	"net/http"
	"strings"

	"gopkg.in/yaml.v3"

	"src/src/internal/model"
	"src/src/internal/service"
)

func CreateService(w http.ResponseWriter, r *http.Request) {
	log.Println("📥 /create-service (YAML) request received")

	// Allow only POST
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Content-Type check (allow charset)
	ct := r.Header.Get("Content-Type")
	if !strings.Contains(ct, "yaml") {
		http.Error(w, "Content-Type must be application/x-yaml", http.StatusUnsupportedMediaType)
		return
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	// Parse YAML
	var req model.CreateServiceRequest
	if err := yaml.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid YAML format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.ServiceName == "" ||
		req.RepoName == "" ||
		req.OwnerTeam == "" ||
		req.Runtime == "" ||
		req.TemplateVersion == "" ||
		req.CICDType == "" ||
		len(req.Environments) == 0 {
		http.Error(
			w,
			"ServiceName, repoName, ownerTeam, runtime, templateVersion, cicdType,Environments are required",
			http.StatusBadRequest,
		)
		return
	}

	// Log (FIXED format)
	log.Printf(
		"🧾 payload → service=%s repo=%s owner=%s runtime=%s template=%s cicd=%s Deployment_type=%s environments=%s",
		req.ServiceName,
		req.RepoName,
		req.OwnerTeam,
		req.Runtime,
		req.TemplateVersion,
		req.CICDType,
		req.DeployType,
		req.Environments,
	)

	// Call service layer
	repoURL, err := service.CreateService(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"repoUrl":"` + repoURL + `"}`))
}