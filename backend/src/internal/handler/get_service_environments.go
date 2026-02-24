package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"src/src/internal/db"
)

func GetServiceEnvironments(w http.ResponseWriter, r *http.Request) {
	// 🔒 Allow GET only
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Expected path: /api/services/{serviceName}/environments
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	serviceName := parts[1]

	rows, err := db.DB.Query(
		`SELECT DISTINCT environment
		 FROM environment_state
		 WHERE service_name = ?`,
		serviceName,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var environments []string

	for rows.Next() {
		var env string
		if err := rows.Scan(&env); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		environments = append(environments, env)
	}

	if len(environments) == 0 {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(environments)
}

