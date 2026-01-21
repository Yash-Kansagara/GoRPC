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

type server struct {
	MathService.UnimplementedMathServiceServer
	StringService.UnimplementedStringServiceServer
}

type streamServer struct {
	ServerStream.UnimplementedRandomStreamServer
}

func (s *server) Add(ctx context.Context, req *MathService.AddRequest) (*MathService.AddResponse, error) {
	return utils.Add(ctx, req)
}

func (s *server) IsPalindrom(ctx context.Context, req *StringService.PalReq) (*StringService.PalRes, error) {
	return utils.IsPalindrom(ctx, req)
}

func (s *streamServer) StreamRandoms(req *ServerStream.Req, stream grpc.ServerStreamingServer[ServerStream.Res]) error {
	return utils.StreamRandoms(req, stream)
}

func main() {

	port := 5051
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Error listening on port ", err)
	}

	grpcServer := grpc.NewServer()
	serverInst := &server{}
	MathService.RegisterMathServiceServer(grpcServer, serverInst)
	StringService.RegisterStringServiceServer(grpcServer, serverInst)
	ServerStream.RegisterRandomStreamServer(grpcServer, &streamServer{})

	log.Printf("listening on :%d", port)
	err = grpcServer.Serve(listner)
	if err != nil {
		log.Fatal("Error on Serve", err)
	}

}
