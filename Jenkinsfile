pipeline {
    agent any
    environment {
        PORT = credentials('BACKEND_PORT_BICHNGOC')
        
        DB_HOST = credentials('DB_HOST')
        DB_PORT = credentials('DB_PORT')
        DB_DATABASE = credentials('DB_DATABASE_BICHNGOC')
        DB_USERNAME = credentials('DB_USERNAME')
        DB_PASSWORD = credentials('DB_PASSWORD') 
        DB_ROOT_PASSWORD = credentials('DB_ROOT_PASSWORD')

        JWT_SECRET = credentials('JWT_SECRET') 
        ALLOWED_ORIGINS = credentials('ALLOWED_ORIGINS')

        AWS_REGION = credentials('AWS_REGION')
        AWS_ACCESS_KEY_ID = credentials('AWS_ACCESS_KEY_ID')
        AWS_SECRET_ACCESS_KEY = credentials('AWS_SECRET_ACCESS_KEY')
        AWS_S3_BUCKET = credentials('AWS_S3_BUCKET_BICHNGOC')
        AWS_S3_ORDER_IMAGES_PREFIX = credentials('AWS_S3_ORDER_IMAGES_PREFIX_BICHNGOC')
        
        DOCKER_TAG = 'latest'
        CONTAINER_NAME = 'order-app-be-container'
    }

    stages {
        stage('Remove Old Docker Image') {
            steps {
                script {
                    echo "Stopping and removing old Docker container..."
                    sh "docker stop ${env.CONTAINER_NAME} || true"
                    sh "docker rm ${env.CONTAINER_NAME} || true"
                    
                    echo "Removing old Docker image..."
                    sh "docker rmi order-app-be:${env.DOCKER_TAG} || true"
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    sh "docker build -t order-app-be:${env.DOCKER_TAG} ."
                }
            }
        }

        stage('Run Container') {
            steps {
                script {
                    sh """
                        docker run -d \\
                            --restart unless-stopped \\
                            --name ${env.CONTAINER_NAME} \\
                            -p ${env.PORT}:${env.PORT} \\
                            -e PORT="${env.PORT}" \\
                            -e DB_HOST="${env.DB_HOST}" \\
                            -e DB_PORT="${env.DB_PORT}" \\
                            -e DB_DATABASE="${env.DB_DATABASE}" \\
                            -e DB_USERNAME="${env.DB_USERNAME}" \\
                            -e DB_PASSWORD="${env.DB_PASSWORD}" \\
                            -e DB_ROOT_PASSWORD="${env.DB_ROOT_PASSWORD}" \\
                            -e AWS_REGION="${env.AWS_REGION}" \\
                            -e AWS_ACCESS_KEY_ID="${env.AWS_ACCESS_KEY_ID}" \\
                            -e AWS_SECRET_ACCESS_KEY="${env.AWS_SECRET_ACCESS_KEY}" \\
                            -e AWS_S3_BUCKET="${env.AWS_S3_BUCKET}" \\
                            -e AWS_S3_ORDER_IMAGES_PREFIX="${env.AWS_S3_ORDER_IMAGES_PREFIX}" \\
                            -e JWT_SECRET="${env.JWT_SECRET}" \\
                            -e ALLOWED_ORIGINS="${env.ALLOWED_ORIGINS}" \\
                            order-app-be:${env.DOCKER_TAG}
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
                sh "docker stop ${env.CONTAINER_NAME} || true"
                sh "docker rm ${env.CONTAINER_NAME} || true"
                cleanWs()
            }
        }
        always {
            echo 'Pipeline completed'
            cleanWs()
        }
    }
}
