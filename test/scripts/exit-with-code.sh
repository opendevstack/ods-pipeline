#!/usr/bin/env bash
set -ue

# This script prints the first argument to stdout, prints the second argument to
# stderr and exits with the exit code given as the third argument.
# It can be used in tests to mimick the behaviour of real binaries.

echo $1
>&2 echo $2
exit $3
