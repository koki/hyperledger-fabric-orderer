#!/bin/bash

cd $(dirname $0)/..

docker build -t build_image .

docker run -v $(pwd):/go/src/github.com/koki/hyperledger-fabric-orderer build_image
