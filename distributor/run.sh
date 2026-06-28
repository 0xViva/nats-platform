#!/bin/sh
set -e

PORT=8080

# kill anything using the port
PID=$(lsof -ti tcp:$PORT)
if [ -n "$PID" ]; then
  kill -9 $PID || true
fi

# wait until port is actually free
while lsof -i tcp:$PORT >/dev/null 2>&1; do
  sleep 0.1
done

templ generate
go run .
