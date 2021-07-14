#!/bin/bash
set -eu

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
