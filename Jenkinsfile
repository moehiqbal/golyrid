pipeline {
    agent any

    environment {
        // Define environment variables
        DOCKER_REGISTRY = "iqbal482"
        DOCKER_REPO = "golyrid"
        KUBE_NAMESPACE = "golyrid"
        KUBE_DEPLOYMENT_NAME = "golang-app"
        KUBE_CONTEXT = "jenkins-context"
        BUILD_NUMBER_ENV = "${BUILD_NUMBER}"
        GITHUB_REPO_URL = "https://github.com/moehiqbal/golyrid.git"
        GIT_CREDENTIALS_ID = "github-credentials"
        KUBECONFIG_CREDENTIALS_ID = "kubeconfig-credentials"
        DOCKER_HUB_CREDENTIALS_ID = "docker-hub-credentials"

    }

    stages {
        stage('Checkout') {
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: GIT_CREDENTIALS_ID, usernameVariable: 'GIT_USERNAME', passwordVariable: 'GIT_PASSWORD')]) {
                        sh "git clone ${GITHUB_REPO_URL}"
                    }
                }
            }
        }


        stage('Build Docker Image') {
            steps {
                script {
                    sh "cd /tmp ${DOCKER_REPO}"
                    sh "docker build -t ${DOCKER_REGISTRY}/${DOCKER_REPO}:${BUILD_NUMBER_ENV} ."
                }
            }
        }

        stage('Push to Docker Registry') {
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: DOCKER_HUB_CREDENTIALS_ID, usernameVariable: 'DOCKER_HUB_USERNAME', passwordVariable: 'DOCKER_HUB_PASSWORD')]) {
                        sh "docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}"
                        sh "docker push ${DOCKER_REGISTRY}/${DOCKER_REPO}:${BUILD_NUMBER_ENV}"
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
