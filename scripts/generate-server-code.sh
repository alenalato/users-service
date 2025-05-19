#!/bin/sh

mkdir -p ./pkg/grpc
docker run \
	-v "$(pwd)":/defs \
	--rm \
	namely/protoc:1.42_2 \
	-I ./proto \
	--go_out ./pkg/grpc \
	--go_opt paths=import \
	--go_opt module=github.com/alenalato/users-service/pkg/grpc \
  --go-grpc_out ./pkg/grpc \
  --go-grpc_opt paths=import \
  --go-grpc_opt module=github.com/alenalato/users-service/pkg/grpc \
  ./proto/*.proto;
