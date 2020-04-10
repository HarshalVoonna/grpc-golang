package main

import (
	"context"
	"fmt"
	"log"

	"github.com/HarshalVoonna/grpc-golang/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello from Client side")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	sum(c)
}

func sum(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC")
	req := &calculatorpb.SumRequest{
		FirstOperand:  3,
		SecondOperand: 10,
	}
	resp, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to calculate sum: %v", err)
	}
	fmt.Println("Sum obtained from server is", resp.SumResult)

}
