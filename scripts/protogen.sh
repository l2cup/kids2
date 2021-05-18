#!/bin/bash

protos=( "vc" "network" "node" "errors")

# shellcheck disable=SC2068
for t in ${protos[@]}; do
  rm ./internal/proto/"$t"/"$t".pb.go
  protoc -I="$PWD"/internal/proto/ --go_out=./internal/proto/ "$PWD"/internal/proto/"$t"/"$t".proto --go_opt=paths=source_relative \
	  --go-grpc_out=./internal/proto/ --go-grpc_opt=paths=source_relative

done
