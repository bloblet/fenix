package main

import (
	"context"
	"log"
	"time"

	pb "github.com/bloblet/fenix/proto/6.0.1"

	"google.golang.org/grpc"
)

type Client struct {
	token    string
	username string
}

func (c *Client) keepalive(a pb.AuthClient, username string, updated chan bool) {
	// Timeout of 10 seconds

	// Keepalive loop
	for true {
		// Send login request
		// Timeout of 10 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		loginAck, err := a.Login(ctx, &pb.ClientAuth{Username: username})
		cancel()

		if err != nil {
			panic(err)
		}
		// Update values
		c.token = loginAck.GetSessionToken()
		c.username = loginAck.GetUsername()
		// Send updated signal on channel
		updated <- true

		// Wait until 30 seconds before token expires, for no "jank"
		time.Sleep(loginAck.GetExpiry().AsTime().Sub(time.Now()))

	}
}

func (c *Client) Connect(username string) {
	timeout := 10 * time.Second

	// Set up a connection to the server.

	conn, err := grpc.Dial("bloblet.com:4000", grpc.WithInsecure(), grpc.WithTimeout(timeout), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	a := pb.NewAuthClient(conn)

	// The channel is to make sure we don't try to do anything with a null SessionToken.
	updated := make(chan bool)
	go c.keepalive(a, username, updated)

	for true {
		<-updated

		log.Printf("Logged in as %s.", c.username)
		log.Printf("Session Token: %s", c.token)
	}
}

func main() {
	c := Client{}
	c.Connect("Alice")
}
