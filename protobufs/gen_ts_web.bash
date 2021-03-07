#!/bin/env bash

if [ -z "$PROTOC_GEN_TS_PATH" ]; then 
	export PROTOC_GEN_TS_PATH="/usr/local/lib/node_modules/ts-protoc-gen/bin/protoc-gen-ts";
fi
cd proto
if [ -f "/usr/local/lib/node_modules/ts-protoc-gen/bin/protoc-gen-ts" ]; then
	protoc \
		--plugin="protoc-gen-ts=${PROTOC_GEN_TS_PATH}" \
		--ts_out="service=grpc-web:../ts_web/" \
		--js_out="import_style=commonjs,binary:../ts_web/" \
		message.proto \
		user.proto 
else
	echo "Cannot find protoc-gen-ts file, please install with sudo npm install -g ts-protoc-gen or set the enviromental variable PROTOC_GEN_TS_PATH to point to your installation."
	exit 1;
fi
