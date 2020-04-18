package main

import (
	"context"
	"log"

	blogpb "github.com/HarshalVoonna/grpc-golang/blog-with-grpc/blog-protobuf"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Blog Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}

	//defer keyword Executes below statement at very end
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	log.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Shreeharsha",
		Title:    "GRPC and MongoDB",
		Content:  "First blog with GRPC and MongoDB",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Unexpected error creating Blog %v", err)
	}
	log.Printf("Blog has been created %v\n", res)

}
