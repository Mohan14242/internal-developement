package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"io"
)

type githubUser struct {
	Login string `json:"login"`
}

var githubClient = &http.Client{
	Timeout: 15 * time.Second,
}


func CreateRepo(token, repoName string) (string, error) {
	log.Println("📦 Creating GitHub repository:", repoName)

	owner, err := GetAuthenticatedUser(token)
	if err != nil {
		return "", err
	}

	exists, err := RepoExists(token, owner, repoName)
	if err != nil {
		return "", err
	}

	repoURL := fmt.Sprintf("https://github.com/%s/%s", owner, repoName)

	if exists {
		log.Println("⚠️ Repo already exists:", repoURL)
		return repoURL, nil
	}

	body, _ := json.Marshal(map[string]interface{}{
		"name":    repoName,
		"private": false,
	})

	req, err := http.NewRequest(
		"POST",
		"https://api.github.com/user/repos",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "platform-backend")

	client := &http.Client{Timeout: 20 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(
			"repo creation failed: status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}

	log.Println("✅ Repo created:", repoURL)
	return repoURL, nil
}



func GetAuthenticatedUser(token string) (string, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "platform-backend")

	resp, err := githubClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(
			"github api /user failed: status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}

	var user githubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}

	if user.Login == "" {
		return "", fmt.Errorf("github user login is empty")
	}

	return user.Login, nil
}


func RepoExists(token, owner, repoName string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repoName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "platform-backend")

	resp, err := githubClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return false, fmt.Errorf("github auth failed (status %d)", resp.StatusCode)
	default:
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf(
			"repo check failed: status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}
}