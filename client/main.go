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
	c := pb.NewAuthClient(conn)

	// Contact the server and get our username accepted.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	loginAck, err := c.Login(ctx, &pb.ClientAuth{Username: "Test"})

	if err != nil {
		log.Fatalf("Failed to log in with username.")
	}

	log.Printf("Logged in as %s.", loginAck.GetUsername())
	log.Printf("Session Token: %s", loginAck.GetSessionToken())
	log.Printf("Expiry: %s", loginAck.GetExpiry().String())
}
