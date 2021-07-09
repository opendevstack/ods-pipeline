#!/bin/bash
set -eu

TEST_LOCATION='build/test-results/test'
COVERAGE_LOCATION='build/test-results/coverage'

echo "Prepare Test Suite" 
python -m venv testsuite
. ./testsuite/bin/activate
pip install -r tests_requirements.txt
pip check
mkdir -p ${TEST_LOCATION}
mkdir -p ${COVERAGE_LOCATION}

echo "Lint"
. ./testsuite/bin/activate
mypy src
flake8 --max-line-length=120 src

echo "Test"
. ./testsuite/bin/activate
PYTHONPATH=src python -m pytest --junitxml=report.xml -o junit_family=xunit2 --cov-report term-missing --cov-report xml --cov=src -o testpaths=tests

mkdir -p .ods/artifacts/xunit-reports
cp report.xml .ods/artifacts/xunit-reports/report.xml

mkdir -p .ods/artifacts/code-coverage
cp coverage.xml .ods/artifacts/code-coverage/coverage.xml
cp .coverage .ods/artifacts/code-coverage/.coverage
