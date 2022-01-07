#!/bin/bash
set -eu

echo "Checking for sonar-project.properties ..."
if [ ! -f sonar-project.properties ]; then
  echo "No sonar-project.properties present, using default:"
  cat /usr/local/default-sonar-project.properties
  cp /usr/local/default-sonar-project.properties sonar-project.properties
fi
