package main

import (
	"context"
	pb "github.com/bloblet/fenix/proto/6.0.1"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	timeout := 10 * time.Second

	// Set up a connection to the server.

	conn, err := grpc.Dial("bloblet.com:4000", grpc.WithInsecure(), grpc.WithTimeout(timeout), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUsersClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Get(ctx, &pb.Authenticate{Token: "Ayy", ID: "Yay"})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetID())
}
