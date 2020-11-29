#!/bin/env bash

# Ensure user has protoc-gen-grpc-web plugin
path_to_plugin=$(which protoc-gen-grpc-web)
if [ ! -x "$path_to_plugin" ] ; then
	echo "Cannot find protoc-gen-grpc-web plugin.  Please download from https://github.com/grpc/grpc-web/releases and make it discoverable from your PATH."
	exit 1;
fi

cd proto
protoc \
	--js_out="import_style=commonjs,binary:../js_web/" \
	--grpc-web_out="import_style=commonjs,mode=grpcwebtext:../js_web/" \
	auth.proto \
	message.proto \
	user.proto 
