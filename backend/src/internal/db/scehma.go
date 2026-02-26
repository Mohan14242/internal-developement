package db

import "log"

func EnsureSchema() error {
	log.Println("🧱 Ensuring database schema exists")

	/* ===================== SERVICES ===================== */

	servicesTable := `
	CREATE TABLE IF NOT EXISTS services (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,

		service_name VARCHAR(150) NOT NULL UNIQUE,

		status VARCHAR(30) NOT NULL DEFAULT 'creating',
		last_error TEXT NULL,
		provisioned_at TIMESTAMP NULL,

		repo_url VARCHAR(255) NULL,
		repo_name VARCHAR(255) NULL,
		webhook_token VARCHAR(64) NULL,

		owner_team VARCHAR(100) NULL,
		runtime VARCHAR(50) NULL,
		cicd_type VARCHAR(50) NULL,
		template_version VARCHAR(50) NULL,
		deploy_type VARCHAR(50) NULL,

		environments JSON NULL,
		enablewebhook BOOLEAN NOT NULL DEFAULT FALSE,

		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			ON UPDATE CURRENT_TIMESTAMP
	);`

	/* ===================== DEPLOYMENTS ===================== */

	deploymentsTable := `
	CREATE TABLE IF NOT EXISTS deployments (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,

		service_id BIGINT NOT NULL,
		environment VARCHAR(20) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'not_deployed',

		last_deployed_at TIMESTAMP NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			ON UPDATE CURRENT_TIMESTAMP,

		UNIQUE KEY uniq_service_env (service_id, environment),
		FOREIGN KEY (service_id)
			REFERENCES services(id)
			ON DELETE CASCADE
	);`

	/* ===================== ARTIFACTS ===================== */

	artifactsTable := `
	CREATE TABLE IF NOT EXISTS artifacts (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,

		service_name VARCHAR(150) NOT NULL,
		environment VARCHAR(20) NOT NULL,

		version VARCHAR(255) NOT NULL,
		artifact_type VARCHAR(20) NOT NULL,
		commit_sha VARCHAR(40) NULL,
		pipeline VARCHAR(30) NULL,
		action VARCHAR(20) NOT NULL,

		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

		INDEX idx_artifacts_service_env (service_name, environment),
		INDEX idx_artifacts_version (version)
	);`

	/* ===================== ENVIRONMENT STATE ===================== */

	environmentStateTable := `
	CREATE TABLE IF NOT EXISTS environment_state (
		service_name VARCHAR(150) NOT NULL,
		environment VARCHAR(20) NOT NULL,

		version VARCHAR(255) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'success',
		deployed_at TIMESTAMP NOT NULL,

		PRIMARY KEY (service_name, environment),
		INDEX idx_env_state_service (service_name)
	);`

	/* ===================== DEPLOYMENT APPROVALS ===================== */

	approvalsTable := `
	CREATE TABLE IF NOT EXISTS deployment_approvals (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,

		service_name VARCHAR(150) NOT NULL,
		environment VARCHAR(50) NOT NULL,
		status ENUM('pending','approved','rejected') NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		approved_at TIMESTAMP NULL,
		INDEX idx_approvals_env_status (environment, status),
		INDEX idx_approvals_service (service_name)
	);`

	/* ===================== EXECUTION ===================== */

	tables := []struct {
		name string
		sql  string
	}{
		{"services", servicesTable},
		{"deployments", deploymentsTable},
		{"artifacts", artifactsTable},
		{"environment_state", environmentStateTable},
		{"deployment_approvals", approvalsTable},
	}

	for _, t := range tables {
		if _, err := DB.Exec(t.sql); err != nil {
			log.Printf("❌ Failed to ensure %s table: %v\n", t.name, err)
			return err
		}
	}

	/* ===================== INDEXES (MYSQL SAFE) ===================== */

	indexes := []string{
		`CREATE INDEX idx_services_owner_team ON services(owner_team);`,
		`CREATE INDEX idx_services_status ON services(status);`,
		`CREATE INDEX idx_services_created_at ON services(created_at);`,
		`CREATE INDEX idx_deployments_service_id ON deployments(service_id);`,
	}

	for _, idx := range indexes {
		if _, err := DB.Exec(idx); err != nil {
			// MySQL will throw "Duplicate key name" if index exists — safe to ignore
			log.Println("ℹ️ Index already exists or skipped:", err)
		}
	}

	log.Println("✅ Database schema is ready (production-grade)")
	return nil
}