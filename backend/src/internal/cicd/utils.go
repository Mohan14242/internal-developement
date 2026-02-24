package cicd

import "strings"

func extractOwner(repoURL string) string {
	parts := strings.Split(repoURL, "/")
	return parts[len(parts)-2]
}

func extractRepo(repoURL string) string {
	name := strings.Split(repoURL, "/")
	return strings.TrimSuffix(name[len(name)-1], ".git")
}