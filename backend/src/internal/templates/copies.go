package templates

import (
	"fmt"
	"os"
	"path/filepath"
)


func CreateServiceFromTemplate(req TemplateRequest, targetRepo string) error {
	// 1️⃣ Ensure repo exists
	if err := os.MkdirAll(targetRepo, 0755); err != nil {
		return err
	}

	// 2️⃣ Resolve template version path
	versionPath, _, err := GetTemplatePaths(req)
	if err != nil {
		return err
	}

	// 3️⃣ Copy base application template (exclude cicd)
	if err := CopyDirExcept(versionPath, targetRepo, "cicd"); err != nil {
		return fmt.Errorf("copy base template failed: %w", err)
	}

	// 4️⃣ Copy CI/CD based on CICD + DeployType
	switch req.CICD {

	// =======================
	// GitHub Actions
	// =======================
	case "github":
		src := filepath.Join(
			versionPath,
			"cicd",
			"github",
			req.DeployType,
			"workflows",
		)

		if _, err := os.Stat(src); err != nil {
			return fmt.Errorf(
				"github workflows not found for deployType '%s' at %s",
				req.DeployType,
				src,
			)
		}

		dest := filepath.Join(targetRepo, ".github", "workflows")

		if err := os.MkdirAll(dest, 0755); err != nil {
			return err
		}

		if err := CopyDir(src, dest); err != nil {
			return fmt.Errorf("copy github workflows failed: %w", err)
		}

	// =======================
	// Jenkins
	// =======================
	case "jenkins":
		src := filepath.Join(
			versionPath,
			"cicd",
			"jenkins",
			req.DeployType,
			"Jenkinsfile",
		)

		if _, err := os.Stat(src); err != nil {
			return fmt.Errorf(
				"jenkinsfile not found for deployType '%s' at %s",
				req.DeployType,
				src,
			)
		}

		dest := filepath.Join(targetRepo, "Jenkinsfile")

		if err := copyFile(src, dest, 0644); err != nil {
			return fmt.Errorf("copy jenkinsfile failed: %w", err)
		}

	default:
		return fmt.Errorf("unsupported cicd type: %s", req.CICD)
	}

	return nil
}



func CopyDirExcept(src, dest, exclude string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if rel == exclude || filepath.HasPrefix(rel, exclude+string(os.PathSeparator)) {
			return filepath.SkipDir
		}

		target := filepath.Join(dest, rel)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		return copyFile(path, target, info.Mode())
	})
}