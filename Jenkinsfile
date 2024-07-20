try {
  node {
    def app
    stage('Clone Repository') {
      final scmVars = checkout(scm)
      env.BRANCH_NAME = scmVars.GIT_BRANCH
      env.SHORT_COMMIT = "${scmVars.GIT_COMMIT[0..7]}"
      env.GIT_REPO_NAME = scmVars.GIT_URL.replaceFirst(/^.*\/([^\/]+?).git$/, '$1')
      env.VERSION = BUILDVERSION()
    }
    
    def registry
    def registryCredentials
    def tag

    if (env.BRANCH_NAME == 'develop') {
        registry = 'https://ocr1.devocp.techbridge.net/'
        registryCredentials = 'service-dev-uat'
        tag = "uat"
    } else if (env.BRANCH_NAME == 'master') {
        registry = 'https://ocr.apps.hqocp.techbridge.net/'
        registryCredentials = 'service-dev-prod'
        tag = "prod"
    }

    if(registry && registryCredentials){
        stage('Build Docker Image') {
            docker.withRegistry('https://techbridge.dkr.ecr.eu-west-1.amazonaws.com/', '') {
                app = docker.build("service-dev/${env.GIT_REPO_NAME}")
            }
        }

        stage("Push Image to OCR ${tag} Registry") {
            retry(3) {
            docker.withRegistry(registry, registryCredentials) {
                app.push("${tag}-${env.SHORT_COMMIT}-${env.VERSION}")
                app.push("latest")
            }
            }
        }
    }

  }
} catch (Error | Exception e) {
  //Finish failing the build after telling someone about it
  throw e
} finally {
  // Post build steps here
  /* Success or failure, always run post build steps */
  // send email
  // publish test results etc etc
}

def BUILDVERSION(){
    timestamp=Calendar.getInstance().getTime().format('YYYYMMddHHmmss',TimeZone.getTimeZone('EAT'))
    return timestamp
}
