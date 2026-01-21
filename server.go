package main

import (
	"context"
	"fmt"
	pb "gorpc/gen/MathService"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMathServiceServer
}

func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {

	ans := make([]byte, (int)(math.Max(float64(len(req.Op1)), float64(len(req.Op2)))+1))
	for i := range ans {
		ans[i] = 'x'
	}
	p1 := len(req.Op1) - 1
	p2 := len(req.Op2) - 1
	a := len(ans) - 1
	carry := byte(0)
	for p1 > -1 || p2 > -1 {
		sum := byte(0)
		if p1 == -1 {
			sum = req.Op2[p2] - '0' + carry
			if sum > 9 {
				carry = 1
				ans[a] = (sum % 10) + '0'
			} else {
				carry = 0
				ans[a] = sum + '0'
			}
			a--
			p2--
			continue
		}
		if p2 == -1 {
			sum = req.Op1[p1] - '0' + carry
			if sum > 9 {
				carry = 1
				ans[a] = (sum % 10) + '0'
			} else {
				carry = 0
				ans[a] = sum + '0'
			}
			a--
			p1--
			continue
		}
		sum = carry
		sum += req.Op1[p1] - '0'
		sum += req.Op2[p2] - '0'

		if sum > 9 {
			carry = 1
			ans[a] = (sum % 10) + '0'
		} else {
			carry = 0
			ans[a] = sum + '0'
		}
		a--
		p1--
		p2--
	}

	if carry > 0 {
		ans[a] = carry + '0'
		a--
	}

	return &pb.AddResponse{
		Ans: string(ans[a+1:]),
	}, nil
}

func main() {

	port := 5051
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Error listening on port ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMathServiceServer(grpcServer, &server{})

	log.Printf("listening on :%d", port)
	err = grpcServer.Serve(listner)
	if err != nil {
		log.Fatal("Error on Serve", err)
	}

}
