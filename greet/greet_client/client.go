package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/HarshalVoonna/grpc-golang/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello from Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}

	//defer keyword Executes below statement at very end
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client %f\n", c)

	doUnary(c)

	doServerStreaming(c)

	doClientStreaming(c)

	dobidirectionalStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Shreeharsha",
			LastName:  "Voonna",
		},
	}

	resp, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet %v\n", resp.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Shreeharsha",
			LastName:  "Voonna",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// Server has finished sending the response
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes %v\n", msg.Result)
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming LongGreet RPC")
	reqStream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet RPC: %v", err)
	}
	for i := 0; i < 10; i++ {
		err := reqStream.Send(&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Shreeharsha " + strconv.Itoa(i),
				LastName:  "Voonna " + strconv.Itoa(i),
			},
		})
		log.Printf("Sending LongGreetRequest request no. %v\n", i)
		time.Sleep(time.Second * 2)
		if err != nil {
			log.Fatalf("Error while sending Client streaming request: %v", err)
		}
	}
	resp, err := reqStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while getting Server response for LongGreet RPC: %v", err)
	}
	log.Printf("Response from LongGreet %v\n", resp.Result)
}

func dobidirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a bidirectional Streaming GreetEveryone RPC")
	reqStream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetEveryone RPC: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Shreeharsha",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mike",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Adam",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Joseph",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tina",
			},
		},
	}

	waitChannel := make(chan struct{})
	//Go routine to stream requests to Server
	go func() {
		for _, req := range requests {
			log.Printf("Streaming message: %v\n", req)
			err := reqStream.Send(req)
			if err != nil {
				log.Fatalf("Error while streaming request to GreetEveryone RPC: %v", err)
				return
			}
			time.Sleep(time.Second * 1)
		}
		reqStream.CloseSend()
	}()

	//Go routine to receive streaming requests from Server
	go func() {
		for {
			resp, err := reqStream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("Error while receiving response from Server: %v", err)
				break
			} else {
				log.Printf("Received message: %v\n", resp.GetResult())
			}
		}
		close(waitChannel)
	}()

	//lock everything till done
	<-waitChannel
}
