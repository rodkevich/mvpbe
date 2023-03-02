#!/bin/bash
set -ex
env CGO_ENABLED=0
env GOOS=darwin
env GOARCH=amd64

mkdir -p bin
go build -ldflags="-w -s" -o bin/cover ./cmd/coverage/main.go;
go build -ldflags="-w -s" -o bin/items_producer ./cmd/items_producer/main.go;
go build -ldflags="-w -s" -o bin/items_processor ./cmd/items_processor/main.go;
