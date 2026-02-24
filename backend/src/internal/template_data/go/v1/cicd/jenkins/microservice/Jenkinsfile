pipeline {
    agent any

    environment {
        SERVICE_NAME = ""
        ENVIRONMENT  = ""
        COMMIT_SHA  = ""
        GIT_TAG     = ""
        VERSION     = ""
        
    }

    stages {

        stage('Load project configuration') {
            steps {
                script {
                    def projectConfig = readJSON file: 'config.json'
                    env.SERVICE_NAME = projectConfig.serviceName
                    env.GITHUB_REPO  = projectConfig.github_repo
                }
            }
        }

        stage('Detect Environment') {
            steps {
                script {
                    if (env.BRANCH_NAME == 'dev') {
                        env.ENVIRONMENT = 'dev'
                    } else if (env.BRANCH_NAME == 'test') {
                        env.ENVIRONMENT = 'test'
                    } else if (env.BRANCH_NAME == 'master' || env.BRANCH_NAME == 'main') {
                        env.ENVIRONMENT = 'prod'
                    } else {
                        error "Unsupported branch: ${env.BRANCH_NAME}"
                    }
                }
            }
        }
        stage('Get Commit SHA') {
            steps {
                script {
                    env.COMMIT_SHA = sh(
                        script: "git rev-parse --short HEAD",
                        returnStdout: true
                    ).trim()
                }
            }
        }

        stage('Get Git Tag Version') {
            steps {
                script {
                    env.GIT_TAG = sh(
                        script: """
                          git describe --tags \$(git rev-list --tags --max-count=1) || echo "v0.0.0"
                        """,
                        returnStdout: true
                    ).trim()

                    echo "Resolved Git Tag: ${env.GIT_TAG}"
                }
            }
        }

        stage('Generate Final Version') {
            steps {
                script {
                    env.VERSION = "${env.SERVICE_NAME}_${env.ENVIRONMENT}_${env.GIT_TAG}_${env.COMMIT_SHA}"
                    echo "Generated VERSION: ${env.VERSION}"
                }
            }
        }

        stage('Docker Build & Push') {
            steps {
                echo "Building Docker image with version ${env.VERSION}"
            }
        }
    }

    post {
        success {
            echo "✅ Pipeline SUCCESS for ${env.VERSION}"

            sh """
              curl -X POST  http://54.196.169.186/api/artifacts \
                -H "Content-Type: application/json" \
                -d '{
                  "serviceName": "${env.SERVICE_NAME}",
                  "environment": "${env.ENVIRONMENT}",
                  "version": "${env.VERSION}",
                  "artifactType": "docker",
                  "artifactId": "myrepo/${env.SERVICE_NAME}:${env.VERSION}",
                  "commitSha": "${env.COMMIT_SHA}",
                  "pipeline": "jenkins",
                  "action":"deploy",
                  "status": "success"
                }'
            """
        }
    }
}