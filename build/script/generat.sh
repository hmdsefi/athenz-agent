#!/usr/bin/env bash

export DEST_DIR=$GOPATH/src
export WORKSPACE=$GOPATH/src/github.com/hamed-yousefi/athenz-agent/grpc

for file in  $(find ${WORKSPACE}/proto -type f -name *.proto)
do
    echo $file
#    protoc -I $WORKSPACE -I $DEST_DIR --go_out=plugins=grpc:$DEST_DIR $file
    protoc -I $WORKSPACE  --go_out=$DEST_DIR --go-grpc_out=require_unimplemented_servers=false:$DEST_DIR $file

done
