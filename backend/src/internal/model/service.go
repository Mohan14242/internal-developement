package model

type EnvironmentStatus string

const (
	Success      EnvironmentStatus = "success"
	Failed       EnvironmentStatus = "failed"
	NotDeployed  EnvironmentStatus = "not_deployed"
)


type CreateServiceRequest struct {
	ServiceName     string `yaml:"serviceName" json:"serviceName"`
	RepoName        string `yaml:"repoName" json:"repoName"`
	OwnerTeam       string `yaml:"ownerTeam" json:"ownerTeam"`
	Runtime         string `yaml:"runtime" json:"runtime"`
	CICDType        string `yaml:"cicdType" json:"cicdType"`
	TemplateVersion string `yaml:"templateVersion" json:"templateVersion"`
	DeployType 		string `yaml:"deploytype" json:"deploytype"`
	Environments    []string `json:"environments" yaml:"environments"`
	EnableWebhook   bool     `json:"enableWebhook" yaml:"enableWebhook"`
	// RuntimeVersion string `yaml:"runtimeVersion"`
	// Environment string `yaml:"environment"`
	// Region      string `yaml:"region"`
	// Repository RepositorySpec `yaml:"repository"`
	// CI         CISpec         `yaml:"ci"`
	// Infra      InfraSpec      `yaml:"infra"`
	// Metadata   MetadataSpec   `yaml:"metadata"`
}

// type RepositorySpec struct {
// 	Visibility    string `yaml:"visibility"`
// 	DefaultBranch string `yaml:"defaultBranch"`
// }

// type CISpec struct {
// 	Enabled bool   `yaml:"enabled"`
// 	Tool    string `yaml:"tool"`
// }

// type InfraSpec struct {
// 	Cloud   string `yaml:"cloud"`
// 	Compute string `yaml:"compute"`
// 	CPU     string `yaml:"cpu"`
// 	Memory  string `yaml:"memory"`
// }

// type MetadataSpec struct {
// 	Description  string `yaml:"description"`
// 	CostCenter   string `yaml:"costCenter"`
// 	BusinessUnit string `yaml:"businessUnit"`
// }



type ArtifactEvent struct {
	ServiceName  string `json:"serviceName"`
	Environment  string `json:"environment"`
	Version      string `json:"version"`
	ArtifactType string `json:"artifactType"`
	ArtifactID   string `json:"artifactId"`
	CommitSHA    string `json:"commitSha"`
	Pipeline     string `json:"pipeline"`
	Action       string `json:"action"`   // deploy | rollback
	Status       string `json:"status"`   // success
}


