package db

import "log"

func EnsureSchema() error {
	log.Println("🧱 Ensuring database schema exists")

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

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_services_owner_team ON services(owner_team);`,
		`CREATE INDEX IF NOT EXISTS idx_services_status ON services(status);`,
		`CREATE INDEX IF NOT EXISTS idx_services_created_at ON services(created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_deployments_service_id ON deployments(service_id);`,
	}

	// Create tables
	if _, err := DB.Exec(servicesTable); err != nil {
		log.Println("❌ Failed to ensure services table:", err)
		return err
	}

	if _, err := DB.Exec(deploymentsTable); err != nil {
		log.Println("❌ Failed to ensure deployments table:", err)
		return err
	}

	// Create indexes
	for _, idx := range indexes {
		if _, err := DB.Exec(idx); err != nil {
			log.Println("⚠️ Failed to create index (may already exist):", err)
		}
	}

	log.Println("✅ Database schema is ready (production-grade)")
	return nil
}