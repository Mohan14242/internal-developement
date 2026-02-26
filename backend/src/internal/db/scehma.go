package db

import "log"

func EnsureSchema() error {
	log.Println("🧱 Ensuring database schema exists")

	/* ===================== SERVICES ===================== */

	servicesTable := `
	CREATE TABLE IF NOT EXISTS services (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,

		-- Identity
		service_name VARCHAR(150) NOT NULL UNIQUE,

		-- Lifecycle
		status VARCHAR(30) NOT NULL DEFAULT 'creating',
		last_error TEXT NULL,
		provisioned_at TIMESTAMP NULL,

		-- Git
		repo_url VARCHAR(255) NULL,
		repo_name VARCHAR(255) NULL,
		webhook_token VARCHAR(64) NULL,

		-- Ownership & runtime
		owner_team VARCHAR(100) NULL,
		runtime VARCHAR(50) NULL,
		cicd_type VARCHAR(50) NULL,
		template_version VARCHAR(50) NULL,
		deploy_type VARCHAR(50) NULL,

		-- Config
		environments JSON NULL,
		enablewebhook BOOLEAN NOT NULL DEFAULT FALSE,

		-- Timestamps
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

		UNIQUE(service_id, environment),
		FOREIGN KEY (service_id)
			REFERENCES services(id)
			ON DELETE CASCADE
	);`

	/* ===================== ARTIFACTS (HISTORY) ===================== */

	artifactsTable := `
	CREATE TABLE IF NOT EXISTS artifacts (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,

		service_name VARCHAR(150) NOT NULL,
		environment VARCHAR(20) NOT NULL,

		version VARCHAR(255) NOT NULL,  -- image:tag or ami-id
		artifact_type VARCHAR(20) NOT NULL, -- docker | ami
		commit_sha VARCHAR(40) NULL,
		pipeline VARCHAR(30) NULL,          -- jenkins | github
		action VARCHAR(20) NOT NULL,        -- deploy | rollback

		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

		INDEX idx_artifacts_service_env (service_name, environment),
		INDEX idx_artifacts_version (version)
	);`

	/* ===================== ENVIRONMENT STATE (CURRENT) ===================== */

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


	approvaltable := `
	CREATE TABLE deployment_approvals (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		service_name VARCHAR(255) NOT NULL,
		environment VARCHAR(50) NOT NULL,
		status ENUM('pending','approved','rejected') NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		approved_at TIMESTAMP NULL
	);`

	/* ===================== INDEXES ===================== */

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_services_owner_team ON services(owner_team);`,
		`CREATE INDEX IF NOT EXISTS idx_services_status ON services(status);`,
		`CREATE INDEX IF NOT EXISTS idx_services_created_at ON services(created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_deployments_service_id ON deployments(service_id);`,
	}

	/* ===================== EXECUTION ===================== */

	tables := []struct {
		name string
		sql  string
	}{
		{"services", servicesTable},
		{"deployments", deploymentsTable},
		{"artifacts", artifactsTable},
		{"environment_state", environmentStateTable},
		{"approvals",approvaltable},
	}

	for _, t := range tables {
		if _, err := DB.Exec(t.sql); err != nil {
			log.Printf("❌ Failed to ensure %s table: %v\n", t.name, err)
			return err
		}
	}

	for _, idx := range indexes {
		if _, err := DB.Exec(idx); err != nil {
			log.Println("⚠️ Index creation skipped (may already exist):", err)
		}
	}

	log.Println("✅ Database schema is ready (production-grade)")
	return nil
}