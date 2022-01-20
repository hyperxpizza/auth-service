#! /bin/sh

echo "generating proto from file auth-service.proto"
protoc --go_out=.  --proto_path=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative auth-service.proto