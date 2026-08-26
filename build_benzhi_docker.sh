#!/bin/sh
set -eu

IMAGE="${1:-benzhi-trainwash:latest}"

docker build -f benzhi.Dockerfile -t "$IMAGE" .

echo "image built: $IMAGE"
