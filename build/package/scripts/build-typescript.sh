#!/bin/bash
set -eu

urlencode() {
    local LC_COLLATE=C
    local length="${#1}"
    for (( i = 0; i < length; i++ )); do
        local c="${1:$i:1}"
        case $c in
            [a-zA-Z0-9.~_-]) printf '%s' "$c" ;;
            *) printf '%%%02X' "'$c" ;;
        esac
    done
}

# Remove the protocol segment from NEXUS_URL
NEXUS_HOST=$(echo "${NEXUS_URL}" | sed -E 's/^\s*.*:\/\///g')

if [ ! -z ${NEXUS_HOST} ] && [ ! -z ${NEXUS_USERNAME} ] && [ ! -z ${NEXUS_PASSWORD} ]; then
    
    printf "\nConfiguring npm\n"

    NEXUS_AUTH="$(urlencode "${NEXUS_USERNAME}"):$(urlencode "${NEXUS_PASSWORD}")"
    
    npm config set registry=$NEXUS_URL/repository/npmjs/
    npm config set always-auth=true
    npm config set _auth=$(echo -n $NEXUS_AUTH | base64)
    npm config set email=no-reply@opendevstack.org
    npm config set ca=null
    npm config set strict-ssl=false
fi;

printf "\nnpm ci and build\n" 
npm ci
npm run build
mkdir -p docker/dist
cp -r dist docker/dist

printf "\nCopying node_modules to docker/dist/node_modules...\n"
# TODO: https://github.com/laktak/rsyncy
# rsync -arh --info=progress2 node_modules/ docker/dist/node_modules
cp -r node_modules docker/dist/node_modules

printf "\nRun tests\n" 
npm run test
# TODO: install junit
# junit 'artifacts/xunit.xml'

#TODO: Copy to .ods/artifacts
mkdir -p .ods/artifacts/xunit-reports
cat build/test-results/test/report.xml
cp build/test-results/test/report.xml .ods/artifacts/xunit-reports/report.xml

# code coverage
mkdir -p .ods/artifacts/code-coverage
cat build/coverage/clover.xml
cp build/coverage/clover.xml .ods/artifacts/code-coverage/clover.xml

cat build/coverage/coverage-final.json
cp build/coverage/coverage-final.json .ods/artifacts/code-coverage/coverage-final.json

cat build/coverage/lcov.info
cp build/coverage/lcov.info .ods/artifacts/code-coverage/lcov.info
