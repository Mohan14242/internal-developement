
package git

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func DeleteRepo(token, repoName string) error {
	log.Println("🗑️ Deleting GitHub repository:", repoName)

	owner, err := GetAuthenticatedUser(token)
	log.Println("the owner of the Github reposiroty is",owner)
	if err != nil {
		log.Println("❌ Failed to determine GitHub owner:", err)
		return err
	}

	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s",
		owner,
		repoName,
	)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Println("❌ Failed to create delete request:", err)
		return err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("❌ GitHub delete request failed:", err)
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		log.Println("✅ GitHub repository deleted:", repoName)
		return nil

	case http.StatusNotFound:
		log.Println("⚠️ GitHub repository not found (already deleted):", repoName)
		return nil

	default:
		err := fmt.Errorf("GitHub repo deletion failed with status %d", resp.StatusCode)
		log.Println("❌", err)
		return err
	}
}