#!/bin/bash
set -eu

make docs
if ! git diff --quiet docs; then
    echo "Docs are not up-to-date! Run 'make docs' to update."
    exit 1
else
    echo "Docs are up-to-date."
fi
