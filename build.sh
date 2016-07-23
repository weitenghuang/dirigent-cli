#!/usr/bin/env bash

set -e

version="latest"
if [ ! -z "$1" ]; then
  version=$1
fi

docker build -t dirigent-cli-install .
docker run --name dirigent-cli dirigent-cli-install go build -o /go/bin/dirigent github.com/weitenghuang/dirigent-cli
docker commit --change 'CMD ["dirigent"]' dirigent-cli dirigent-cli:$version
docker rm dirigent-cli
