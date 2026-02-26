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
	"src/src/internal/git"
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




func TriggerGitHubDeploy(repo, branch string) error {
	token, err := aws.GetGitToken("git-token")
	if err != nil {
		log.Println("[GITHUB][ERROR] Failed to fetch GitHub token:", err)
		return err
	}

	workflow := "cicd.yaml"

	payload := map[string]interface{}{
		"ref": branch,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	owner, err := git.GetAuthenticatedUser(token)
	if err != nil {
		return err
	}

	fmt.Println("Authenticated GitHub User:", owner)

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/actions/workflows/%s/dispatches",
			owner, repo, workflow,
		),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "platform-backend")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	return nil
}


func TriggerGitHubRollback(repo, environment, version string) error {
	log.Println("[GITHUB][ROLLBACK] Starting GitHub rollback trigger")

	// 🔐 Fetch GitHub token
	token, err := aws.GetGitToken("git-secrete")
	if err != nil {
		log.Printf("[GITHUB][ROLLBACK][ERROR] Failed to fetch GitHub token: %v\n", err)
		return err
	}
	log.Println("[GITHUB][ROLLBACK] GitHub token fetched successfully")

	workflow := "cicd.yaml" // same workflow, handles rollback via inputs

	ref := environmentBranch(environment)
	log.Printf(
		"[GITHUB][ROLLBACK] Using ref=%s for environment=%s\n",
		ref,
		environment,
	)

	payload := map[string]interface{}{
		"ref": ref,
		"inputs": map[string]string{
			"rollback_version": "true",
			"version":          version,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[GITHUB][ROLLBACK][ERROR] Failed to marshal payload: %v\n", err)
		return err
	}

	log.Printf(
		"[GITHUB][ROLLBACK] Dispatch payload: %s\n",
		string(body),
	)

	owner, err := git.GetAuthenticatedUser(token)
	if err != nil {
		log.Printf("[GITHUB][ROLLBACK][ERROR] Failed to get authenticated GitHub user: %v\n", err)
		return err
	}

	log.Printf("[GITHUB][ROLLBACK] Authenticated GitHub user: %s\n", owner)

	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/actions/workflows/%s/dispatches",
		owner,
		repo,
		workflow,
	)

	log.Printf("[GITHUB][ROLLBACK] Dispatch URL: %s\n", url)

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Printf("[GITHUB][ROLLBACK][ERROR] Failed to create HTTP request: %v\n", err)
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	log.Println("[GITHUB][ROLLBACK] Sending request to GitHub")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[GITHUB][ROLLBACK][ERROR] HTTP request failed: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	log.Printf(
		"[GITHUB][ROLLBACK] GitHub response status: %s (%d)\n",
		resp.Status,
		resp.StatusCode,
	)

	if resp.StatusCode != http.StatusNoContent {
		log.Printf(
			"[GITHUB][ROLLBACK][ERROR] GitHub rejected rollback trigger: status=%s\n",
			resp.Status,
		)
		return fmt.Errorf("github rollback trigger failed: %s", resp.Status)
	}

	log.Println("[GITHUB][ROLLBACK][SUCCESS] GitHub rollback workflow dispatched successfully")
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
