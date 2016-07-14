#!/usr/bin/env bash

set -e

if [ -d /opt/credentials ]; then
  echo "Setting cluster entry from /opt/credentials ..."
  echo "${QUOIN_NAME:?Please provide QUOIN_NAME value}"
  echo "${LOAD_BALANCER:?Please provide LOAD_BALANCER value}"
  pushd /opt/credentials > /dev/null
  kubectl config set-cluster $QUOIN_NAME --server=$LOAD_BALANCER --certificate-authority=${PWD}/certs/ca-chain.pem
  kubectl config set-credentials $QUOIN_NAME-admin --certificate-authority=${PWD}/certs/ca-chain.pem --client-key=${PWD}/certs/admin-key.pem --client-certificate=${PWD}/certs/admin.pem
  kubectl config set-context $QUOIN_NAME --cluster=$QUOIN_NAME --user=$QUOIN_NAME-admin
  kubectl config use-context $QUOIN_NAME
  popd > /dev/null
  kubectl cluster-info
  echo "Connected to cluster"
fi

exec "$@"