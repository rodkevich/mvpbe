#!/bin/bash
set -ex
env CGO_ENABLED=0
env GOOS=darwin
env GOARCH=amd64

mkdir -p bin
go build -ldflags="-w -s" -o bin/items_producer ./cmd/itemsproducer/main.go;
go build -ldflags="-w -s" -o bin/items_processor ./cmd/itemsprocessor/main.go;
