package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"src/src/internal/cicd"
	"src/src/internal/db"
)

type Approval struct {
	ID          int64     `json:"id"`
	ServiceName string    `json:"serviceName"`
	Environment string    `json:"environment"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	ApprovedAt *time.Time `json:"approvedAt"`
}

func GetApprovals(w http.ResponseWriter, r *http.Request) {
	log.Println("[APPROVAL] Fetching approvals")

	env := r.URL.Query().Get("environment")
	if env == "" {
		http.Error(w, "environment is required", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
		SELECT id, service_name, environment,
		status, created_at
		FROM deployment_approvals
		WHERE environment = ? AND status = 'pending'
		ORDER BY created_at ASC
	`, env)

	if err != nil {
		log.Println("[APPROVAL][ERROR]", err)
		http.Error(w, "failed to fetch approvals", 500)
		return
	}
	defer rows.Close()

	var approvals []Approval

	for rows.Next() {
		var a Approval
		err := rows.Scan(
			&a.ID,
			&a.ServiceName,
			&a.Environment,
			&a.Status,
			&a.CreatedAt,
			&a.ApprovedAt,
		)
		if err != nil {
			log.Println("[APPROVAL][ERROR]", err)
			continue
		}
		approvals = append(approvals, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(approvals)
}




func ApproveDeployment(w http.ResponseWriter, r *http.Request) {
	log.Println("[APPROVAL] Approve request received")

	id, err := extractApprovalID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid approval id", 400)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}
	defer tx.Rollback()

	var service, env string

	err = tx.QueryRow(`
		SELECT service_name, environment
		FROM deployment_approvals
		WHERE id = ? AND status = 'pending'
	`, id).Scan(&service, &env)

	if err == sql.ErrNoRows {
		http.Error(w, "approval not found or already processed", 404)
		return
	}
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	log.Printf("[APPROVAL] Approving service=%s env=%s version=%s\n",
		service, env)

	_, err = tx.Exec(`
		UPDATE deployment_approvals
		SET status='approved', approved_at=NOW()
		WHERE id=?
	`, id)
	if err != nil {
		http.Error(w, "failed to update approval", 500)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "failed to commit approval", 500)
		return
	}

	branch := "master"
	var cicdType, repo string
	err = db.DB.QueryRow(`
		SELECT cicd_type, repo_name
		FROM services
		WHERE service_name = ?`,
		service,
	).Scan(&cicdType, &repo)
	if err != nil {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}
	log.Printf("Selecting the pipeline type and triggering the deployment")
	// 🚀 Trigger correct CICD
	switch cicdType {
	case "jenkins":
		err = cicd.TriggerJenkinsDeploy(service, branch)

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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"production deployment triggered"}`))

}

func RejectDeployment(w http.ResponseWriter, r *http.Request) {
	log.Println("[APPROVAL] Reject request received")

	id, err := extractApprovalID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid approval id", 400)
		return
	}
	_, err = db.DB.Exec(`
		UPDATE deployment_approvals
		SET status='rejected', approved_at=NOW()
		WHERE id=? AND status='pending'
	`, id)

	if err != nil {
		http.Error(w, "failed to reject approval", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"production deployment rejected"}`))
}

func extractApprovalID(path string) (int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	idStr := parts[len(parts)-2]
	return strconv.ParseInt(idStr, 10, 64)
}



