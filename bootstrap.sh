#!/bin/sh
set -euo pipefail

rm -f bootstrap api.zip || true

if [ "$(uname -m)" = "arm64" ]; then
    ARCH="arm64"
else
    ARCH="amd64"
fi
GOARCH="${GOARCH:-$ARCH}"

GOOS=linux go build -o bootstrap main.go
zip api.zip bootstrap