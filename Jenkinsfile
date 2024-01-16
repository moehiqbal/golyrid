pipeline {
    agent any

    environment {
        // Define environment variables
        DOCKER_REGISTRY = "docker.io"
        DOCKER_REPO = "iqbal482/golyrid"
        KUBE_NAMESPACE = "apps"
        KUBE_DEPLOYMENT_NAME = "golyrid"
        KUBE_CONTEXT = "jenkins-context"
        BUILD_NUMBER_ENV = "${BUILD_NUMBER}"
        GITHUB_REPO_URL = "https://github.com/moehiqbal/golyrid.git"
        GIT_CREDENTIALS_ID = "github-credentials"
        KUBECONFIG_CREDENTIALS_ID = "kubeconfig-credentials"
    }

    stages {
        stage('Checkout') {
            steps {
                // Use specified Git credentials for the repository
                git url: GITHUB_REPO_URL, credentialsId: GIT_CREDENTIALS_ID
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    docker.build("${DOCKER_REGISTRY}/${DOCKER_REPO}:${BUILD_NUMBER_ENV}")
                }
            }
        }

        stage('Push to Docker Registry') {
            steps {
                script {
                    docker.withRegistry("${DOCKER_REGISTRY}", 'docker-registry-credentials') {
                        docker.image("${DOCKER_REGISTRY}/${DOCKER_REPO}:${BUILD_NUMBER_ENV}").push()
                    }
                }
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                script {
                    withCredentials([file(credentialsId: 'your-kubeconfig-credentials-id', variable: 'KUBECONFIG_FILE')]) {
                        sh "kubectl --kubeconfig=${KUBECONFIG_FILE} config use-context ${KUBE_CONTEXT}"
                        sh "kubectl --kubeconfig=${KUBECONFIG_FILE} set image deployment/${KUBE_DEPLOYMENT_NAME} ${KUBE_DEPLOYMENT_NAME}=${DOCKER_REGISTRY}/${DOCKER_REPO}:${BUILD_NUMBER_ENV} -n ${KUBE_NAMESPACE}"
                    }
                }
            }
        }
    }

    post {
        success {
            echo 'Pipeline succeeded! Clean up or notify here if needed.'
        }
        failure {
            echo 'Pipeline failed! Handle cleanup or notifications here.'
        }
    }
}
