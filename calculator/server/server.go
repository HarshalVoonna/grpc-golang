package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/HarshalVoonna/grpc-golang/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)
	firstOperator := req.GetFirstOperand()
	secondOperator := req.GetSecondOperand()
	result := firstOperator + secondOperator
	resp := &calculatorpb.SumResponse{
		SumResult: result,
	}
	return resp, nil
}

func main() {
	fmt.Println("Hello from Server side")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
