#!/bin/sh

set -e
(
  cd "$(dirname "$0")"
  go build -o /tmp/shell-go app/*.go
)
exec /tmp/shell-go "$@"
