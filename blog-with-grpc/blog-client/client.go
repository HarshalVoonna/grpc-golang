package main

import (
	"context"
	"log"

	blogpb "github.com/HarshalVoonna/grpc-golang/blog-with-grpc/blog-protobuf"

	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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

	callUpdateBlog(c, blogID)
	callUpdateBlog(c, "5e9ac1cc4405dc7ca7592fc2")
	callReadBlog(c, blogID)

	callDeleteBlog(c, blogID)
	callReadBlog(c, blogID)
	callDeleteBlog(c, "5e9ac1cc4405dc7ca7592fc2")
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
		log.Fatalf("Unexpected error creating Blog %v\n", err)
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
		log.Printf("Error occurred while reading: %v\n", err)
	} else {
		log.Printf("Blog read is %v\n", readBlogRes)
	}
}

func callUpdateBlog(c blogpb.BlogServiceClient, blogID string) {
	log.Println("Updating the blog")
	blog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Mark",
		Title:    "GRPC and MongoDB Part (Updated)",
		Content:  "Updated blog with GRPC and MongoDB",
	}
	updateBlogRes, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Printf("Error occurred while updating: %v\n", err)
	} else {
		log.Printf("Updated blog is %v\n", updateBlogRes)
	}
}

func callDeleteBlog(c blogpb.BlogServiceClient, blogID string) {
	log.Println("Deleting the blog")
	deleteBlogRes, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: blogID,
	})
	if err != nil {
		log.Printf("Error occurred while deleting: %v\n", err)
	} else {
		log.Printf("Deleted blog ID is %v\n", deleteBlogRes.GetBlogId())
	}
}
