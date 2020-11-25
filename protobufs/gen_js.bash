#!/usr/bin/env bash
cd proto

protoc -I=. \
    --js_out=import_style=commonjs:../web \
    --grpc-web_out=import_style=commonjs,mode=grpcwebtext:../web \
    user.proto \
    auth.proto \
    message.proto
