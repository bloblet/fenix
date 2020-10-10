#!/usr/bin/env bash

go get github.com/golang/protobuf/protoc-gen-go
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

protoc --go_out=proto/go \
       --go-grpc_out=proto/go \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative \
       proto/user.proto \
       proto/auth.proto \
       proto/message.proto
