package main

import (
	"log"
	"net/http"

	"src/src/internal/db"
	"src/src/internal/handler"
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
	http.HandleFunc("/api/services/", handler.GetServiceDashboard)
	http.HandleFunc("/service-by-env/", handler.GetServiceEnvironments)
	http.HandleFunc("/artifact-by-env/", handler.GetServiceArtifacts)
	http.HandleFunc("/deploy-services/", handler.DeployServices)
	http.HandleFunc("/rollback-services/", handler.RollbackService)

	log.Println("🚀 Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}