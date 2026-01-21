package main

import (
	"context"
	"fmt"
	"gorpc/gen/MathService"
	"gorpc/gen/StringService"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	MathService.UnimplementedMathServiceServer
	StringService.UnimplementedStringServiceServer
}

func (s *server) Add(ctx context.Context, req *MathService.AddRequest) (*MathService.AddResponse, error) {

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

	return &MathService.AddResponse{
		Ans: string(ans[a+1:]),
	}, nil
}

func (s *server) IsPalindrom(ctx context.Context, req *StringService.PalReq) (*StringService.PalRes, error) {

	l, r := 0, len(req.S)-1
	isPal := true
	for l < r {
		if req.S[l] != req.S[r] {
			isPal = false
			break
		}
		l++
		r--
	}
	return &StringService.PalRes{
		IsPal: isPal,
	}, nil
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

	log.Printf("listening on :%d", port)
	err = grpcServer.Serve(listner)
	if err != nil {
		log.Fatal("Error on Serve", err)
	}

}
