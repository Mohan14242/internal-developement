package repository

import (
	"database/sql"
	"src/src/internal/db"
)

func GetServices() ([]map[string]interface{}, error) {
	rows, err := db.DB.Query(`
		SELECT s.id, s.service_name, s.repo_name, s.owner_team,
		       s.runtime, s.cicd_type, s.template_version, s.deploy_type,
		       d.environment, d.status
		FROM services s
		LEFT JOIN deployments d ON s.id = d.service_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]map[string]interface{})

	for rows.Next() {
		var (
			id int64
			env, status sql.NullString
			serviceName, repoName, ownerTeam, runtime, cicd, tpl, deploy string
		)

		rows.Scan(
			&id, &serviceName, &repoName, &ownerTeam,
			&runtime, &cicd, &tpl, &deploy,
			&env, &status,
		)

		if _, ok := result[id]; !ok {
			result[id] = map[string]interface{}{
				"serviceName":     serviceName,
				"repoName":        repoName,
				"ownerTeam":       ownerTeam,
				"runtime":         runtime,
				"cicdType":        cicd,
				"templateVersion": tpl,
				"deployType":      deploy,
				"environments":    map[string]string{},
			}
		}

		if env.Valid {
			result[id]["environments"].(map[string]string)[env.String] = status.String
		}
	}

	var services []map[string]interface{}
	for _, v := range result {
		services = append(services, v)
	}

	return services, nil
}

func UpdateDeployment(serviceName, env, status string) error {
	_, err := db.DB.Exec(`
		INSERT INTO deployments (service_id, environment, status)
		SELECT id, ?, ? FROM services WHERE service_name = ?
		ON DUPLICATE KEY UPDATE status = ?, updated_at = NOW()
	`, env, status, serviceName, status)

	return err
}