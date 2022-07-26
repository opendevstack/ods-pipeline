#!/usr/bin/env bash
set -ue

# This script prints the first argument to stdout, the second argument to
# stderr and exits with the code given as the third argument.
# It can be used in tests to mimick the behaviour of real binaries.
# The strings printed to stdout and stderr are eval'd to facilitate testing of
# environment variables.

eval echo "$1"
>&2 eval echo "$2"
exit "$3"
