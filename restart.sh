#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"

go mod tidy
go build

pkill scc
