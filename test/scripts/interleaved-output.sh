#!/usr/bin/env bash
set -ue

bash -c 'sleep 1; >&2 echo stderr after sleep' &
bash -c 'sleep 2; echo stdout after sleep' &

echo "some stdout"
sleep 0.1
>&2 echo "some stderr"
sleep 0.1
echo "more stdout"
sleep 0.1
>&2 echo "more stderr"

exit 0
