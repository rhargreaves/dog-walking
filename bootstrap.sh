#!/bin/sh
set -euo pipefail

GOOS=linux GOARCH=arm64 go build -o bootstrap main.go
zip api.zip bootstrap
rm bootstrap