package service

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type RepoConfig struct {
	ServiceName  string   `json:"serviceName"`
	RepoURL      string   `json:"repoUrl"`
}

func UpdateConfigJSON(repoPath string, serviceName string, repoURL string) error {
	configPath := filepath.Join(repoPath, "config.json")

	// Read existing file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Use map to avoid touching other fields
	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	// Update ONLY required fields
	cfg["serviceName"] = serviceName
	cfg["repoUrl"] = repoURL

	// Write back
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, out, 0644)
}


