#!/usr/bin/env bash

# These tests works only when ethereum-http-secured is deployed

$(cat ../ethereum-http-secured/app/eth-node-mainnet/credentials.txt)
export HTTP_ADDR=http://localhost
export  | grep HTTP
go test -tags evm -v ./internal/engines/evm/...