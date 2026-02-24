package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetGitToken(secretName string) (string, error) {
	log.Println("🔐 Fetching secret from AWS Secrets Manager:", secretName)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println("❌ Failed to load AWS config:", err)
		return "", err
	}

	client := secretsmanager.NewFromConfig(cfg)

	secret, err := client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		log.Println("❌ Failed to fetch secret:", err)
		return "", err
	}

	if secret.SecretString == nil {
		log.Println("❌ SecretString is nil for secret:", secretName)
		return "", err
	}

	log.Println("✅ Secret fetched successfully:", secretName)

	// ❗ NEVER log secret.SecretString
	return *secret.SecretString, nil
}