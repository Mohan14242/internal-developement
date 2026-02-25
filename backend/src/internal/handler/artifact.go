package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"src/src/internal/db"
	"src/src/internal/model"
)

func RegisterArtifact(w http.ResponseWriter, r *http.Request) {
	// 🔒 Allow POST only
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 🔒 Limit request body (1MB max)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var req model.ArtifactEvent

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	// 🔒 Hard guardrails
	if req.Status != "success" {
		http.Error(w, "only successful pipelines are accepted", http.StatusBadRequest)
		return
	}

	if err := validateArtifactRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := saveArtifact(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "artifact registered successfully",
	})
}


func validateArtifactRequest(req model.ArtifactEvent) error {
	if req.ServiceName == "" {
		return errors.New("serviceName is required")
	}
	if !isValidEnv(req.Environment) {
		return errors.New("invalid environment")
	}
	if req.Version == "" {
		return errors.New("version is required")
	}
	if req.ArtifactID == "" {
		return errors.New("artifactId is required")
	}

	if req.Action != "deploy" && req.Action != "rollback" {
		return errors.New("action must be deploy or rollback")
	}

	if req.Pipeline != "jenkins" && req.Pipeline != "github" {
		return errors.New("pipeline must be jenkins or github")
	}

	return nil
}



func saveArtifact(a model.ArtifactEvent) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1️⃣ Insert into artifacts (HISTORY)
	_, err = tx.Exec(`
		INSERT INTO artifacts
		(service_name, environment, version, artifact_type,
		 commit_sha, pipeline, action)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		a.ServiceName,
		a.Environment,
		a.Version,
		a.ArtifactType,
		a.ArtifactID,
		a.CommitSHA,
		a.Pipeline,
		a.Action,
	)
	if err != nil {
		return err
	}

	// 2️⃣ Update current environment state (UPSERT)
	_, err = tx.Exec(`
		REPLACE INTO environment_state
		(service_name, environment, version,
		 status, deployed_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		a.ServiceName,
		a.Environment,
		a.Version,
		a.ArtifactID,
		"success",
		time.Now(),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}


func isValidEnv(env string) bool {
	switch env {
	case "dev", "test","pre-prod","prod":
		return true
	default:
		return false
	}
}



