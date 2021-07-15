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
    
    NEXUS_AUTH="$(urlencode "${NEXUS_USERNAME}"):$(urlencode "${NEXUS_PASSWORD}")"
    
    pip3 config set global.index-url https://${NEXUS_AUTH}@${NEXUS_HOST}/repository/pypi-all/simple
    pip3 config set global.trusted-host ${NEXUS_HOST}
    pip3 config set global.extra-index-url https://pypi.org/simple
fi;

printf "\nInstall test requirements\n" 
. /opt/venv/bin/activate
pip install --upgrade pip
pip install -r tests_requirements.txt
pip check

printf "\nExecute linting\n"
mypy src
flake8 --max-line-length=120 src

printf "\nExecute testing\n"
mkdir -p build/test-results/test
mkdir -p build/test-results/coverage
PYTHONPATH=src python -m pytest --junitxml=build/test-results/test/report.xml -o junit_family=xunit2 --cov-report term-missing --cov-report xml:build/test-results/coverage/coverage.xml --cov=src -o testpaths=tests

# xunit test report
mkdir -p .ods/artifacts/xunit-reports
cat build/test-results/test/report.xml
cp build/test-results/test/report.xml .ods/artifacts/xunit-reports/report.xml

# code coverage
mkdir -p .ods/artifacts/code-coverage
cat build/test-results/coverage/coverage.xml
cp build/test-results/coverage/coverage.xml .ods/artifacts/code-coverage/coverage.xml

printf "\nCopy src and requirements.txt to docker/app\n"
cp -rv src docker/app
cp -rv requirements.txt docker/app
