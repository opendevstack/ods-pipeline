#!/bin/bash
set -eu

NO_PROXY="${NO_PROXY:-}"
HTTPS_PROXY="${HTTPS_PROXY:-}"

printf "\nInstall test requirements\n" 
. /opt/venv/bin/activate
pip install --upgrade pip
pip install --proxy $HTTPS_PROXY -r tests_requirements.txt
pip check

printf "\nExecute linting\n"
mypy src
flake8 --max-line-length=120 src

printf "\nExecute testing\n"
PYTHONPATH=src python -m pytest --junitxml=report.xml -o junit_family=xunit2 --cov-report term-missing --cov-report xml --cov=src -o testpaths=tests

mkdir -p .ods/artifacts/xunit-reports
cp report.xml .ods/artifacts/xunit-reports/report.xml

mkdir -p .ods/artifacts/code-coverage
cp coverage.xml .ods/artifacts/code-coverage/coverage.xml

printf "\nCopy src and requirements.txt to docker/app\n"
cp -rv src docker/app
cp -rv requirements.txt docker/app
