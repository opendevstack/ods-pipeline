#!/bin/bash
set -eu

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
