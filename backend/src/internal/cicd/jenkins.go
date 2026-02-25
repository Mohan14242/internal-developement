package cicd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type JenkinsClient struct {
	BaseURL string
	User    string
	Token   string
}

// 🔐 Create Jenkins client with normalized base URL
func NewJenkinsClient() *JenkinsClient {
	baseURL := os.Getenv("JENKINS_URL")
	user := os.Getenv("JENKINS_USER")
	token := os.Getenv("JENKINS_API_TOKEN")

	log.Println("[JENKINS] Initializing Jenkins client")
	log.Println("[JENKINS] Raw JENKINS_URL:", baseURL)
	log.Println("[JENKINS] Jenkins user:", user)

	if baseURL == "" || user == "" || token == "" {
		log.Fatal("[JENKINS] Missing required Jenkins environment variables")
	}

	client := &JenkinsClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		User:    user,
		Token:   token,
	}

	log.Println("[JENKINS] Normalized BaseURL:", client.BaseURL)
	return client
}

//
// ─────────────────────────────────────────────
// 🔐 CSRF CRUMB
// ─────────────────────────────────────────────
//

func (j *JenkinsClient) getCrumb() (string, string, error) {
	url := j.BaseURL + "/crumbIssuer/api/json"
	log.Println("[JENKINS] Fetching CSRF crumb from:", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("[JENKINS][ERROR] Failed to create crumb request:", err)
		return "", "", err
	}

	req.SetBasicAuth(j.User, j.Token)

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	log.Println("[JENKINS] Crumb request latency:", time.Since(start))

	if err != nil {
		log.Println("[JENKINS][ERROR] Crumb request failed:", err)
		return "", "", err
	}
	defer resp.Body.Close()

	log.Println("[JENKINS] Crumb response status:", resp.Status)

	if resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("crumb fetch failed: %s", resp.Status)
	}

	var data struct {
		Crumb             string `json:"crumb"`
		CrumbRequestField string `json:"crumbRequestField"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("[JENKINS][ERROR] Failed to decode crumb response:", err)
		return "", "", err
	}

	log.Println("[JENKINS] Crumb field:", data.CrumbRequestField)
	return data.CrumbRequestField, data.Crumb, nil
}

//
// ─────────────────────────────────────────────
// 🚀 CREATE MULTIBRANCH JOB
// ─────────────────────────────────────────────
//


func (j *JenkinsClient) CreateMultibranchJob(
	jobName, repoURL, credentialsID, webhookToken string,
) error {

	log.Println("[JENKINS] Creating multibranch job:", jobName)

	configXML := fmt.Sprintf(`
<org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject plugin="workflow-multibranch">
  <description>Auto-created by Platform</description>

  <properties>
    <com.igalg.jenkins.plugins.mswt.trigger.ComputedFolderWebHookTrigger>
      <token>%s</token>
    </com.igalg.jenkins.plugins.mswt.trigger.ComputedFolderWebHookTrigger>
  </properties>
  
  <orphanedItemStrategy class="com.cloudbees.hudson.plugins.folder.computed.DefaultOrphanedItemStrategy">
    <pruneDeadBranches>true</pruneDeadBranches>
    <daysToKeep>-1</daysToKeep>
    <numToKeep>-1</numToKeep>
  </orphanedItemStrategy>

  <sources class="jenkins.branch.MultiBranchProject$BranchSourceList">
    <data>
      <jenkins.branch.BranchSource>
        <source class="org.jenkinsci.plugins.github_branch_source.GitHubSCMSource">
          <id>%s</id>
          <repoOwner>%s</repoOwner>
          <repository>%s</repository>
          <credentialsId>%s</credentialsId>
        </source>
      </jenkins.branch.BranchSource>
    </data>
  </sources>

  <factory class="org.jenkinsci.plugins.workflow.multibranch.WorkflowBranchProjectFactory">
    <scriptPath>Jenkinsfile</scriptPath>
  </factory>
</org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject>
`,
		webhookToken,
		jobName,
		extractOwner(repoURL),
		extractRepo(repoURL),
		credentialsID,
	)

	endpoint := fmt.Sprintf(
		"%s/createItem?name=%s",
		j.BaseURL,
		url.QueryEscape(jobName),
	)

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(configXML)))
	req.SetBasicAuth(j.User, j.Token)
	req.Header.Set("Content-Type", "application/xml")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("jenkins job creation failed: %s", resp.Status)
	}

	log.Println("[JENKINS] Job created:", jobName)
	return nil
}




func TriggerJenkinsDeploy(jobName, branch string) error {
	jenkinsURL := os.Getenv("JENKINS_URL")
	user := os.Getenv("JENKINS_USER")
	apiToken := os.Getenv("JENKINS_API_TOKEN")

	url := fmt.Sprintf(
		"%s/job/%s/job/%s/build",
		jenkinsURL, jobName, branch,
	)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(user, apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("jenkins trigger failed: %s", resp.Status)
	}
	return nil
}



func TriggerJenkinsRollback(serviceName, environment, version string) error {
	jenkinsURL := os.Getenv("JENKINS_URL")
	token := os.Getenv("JENKINS_API_TOKEN")

	url := fmt.Sprintf(
		"%s/job/%s/buildWithParameters?token=%s&ROLLBACK=true&ROLLBACK_VERSION=%s&ENVIRONMENT=%s",
		jenkinsURL,
		serviceName,
		token,
		version,
		environment,
	)

	resp, err := http.Post(url, "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("jenkins rollback trigger failed: %s", resp.Status)
	}

	return nil
}