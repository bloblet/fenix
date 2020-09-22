package main

import (
	"context"
	"log"
	"time"

	pb "github.com/bloblet/fenix/proto/6.0.1"

	"google.golang.org/grpc"
)

type Client struct {
	token string
	username string

}

func (c *Client) keepalive(a pb.AuthClient, username string, updated chan bool)  {
	// Timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	// Keepalive loop
	for true {
		// Send login request
		loginAck, err := a.Login(ctx, &pb.ClientAuth{Username: username})

		if err == nil {
			// Update values
			c.token = loginAck.GetSessionToken()
			c.username = loginAck.GetUsername()
			// Send updated signal on channel
			updated <- true

			// Wait until 30 seconds before token expires, for no "jank"
			time.Sleep(loginAck.GetExpiry().AsTime().Sub(time.Now()) - (30 * time.Second))
		}
	}
}


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
