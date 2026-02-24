package service

import "src/src/internal/repository"

func ListServices() ([]map[string]interface{}, error) {
	return repository.GetServices()
}

func TriggerDeploy(serviceName, env string) error {
	// 1️⃣ Mark deployment as IN_PROGRESS
	if err := repository.UpdateDeployment(serviceName, env, "in_progress"); err != nil {
		return err
	}

	// 2️⃣ Trigger CI/CD asynchronously (later)
	// go triggerPipeline(serviceName, env)

	return nil
}