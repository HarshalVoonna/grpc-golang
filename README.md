# grpc-golang
Implementing Unary, ServerStreaming, ClientStreaming, Bidirectional RPCs with the help of GRPC.

> Use commands mentioned in generate.sh file present in root directory to generate pb file for the proto.
> Eg: protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.

> To run server:
go run greet/greet_server/server.go

> To run client:
go run greet/greet_client/client.go