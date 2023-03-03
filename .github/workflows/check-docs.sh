#!/bin/bash
set -eu

make docs
if ! git diff --quiet docs tasks; then
    echo "Docs / tasks are not up-to-date! Run 'make docs' to update."
    exit 1
else
    echo "Docs /tasks are up-to-date."
fi
