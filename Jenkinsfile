pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                checkout scm
                script {
                    echo "Checked out branch: ${env.BRANCH_NAME}"
                }
            }
        }

        stage('Static Analysis') {
            steps {
                echo 'Running Linters (Simulated)...'
                echo 'Checking for forbidden files (.env, .pem)...' [cite: 55]
            }
        }

        stage('Build') {
            steps {
                echo 'Building Application (Simulated)...'
            }
        }

        stage('Deploy to Staging') {
            when {
                branch 'beta'
            }
            steps {
                echo 'Deploying to Staging Environment...'
                echo 'Status: Active Integration Trunk Updated.'
            }
        }
        
        stage('Deploy to Production') {
            when {
                branch 'main'
            }
            steps {
                echo 'Deploying to Production (Golden Source)...' 
            }
        }
    }

    post {
        always {
            echo 'Pipeline execution finished.'
        }
        failure {
            echo 'Build Failed! Check logs.'
        }
    }
}
