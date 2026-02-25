package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"src/src/internal/db"
)

type EnvironmentDashboard struct {
	CurrentVersion *string    `json:"currentVersion"`
	ArtifactID     *string    `json:"artifactId"`
	Status         string     `json:"status"`
	DeployedAt     *time.Time `json:"deployedAt"`
}

type ServiceDashboardResponse struct {
	ServiceName  string                           `json:"serviceName"`
	Environments map[string]EnvironmentDashboard  `json:"environments"`
}

func GetServiceDashboard(w http.ResponseWriter, r *http.Request) {
	// 🔒 Allow GET only
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ⏱️ Timeout protection
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// 🔍 Extract serviceName from URL
	// Expected: /api/services/{serviceName}/dashboard
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "servicesdashboard" || parts[2] != "dashboard" {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}
	serviceName := parts[1]

	rows, err := db.DB.QueryContext(
		ctx,
		`SELECT environment, version, artifact_id, status, deployed_at
		 FROM environment_state
		 WHERE service_name = ?`,
		serviceName,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Default response (all envs not deployed)
	resp := ServiceDashboardResponse{
		ServiceName: serviceName,
		Environments: map[string]EnvironmentDashboard{
			"dev":  emptyEnv(),
			"test": emptyEnv(),
			"prod": emptyEnv(),
		},
	}

	found := false

	for rows.Next() {
		found = true

		var (
			env        string
			version    sql.NullString
			artifactID sql.NullString
			status     string
			deployedAt sql.NullTime
		)

		if err := rows.Scan(&env, &version, &artifactID, &status, &deployedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp.Environments[env] = EnvironmentDashboard{
			CurrentVersion: nullableString(version),
			ArtifactID:     nullableString(artifactID),
			Status:         status,
			DeployedAt:     nullableTime(deployedAt),
		}
	}

	if !found {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}


func emptyEnv() EnvironmentDashboard {
	return EnvironmentDashboard{
		CurrentVersion: nil,
		ArtifactID:     nil,
		Status:         "not_deployed",
		DeployedAt:     nil,
	}
}

func nullableString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullableTime(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}