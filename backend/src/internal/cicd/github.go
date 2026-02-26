package cicd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"src/src/internal/aws"
)

type GitHubClient struct {
	Token string
}

// ------------------------------------------------------------
// Create GitHub client (token from AWS)
// ------------------------------------------------------------
func NewGitHubClient() (*GitHubClient, error) {
	log.Println("[GITHUB][STEP 1] Fetching GitHub token from AWS Secrets Manager")

	token, err := aws.GetGitToken("git-token")
	if err != nil {
		log.Println("[GITHUB][ERROR] Failed to fetch GitHub token:", err)
		return nil, err
	}

	if token == "" {
		log.Println("[GITHUB][ERROR] GitHub token is EMPTY")
		return nil, fmt.Errorf("github token is empty")
	}

	log.Println("[GITHUB][STEP 1] GitHub token fetched successfully")
	return &GitHubClient{Token: token}, nil
}

// ------------------------------------------------------------
// CreateWebhook – WITH FULL DIAGNOSTIC LOGGING
// ------------------------------------------------------------
func (g *GitHubClient) CreateWebhook(owner, repo, webhookURL string) error {
	startTotal := time.Now()

	log.Println("--------------------------------------------------")
	log.Println("[GITHUB][STEP 2] Creating webhook")
	log.Println("[GITHUB] Repo owner:", owner)
	log.Println("[GITHUB] Repo name :", repo)
	log.Println("[GITHUB] Webhook URL:", webhookURL)

	// 1️⃣ Build payload
	payload := map[string]interface{}{
		"name":   "web",
		"active": true,
		"events": []string{"push"},
		"config": map[string]string{
			"url":          webhookURL,
			"content_type": "json",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Println("[GITHUB][ERROR] Failed to marshal webhook payload:", err)
		return err
	}

	log.Println("[GITHUB][STEP 3] Webhook payload created (size:", len(body), "bytes)")

	// 2️⃣ Create HTTP request
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/hooks",
		owner,
		repo,
	)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("[GITHUB][ERROR] Failed to create HTTP request:", err)
		return err
	}

	req.Header.Set("Authorization", "token "+g.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	log.Println("[GITHUB][STEP 4] HTTP request created")
	log.Println("[GITHUB] POST", url)

	// 3️⃣ Execute request
	client := &http.Client{Timeout: 15 * time.Second}

	startHTTP := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(startHTTP)

	if err != nil {
		log.Println("[GITHUB][ERROR] HTTP request failed:", err)
		log.Println("[GITHUB] HTTP duration:", duration)
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	log.Println("[GITHUB][STEP 5] HTTP response received")
	log.Println("[GITHUB] Status   :", resp.Status)
	log.Println("[GITHUB] Duration :", duration)
	log.Println("[GITHUB] Body     :", strings.TrimSpace(string(respBody)))

	// 4️⃣ Handle errors explicitly
	if resp.StatusCode == 401 {
		log.Println("[GITHUB][ERROR] Unauthorized – token is invalid or expired")
	}

	if resp.StatusCode == 403 {
		log.Println("[GITHUB][ERROR] Forbidden – missing permissions (admin:repo_hook)")
	}

	if resp.StatusCode == 404 {
		log.Println("[GITHUB][ERROR] Repo not found – check owner/repo name")
	}

	if resp.StatusCode == 422 {
		log.Println("[GITHUB][WARN] Webhook already exists OR validation failed")
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf(
			"github webhook creation failed: %s",
			resp.Status,
		)
	}

	log.Println("[GITHUB][SUCCESS] Webhook created successfully")
	log.Println("[GITHUB] Total time:", time.Since(startTotal))
	log.Println("--------------------------------------------------")

	return nil
}




func TriggerGitHubDeploy(owner, repo, branch string) error {
	token, err := aws.GetGitToken("git-secrete")
	workflow := "cicd.yml" // name of workflow file
	payload := map[string]interface{}{
		"ref": branch,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/actions/workflows/%s/dispatches",
			owner, repo, workflow,
		),
		bytes.NewBuffer(body),
	)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("github actions trigger failed: %s", resp.Status)
	}

	return nil
}


func TriggerGitHubRollback(owner, repo, environment, version string) error {
	token, err := aws.GetGitToken("git-secrete")
	workflow := "cicd.yml" // same workflow, handles rollback via inputs
	payload := map[string]interface{}{
		"ref": environmentBranch(environment),
		"inputs": map[string]string{
			"rollback_version": "true",
			"version":  version,
		},
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/actions/workflows/%s/dispatches",
			owner, repo, workflow,
		),
		bytes.NewBuffer(body),
	)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("github rollback trigger failed: %s", resp.Status)
	}

	return nil
}

func environmentBranch(env string) string {
	switch env {
	case "dev":
		return "dev"
	case "test":
		return "test"
	case "prod":
		return "master"
	default:
		return "master"
	}
}
