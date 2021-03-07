#!/bin/env bash
cd proto
protoc --go_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:../go \
       --go-grpc_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:../go \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative \
       user.proto \
       message.proto