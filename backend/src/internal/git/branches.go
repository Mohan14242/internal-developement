package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateBranch(token, owner, repo, newBranch, sourceBranch string) error {
	// 1️⃣ Get source branch SHA
	sha, err := getBranchSHA(token, owner, repo, sourceBranch)
	if err != nil {
		return err
	}

	// 2️⃣ Create new branch ref
	payload := map[string]string{
		"ref": fmt.Sprintf("refs/heads/%s", newBranch),
		"sha": sha,
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs", owner, repo),
		bytes.NewBuffer(body),
	)
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 422 {
		// branch already exists → safe
		return nil
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to create branch %s", newBranch)
	}

	return nil
}


func getBranchSHA(token, owner, repo, branch string) (string, error) {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.github.com/repos/%s/%s/git/ref/heads/%s", owner, repo, branch),
		nil,
	)
	req.Header.Set("Authorization", "token "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("branch %s not found", branch)
	}

	var res struct {
		Object struct {
			SHA string `json:"sha"`
		} `json:"object"`
	}

	json.NewDecoder(resp.Body).Decode(&res)
	return res.Object.SHA, nil
}


