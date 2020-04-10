package main

import (
	"fmt"
	"log"

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
	fmt.Print("Created client %f", c)
}
