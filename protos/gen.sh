#!/bin/bash

# You might need to go get -v github.com/gogo/protobuf/...

lehajam=${GOPATH-$HOME/go}/src/github.com/lehajam
protos=$lehajam/dgo/protos
pushd $protos > /dev/null
protoc --gofast_out=plugins=grpc:api api.proto
