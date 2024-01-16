pipeline {
    agent {
        kubernetes {
            yaml '''
                apiVersion: v1
                kind: Pod
                spec:
                  containers:
                  - name: agent
                    image: jenkins/inbound-agent
                    command:
                    - cat
                    tty: true
                  - name: docker
                    image: docker:latest
                    command:
                    - cat
                    tty: true
                    volumeMounts:
                    - mountPath: /var/run/docker.sock
                      name: docker-sock
                  - name: kubectl
                    image: lachlanevenson/k8s-kubectl
                    command:
                    - cat
                    tty: true
                  volumes:
                  - name: docker-sock
                    hostPath:
                      path: /var/run/docker.sock
            '''
        }
    }

    environment {
        DOCKER_REGISTRY = "iqbal482"
        DOCKER_REPO = "golyrid"
        KUBE_NAMESPACE = "golyrid"
        KUBE_DEPLOYMENT_NAME = "golang-app"
        KUBE_CONTEXT = "jenkins-context"
        BUILD_NUMBER_ENV = "${BUILD_NUMBER}"
        GITHUB_REPO_URL = "https://github.com/moehiqbal/golyrid.git"
        GITHUB_REPO = "moehiqbal/golyrid"
        GIT_CREDENTIALS_ID = "github-credentials"
        DOCKER_HUB_CREDENTIALS_ID = "docker-registry-credentials"
        KUBECONFIG_CREDENTIALS_ID = "kubeconfig-credentials"
    }

    stages {
        stage('Checkout') {
            steps {
                script {
                    // Clone the GitHub repository using credentials
                    withCredentials([usernamePassword(credentialsId: GIT_CREDENTIALS_ID, usernameVariable: 'GIT_USERNAME', passwordVariable: 'GIT_PASSWORD')]) {
                        sh "git clone https://${GIT_USERNAME}:${GIT_PASSWORD}@github.com/${GITHUB_REPO}.git"
                    }
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                container('docker') {
                    script {
                        // Change to the Docker repository directory
                        dir("${DOCKER_REPO}") {
                            // Build Docker image
                            sh "docker build -t ${DOCKER_REGISTRY}/${DOCKER_REPO}:${BUILD_NUMBER_ENV} ."
                        }
                    }
                }
            }
        }

        stage('Push to Docker Registry') {
            steps {
                container('docker') {
                    script {
                        // Log in to Docker Hub and push the image
                        withCredentials([usernamePassword(credentialsId: DOCKER_HUB_CREDENTIALS_ID, usernameVariable: 'DOCKER_HUB_USERNAME', passwordVariable: 'DOCKER_HUB_PASSWORD')]) {
                            sh "docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}"
                            sh "docker push ${DOCKER_REGISTRY}/${DOCKER_REPO}:${BUILD_NUMBER_ENV}"
                        }
                    }
                }
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                script {
                    // Install kubectl and deploy to Kubernetes
                    withCredentials([file(credentialsId: KUBECONFIG_CREDENTIALS_ID, variable: 'KUBECONFIG_FILE')]) {
                        sh "cat /etc/os-release"
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
