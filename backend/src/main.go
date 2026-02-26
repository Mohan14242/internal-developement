package main

import (
	"log"
	"net/http"

	"src/src/internal/db"
	"src/src/internal/handler"
	"strings"
)

func main() {
	db.InitMySQL()
	if err := db.EnsureSchema(); err != nil {
		log.Fatal("❌ Database schema initialization failed:", err)
	}
	
	http.HandleFunc("/create-service", handler.CreateService)
	http.HandleFunc("/services", handler.GetServices)
	http.HandleFunc("/services/", handler.DeployService)
	http.HandleFunc("/artifacts", handler.RegisterArtifact)
	http.HandleFunc("/servicesdashboard/", handler.GetServiceDashboard)
	http.HandleFunc("/service-by-env/", handler.GetServiceEnvironments)
	http.HandleFunc("/artifact-by-env/", handler.GetServiceArtifacts)
	http.HandleFunc("/deploy-services/", handler.DeployServices)
	http.HandleFunc("/rollback-services/", handler.RollbackService)
	http.HandleFunc("/approvals", handler.GetApprovals)
	http.HandleFunc("/approvals/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/approve") {
			handler.ApproveDeployment(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/reject") {
			handler.RejectDeployment(w, r)
			return
		}
		http.NotFound(w, r)
	})

	log.Println("🚀 Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}