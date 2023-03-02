#!/bin/bash
set -ex
env CGO_ENABLED=0
env GOOS=darwin
env GOARCH=amd64

mkdir -p bin
go build -ldflags="-w -s" -o bin/items ./cmd/items/main.go;
go build -ldflags="-w -s" -o bin/cover ./cmd/coverage/main.go;
