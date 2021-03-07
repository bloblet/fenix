#!/usr/bin/env bash

#go get github.com/golang/protobuf/protoc-gen-go
#go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

protoc -I=. \
    --js_out=import_style=commonjs:./web \
    --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./web \
       user.proto \
       message.proto
