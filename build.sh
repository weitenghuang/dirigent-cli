#!/usr/bin/env bash

set -e

version="latest"
if [ ! -z "$1" ]; then
  version=$1
fi
echo ">>> Build dirigent base install image."
docker build -t dirigent-cli-install .
echo ">>> Install dirigent command."
docker run --name dirigent-cli dirigent-cli-install go build -i -x -o /go/bin/dirigent github.com/weitenghuang/dirigent-cli
echo ">>> Create dirigent command image"
docker commit --change 'CMD ["dirigent"]' dirigent-cli dirigent-cli:$version
docker rm dirigent-cli
