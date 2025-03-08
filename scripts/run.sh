#!/bin/sh

set -e

PROJECT_ROOT="$(dirname "$(dirname "$0")")"
(
  cd "$PROJECT_ROOT"
  go build -o /tmp/codecrafters-build-redis-go ./cmd/redis
)

exec /tmp/codecrafters-build-redis-go "$@"
