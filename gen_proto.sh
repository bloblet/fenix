#!/usr/bin/env bash

go get github.com/golang/protobuf/protoc-gen-go
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
protoc --go_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
       --go-grpc_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative \
       proto/6.0.1/user.proto \
       proto/6.0.1/auth.proto \
       proto/6.0.1/message.proto
