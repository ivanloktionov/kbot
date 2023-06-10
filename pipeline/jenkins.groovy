pipeline {
    agent any
    parameters {

        choice(name: 'OS', choices: ['linux', 'arm', 'windows', 'macos'], description: 'Pick OS')
        choice(name: 'ARCH', choices: ['amd64', 'arm64'], description: 'Pick Arch')

    }
    stages {
        stage('Buiild') {
            steps {
                echo "Build for platform ${params.OS}"
                sh 'make build ${params.OS}'    
                echo "Build for arch: ${params.ARCH}"
                sh 'make build ${params.ARCH}' 
            }
        }
    }
}
