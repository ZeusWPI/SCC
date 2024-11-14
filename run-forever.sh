#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"

while true
do
  PORT=8888 ./scc
  echo 'scc has quit! restarting in 1 second'
  sleep 1
done