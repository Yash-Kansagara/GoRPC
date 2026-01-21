protoc -I=./proto --go_out=./ --go-grpc_out=./ --grpc-gateway_out=./ ./proto/MathService.proto

