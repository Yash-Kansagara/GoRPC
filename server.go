package main

import (
	"context"
	"fmt"
	"gorpc/gen/MathService"
	"gorpc/gen/ServerStream"
	"gorpc/gen/StringService"
	"gorpc/utils"
	"log"
	"net"

	"google.golang.org/grpc"
)

// struct implementing unary rpcs
type server struct {
	MathService.UnimplementedMathServiceServer
	StringService.UnimplementedStringServiceServer
}

// struct implementing server side streaming
type streamServer struct {
	ServerStream.UnimplementedRandomStreamServer
}

// Unary RPC
func (s *server) Add(ctx context.Context, req *MathService.AddRequest) (*MathService.AddResponse, error) {
	return utils.Add(ctx, req)
}

// Unary RPC
func (s *server) IsPalindrom(ctx context.Context, req *StringService.PalReq) (*StringService.PalRes, error) {
	return utils.IsPalindrom(ctx, req)
}

// Server side stream
func (s *streamServer) StreamRandoms(req *ServerStream.Req, stream grpc.ServerStreamingServer[ServerStream.Res]) error {
	return utils.StreamRandoms(req, stream)
}

func main() {

	// connect tcp
	port := 5051
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Error listening on port ", err)
	}

	// crate grpc server
	grpcServer := grpc.NewServer()

	// enable reflection on dev
	// reflection.Register(grpcServer)

	// register service implementations
	serverInst := &server{}
	MathService.RegisterMathServiceServer(grpcServer, serverInst)
	StringService.RegisterStringServiceServer(grpcServer, serverInst)
	ServerStream.RegisterRandomStreamServer(grpcServer, &streamServer{})

	log.Printf("listening on :%d", port)

	// Serve
	err = grpcServer.Serve(listner)
	if err != nil {
		log.Fatal("Error on Serve", err)
	}

}
