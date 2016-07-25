# dirigent-cli
`dirigent` command line interface 

[![Build Status](https://travis-ci.org/weitenghuang/dirigent-cli.svg?branch=master)](https://travis-ci.org/weitenghuang/dirigent-cli)  

Create K8s Pods from docker-compose YAML file
-------
1. Build dirigent command:
`$ ./build.sh`

2. Run dirigent deploy command:

```shell
docker run --rm -v /local/docker-compose.yml:/opt/docker-compose.yml \
    -v ./quoins-data:/opt/credentials \
    -v /local/deploy:/opt/deploy \
    -v $PWD:/go/src/github.com/weitenghuang/dirigent-cli \
    -e QUOIN_NAME="your-quoin-name" \
    -e LOAD_BALANCER="https://your-quoin-loadbalancer:443" \
    dirigent-cli:latest dirigent deploy
```
