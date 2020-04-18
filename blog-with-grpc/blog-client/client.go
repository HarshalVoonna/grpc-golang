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

	blogID := callCreateBlog(c)

	callReadBlog(c, blogID)
	callReadBlog(c, "5e9ac1cc4405dc7ca7592fc2")
}

func callCreateBlog(c blogpb.BlogServiceClient) string {
	log.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Shreeharsha",
		Title:    "GRPC and MongoDB",
		Content:  "First blog with GRPC and MongoDB",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Unexpected error creating Blog %v", err)
	}
	log.Printf("Blog has been created %v\n", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()
	return blogID
}

func callReadBlog(c blogpb.BlogServiceClient, blogID string) {
	log.Println("Reading the blog")
	readBlogRes, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: blogID,
	})
	if err != nil {
		log.Printf("Error occurred while reading: %v", err)
	} else {
		log.Printf("Blog read is %v\n", readBlogRes)
	}
}
