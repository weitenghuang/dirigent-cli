#!/usr/bin/env bash

set -e

docker build -t dirigent-cli-install .
docker run --name dirigent-cli dirigent-cli-install go build -o /go/bin/dirigent github.com/weitenghuang/dirigent-cli
docker commit --change 'CMD ["dirigent"]' dirigent-cli dirigent-cli:$1
docker rm dirigent-cli
