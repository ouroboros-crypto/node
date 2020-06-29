#!/usr/bin/env bash
# Generates builds
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o ./build/ouroborosd ./cmd/ouroborosd
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o ./build/ouroboroscli ./cmd/ouroboroscli
