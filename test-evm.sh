#!/usr/bin/env bash

export SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
set -e

reindex() {
    local network=$1
    cd ../ethereum-http-secured/app/eth-node-$network
    # bring up ETH network (other plans must be down)
    terraform apply -auto-approve
    # wait and reindex
    cd $SCRIPT_DIR/bin/blockstime-server \
        -c $SCRIPT_DIR/config.yml -i eth.$network
    # bring down ETH
    terraform destroy -auto-approve \
        -target module.httpproxy.docker_container.http \
        -target module.ethereum.docker_container.ethereum
}

# These tests works only when ethereum-http-secured is deployed
[[ "$1" == "mainnet"]] && { reindex $1; exit 0; }
[[ "$1" == "rinkeby"]] && { reindex $1; exit 0; }
[[ "$1" == "goerli"]] && { reindex $1; exit 0; }

[[ "$1" == "test" ]] && {
    $(cat ../ethereum-http-secured/app/eth-node-mainnet/credentials.txt)
    export HTTP_ADDR=http://localhost
    export  | grep HTTP
    go test -tags evm -v ./internal/engines/evm/...
}