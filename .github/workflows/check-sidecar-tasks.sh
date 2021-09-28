#!/bin/bash
set -eu

make sidecar-tasks
if ! git diff --quiet deploy/central/tasks-chart/templates; then
    echo "Sidecar Tasks are not up-to-date! Run 'make sidecar-tasks' to update."
    exit 1
else
    echo "Sidecar Tasks are up-to-date."
fi
