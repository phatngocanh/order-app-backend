pipeline {
    agent any
    environment {
        PORT = credentials('BACKEND_PORT')
        
        DB_HOST = credentials('DB_HOST')
        DB_PORT = credentials('DB_PORT')
        DB_DATABASE = credentials('DB_DATABASE')
        DB_USERNAME = credentials('DB_USERNAME')
        DB_PASSWORD = credentials('DB_PASSWORD') 
        DB_ROOT_PASSWORD = credentials('DB_ROOT_PASSWORD')
        NOTIFICATION_ACCESS_TOKEN = credentials('NOTIFICATION_ACCESS_TOKEN')

        JWT_SECRET = credentials('JWT_SECRET') 
        ALLOWED_ORIGINS = credentials('ALLOWED_ORIGINS')
    }

    stages {
        stage('Set Environment Tag') {
            steps {
                script {
                    // Set dynamic values
                    def dockerTag = 'main'
                    def exposePort = env.PORT

                    if (env.BRANCH_NAME == 'develop') {
                        dockerTag = 'develop'
                        exposePort = '6868'
                        echo "Environment set to develop with tag 'develop' and EXPOSE_PORT 6868"
                    } else {
                        echo "Environment set to main with tag 'main'"
                    }

                    // Save to env so later stages can access
                    env.DOCKER_TAG = dockerTag
                    env.EXPOSE_PORT = exposePort
                }
            }
        }

        stage('Remove Old Docker Image') {
            steps {
                script {
                    echo "Stopping and removing old Docker container for ${env.DOCKER_TAG}..."
                    sh "docker stop travel-be-container-${env.DOCKER_TAG} || true"
                    sh "docker rm travel-be-container-${env.DOCKER_TAG} || true"
                    
                    echo "Removing old Docker image for ${env.DOCKER_TAG}..."
                    sh "docker rmi travel-be:${env.DOCKER_TAG} || true"
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    sh "docker build -t travel-be:${env.DOCKER_TAG} ."
                }
            }
        }

        stage('Run Container') {
            steps {
                script {
                    sh """
                        docker run -d \\
                            --restart unless-stopped \\
                            --name travel-be-container-${env.DOCKER_TAG} \\
                            -p ${env.EXPOSE_PORT}:${env.PORT} \\
                            -e PORT="${env.PORT}" \\
                            -e DB_HOST="${env.DB_HOST}" \\
                            -e DB_PORT="${env.DB_PORT}" \\
                            -e DB_DATABASE="${env.DB_DATABASE}" \\
                            -e DB_USERNAME="${env.DB_USERNAME}" \\
                            -e DB_PASSWORD="${env.DB_PASSWORD}" \\
                            -e DB_ROOT_PASSWORD="${env.DB_ROOT_PASSWORD}" \\
                            -e MAIL_HOST="${env.MAIL_HOST}" \\
                            -e MAIL_PORT="${env.MAIL_PORT}" \\
                            -e MAIL_USERNAME="${env.MAIL_USERNAME}" \\
                            -e MAIL_PASSWORD="${env.MAIL_PASSWORD}" \\
                            -e MAIL_FROM="${env.MAIL_FROM}" \\
                            -e MAIL_FROM_NAME="${env.MAIL_FROM_NAME}" \\
                            -e REDIS_HOST="${env.REDIS_HOST}" \\
                            -e REDIS_PORT="${env.REDIS_PORT}" \\
                            -e REDIS_PASSWORD="${env.REDIS_PASSWORD}" \\
                            -e JWT_SECRET="${env.JWT_SECRET}" \\
                            -e ALLOWED_ORIGINS="${env.ALLOWED_ORIGINS}" \\
                            -e GEN_TOKEN_URL="${env.GEN_TOKEN_URL}" \\
                            -e CREATE_TOUR_URL="${env.CREATE_TOUR_URL}" \\
                            -e CORE_SECRET_KEY="${env.CORE_SECRET_KEY}" \\
                            -e PLACE_INFO_URL="${env.PLACE_INFO_URL}" \\
                            -e NOTIFICATION_ACCESS_TOKEN="${env.NOTIFICATION_ACCESS_TOKEN}" \\
                            travel-be:${env.DOCKER_TAG}
                    """
                }
            }
        }
    }

    post {
        success {
            echo 'Pipeline completed successfully'
            cleanWs()
        }
        failure {
            echo 'Pipeline failed'
            script {
                sh "docker stop travel-be-container-${env.DOCKER_TAG} || true"
                sh "docker rm travel-be-container-${env.DOCKER_TAG} || true"
                cleanWs()
            }
        }
        always {
            echo 'Pipeline completed'
            cleanWs()
        }
    }
}
