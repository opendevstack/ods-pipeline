#!/bin/bash
set -eu

working_dir="."

while [ "$#" -gt 0 ]; do
    case $1 in
    --working-dir) working_dir="$2"; shift;;
    --working-dir=*) working_dir="${1#*=}";;
  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

echo "Checking for sonar-project.properties ..."
if [ ! -f "${working_dir}/sonar-project.properties" ]; then
  echo "No sonar-project.properties present, using default:"
  cat /usr/local/default-sonar-project.properties
  cp /usr/local/default-sonar-project.properties "${working_dir}/sonar-project.properties"
fi
