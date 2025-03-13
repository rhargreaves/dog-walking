#!/bin/sh
set -euo pipefail

rm -f bootstrap api.zip || true
GOOS=linux GOARCH=arm64 go build -o bootstrap main.go
zip api.zip bootstrap