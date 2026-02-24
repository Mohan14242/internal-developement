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
	var req model.ArtifactEvent

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	// 🔒 HARD GUARDRAILS
	if req.Status != "success" {
		http.Error(w, "only successful pipelines are accepted", http.StatusBadRequest)
		return
	}
	if req.ServiceName == "" || req.Environment == "" ||
		req.Version == "" || req.ArtifactID == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	if !isValidEnv(req.Environment) {
		http.Error(w, "invalid environment", http.StatusBadRequest)
		return
	}

	err := saveArtifact(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func saveArtifact(a ArtifactEvent) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1️⃣ Insert into artifacts (HISTORY)
	_, err = tx.Exec(`
		INSERT INTO artifacts
		(service_name, environment, version, artifact_type,
		 artifact_id, commit_sha, pipeline, action)
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
		 artifact_id, status, deployed_at)
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



