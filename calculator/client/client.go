package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/HarshalVoonna/grpc-golang/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello from Client side")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server %v\n", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	sum(c)

	primeNumberDecomposition(c)

	computeAverage(c)

	findMaximum(c)
}

func sum(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC")
	req := &calculatorpb.SumRequest{
		FirstOperand:  3,
		SecondOperand: 10,
	}
	resp, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to calculate sum: %v\n", err)
	}
	fmt.Println("Sum obtained from server is", resp.SumResult)
}

func primeNumberDecomposition(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		InputNumber: 210,
	}
	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to calculate Prime Number Decomposition: %v\n", err)
	}
	for {
		resp, respErr := resStream.Recv()
		if respErr == io.EOF {
			fmt.Printf("Reached end of response from Server side: %v\n", respErr)
			break
		} else if respErr != nil {
			log.Fatalf("Failed to stream Server response: %v\n", respErr)
		} else {
			fmt.Println("Prime number decomposition is", resp.GetPrimeNumberDecompositionResult())
		}
	}
}

func computeAverage(c calculatorpb.CalculatorServiceClient) {
	inputNumbers := []int32{1, 2, 3, 4}
	reqStream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Failed to send ComputeAverage RPC request to Server: %v\n", err)
	}
	for _, val := range inputNumbers {
		req := &calculatorpb.ComputeAverageRequest{
			InputNumber: val,
		}
		err := reqStream.Send(req)
		if err != nil {
			log.Fatalf("Failed to send Client-Streaming request to Server: %v\n", err)
		}
		fmt.Printf("Sent inputNumber %v\n", val)
		time.Sleep(time.Second * 1)
	}
	resp, err := reqStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to close Client-Streaming RPC: %v\n", err)
	}
	fmt.Printf("Computed average of numbers sent is %v\n", resp.GetAverageResult())
}

func findMaximum(c calculatorpb.CalculatorServiceClient) {
	inputNumbers := []int32{1, 5, 3, 6, 2, 20}
	reqStream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Failed to send FindMaximum RPC request to Server: %v\n", err)
		return
	}

	lockChannel := make(chan struct{})
	// Send
	go func() {
		for _, num := range inputNumbers {
			log.Println("Sending input number", num)
			err := reqStream.Send(&calculatorpb.FindMaximumRequest{
				InputNumber: num,
			})
			if err != nil {
				log.Fatalf("Failed to stream request to Server: %v\n", err)
				break
			}
			time.Sleep(time.Second * 1)
		}
		reqStream.CloseSend()
	}()

	//Recv
	go func() {
		for {
			resp, err := reqStream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("Failed to receive response from Server: %v\n", err)
				break
			} else {
				log.Println("Max numer till now is", resp.GetMaxNumber())
			}
		}
		close(lockChannel)
	}()

	<-lockChannel
}
