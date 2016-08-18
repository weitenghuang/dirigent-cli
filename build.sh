#!/usr/bin/env bash

set -e

version="latest"
if [ ! -z "$1" ]; then
  version=$1
fi
echo ">>> Build dirigent development image."
docker build -t dirigent-cli-development .
echo ">>> Build dirigent cli command image."
docker run --rm dirigent-cli-development | docker build -t dirigent-cli:$version --no-cache -f Dockerfile.install -
