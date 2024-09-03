cd userservice
rm -rf ./generated/grpcinterface
protoc --go_out=./generated/ --go-grpc_out=./generated/ interface/grpc/user-service.proto
cd -