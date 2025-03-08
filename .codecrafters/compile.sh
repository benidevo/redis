#!/bin/sh

set -e

PROJECT_ROOT="$(dirname "$(dirname $"0")")"

go build -o /tmp/codecrafters-build-redis-go ./cmd/redis
