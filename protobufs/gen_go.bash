#!/bin/bash
cd proto

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     INCLUDE_DIR="/usr/include/";;
    Darwin*)    INCLUDE_DIR="/opt/homebrew/Cellar/protobuf/3.15.8/include/";
esac


protoc --go_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:../go \
       --go-grpc_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:../go \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative \
       -I "$INCLUDE_DIR" -I . \
       user.proto \
       message.proto \
       channels.proto