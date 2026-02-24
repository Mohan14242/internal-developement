package templates

import (
	"fmt"
	"os"
	"path/filepath"
)

// 🔑 Resolve template root dynamically (portable)
func templateRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(
		wd,
		"src",
		"internal",
		"template_data",
	), nil
}


type TemplateRequest struct {
	Language   string
	Version    string
	CICD       string          // github | jenkins
	DeployType string          // ec2 | microservice
}


// ✅ Validate & return template paths
func GetTemplatePaths(req TemplateRequest) (versionPath, cicdPath string, err error) {
	if req.Language == "" || req.Version == "" || req.CICD == "" || req.DeployType == "" {
		return "", "", fmt.Errorf("language, version, cicd, and deployType are required")
	}

	root, err := templateRoot()
	if err != nil {
		return "", "", err
	}

	// 1️⃣ Language
	langPath := filepath.Join(root, req.Language)
	if _, err := os.Stat(langPath); err != nil {
		return "", "", fmt.Errorf("language '%s' not found", req.Language)
	}

	// 2️⃣ Version
	versionPath = filepath.Join(langPath, req.Version)
	if _, err := os.Stat(versionPath); err != nil {
		return "", "", fmt.Errorf(
			"version '%s' not found for language '%s'",
			req.Version,
			req.Language,
		)
	}

	// 3️⃣ CICD + DeployType validation
	cicdPath = filepath.Join(
		versionPath,
		"cicd",
		req.CICD,
		req.DeployType,
	)

	if _, err := os.Stat(cicdPath); err != nil {
		return "", "", fmt.Errorf(
			"deployType '%s' not supported for cicd '%s'",
			req.DeployType,
			req.CICD,
		)
	}


	

	// 4️⃣ Ensure app source exists
	if _, err := os.Stat(filepath.Join(versionPath, "src")); err != nil {
		return "", "", fmt.Errorf("src folder missing in template")
	}

	return versionPath, cicdPath, nil
}