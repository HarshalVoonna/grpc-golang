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

func (s *server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function was invoked with %v\n", req)
	inputNumber := req.GetInputNumber()
	var p int32 = 2
	// p := int64(2)
	for inputNumber > 1 {
		if inputNumber%p == 0 {
			resp := &calculatorpb.PrimeNumberDecompositionResponse{
				PrimeNumberDecompositionResult: p,
			}
			stream.Send(resp)
			inputNumber = inputNumber / p
		} else {
			p = p + 1
		}
	}
	return nil
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
