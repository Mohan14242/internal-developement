package cicd

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func RegisterJenkins(
	repoURL, serviceName string,
	enableWebhook bool,
) (string, error) {

	log.Println("[CICD] Registering Jenkins for service:", serviceName)

	jenkins := NewJenkinsClient()

	// GitHub client only needed if webhook enabled
	var github *GitHubClient
	var err error
	if enableWebhook {
		github, err = NewGitHubClient()
		if err != nil {
			log.Println("[CICD][ERROR] Failed to create GitHub client:", err)
			return "", err
		}
	}

	// 🔐 Generate webhook token
	webhookToken, err := GenerateWebhookToken()
	if err != nil {
		log.Println("[CICD][ERROR] Failed to generate webhook token:", err)
		return "", err
	}

	log.Println("[CICD] Webhook token generated")

	// 1️⃣ Create Jenkins job
	log.Println("[CICD] Creating Jenkins multibranch job")
	if err := jenkins.CreateMultibranchJob(
		serviceName,
		repoURL,
		os.Getenv("JENKINS_GITHUB_CREDENTIALS_ID"),
		webhookToken,
	); err != nil {
		log.Println("[CICD][ERROR] Jenkins job creation failed:", err)
		return "", err
	}

	log.Println("[CICD] Jenkins job created successfully")

	// 2️⃣ Create GitHub webhook (optional)
	if enableWebhook {
		webhookURL := fmt.Sprintf(
			"%s/multibranch-webhook-trigger/invoke?token=%s",
			strings.TrimRight(os.Getenv("JENKINS_URL"), "/"),
			webhookToken,
		)

		log.Println("[CICD] Creating GitHub webhook:", webhookURL)

		if err := github.CreateWebhook(
			extractOwner(repoURL),
			extractRepo(repoURL),
			webhookURL,
		); err != nil {
			log.Println("[CICD][ERROR] GitHub webhook creation failed:", err)
			return "", err
		}

		log.Println("[CICD] GitHub webhook created successfully")
	}

	// 🔥 IMPORTANT: return token so caller can persist it
	return webhookToken, nil
}