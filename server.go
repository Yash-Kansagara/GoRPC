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
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	wg := sync.WaitGroup{}

	// Serve grpc
	go func() {
		wg.Add(1)
		err = grpcServer.Serve(listner)
		if err != nil {
			log.Fatal("Error on Serve", err)
		}
		wg.Done()
	}()

	// serve reverse-proxy http -> grpc
	conn, err := grpc.NewClient("localhost:5051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to create gateway grpc client", err)
	}
	mux := runtime.NewServeMux()
	err = MathService.RegisterMathServiceHandler(context.Background(), mux, conn)

	if err != nil {
		log.Fatal("Failed to register gateway", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}

	go func() {
		wg.Add(1)
		log.Printf("gateway listening on :%d", port)
		err = gwServer.ListenAndServe()
		if err != nil {
			log.Fatal("Error starting gateway server")
		}
		wg.Done()
	}()

	wg.Wait()
}
