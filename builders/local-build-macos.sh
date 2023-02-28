#!/bin/bash
set -ex
env CGO_ENABLED=0
env GOOS=darwin
env GOARCH=amd64

go build -ldflags="-w -s" -o bin/server ./cmd/sample/main.go;
go build -ldflags="-w -s" -o bin/cover ./cmd/coverage/main.go;
