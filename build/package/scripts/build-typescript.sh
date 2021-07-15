#!/bin/bash
set -eu

printf "\nnpm ci and build\n" 
# npm install
npm ci
npm run build
mkdir -p docker/dist
cp -r dist docker/dist
cp -r node_modules docker/dist/node_modules

printf "\nRun tests\n" 
npm run test
junit 'artifacts/xunit.xml'
