#!/bin/bash

git fetch
git pull

go mod tidy
echo "Building..."
go build .

pkill scc
