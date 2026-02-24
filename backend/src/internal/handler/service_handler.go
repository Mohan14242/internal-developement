package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"src/src/internal/service"
)

func GetServices(w http.ResponseWriter, _ *http.Request) {
	data, err := service.ListServices()
	if err != nil {
		http.Error(w, "failed to fetch services", 500)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func DeployService(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		http.Error(w, "invalid path", 400)
		return
	}

	serviceName := parts[3]
	env := parts[5]

	if err := service.TriggerDeploy(serviceName, env); err != nil {
		http.Error(w, "deploy failed", 500)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}