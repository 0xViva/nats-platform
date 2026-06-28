#!/bin/bash

mkdir -p jetstream

docker run -d --name nats-server \
  -p 4222:4222 \
  -p 8222:8222 \
  -v "$PWD/nats.conf:/etc/nats/nats.conf" \
  -v "$PWD/jetstream:/data" \
  nats:latest \
  -c /etc/nats/nats.conf
