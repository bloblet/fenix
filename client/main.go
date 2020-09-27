package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/bloblet/fenix/proto/6.0.1"
	"github.com/pborman/ansi"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	token         string
	username      string
	conn          *grpc.ClientConn
	messageStream pb.Messages_HandleMessagesClient
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
			log.Fatal(err)
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

func (c *Client) Connect(username string) chan bool {
	timeout := 10 * time.Second

	// Set up a connection to the server.
	conn, err := grpc.Dial("bloblet.com:4000", grpc.WithInsecure(), grpc.WithTimeout(timeout), grpc.WithBlock())
	c.conn = conn
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	a := pb.NewAuthClient(conn)
	// The channel is to make sure we don't try to do anything with a null SessionToken.
	updated := make(chan bool)
	go c.keepalive(a, username, updated)
	<-updated

	msgClient := pb.NewMessagesClient(c.conn)

	md := metadata.New(map[string]string{"session-token": c.token})

	messageStream, err := msgClient.HandleMessages(metadata.NewOutgoingContext(context.Background(), md))

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for true {
			msg, err := c.messageStream.Recv()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("<%v> %v\n", msg.UserID, msg.Content)
		}
	}()

	c.messageStream = messageStream

	return updated
}

func (c *Client) addMetadata(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"session-token": c.token}))
}

func (c *Client) SendMessage(message string) {
	message = strings.ReplaceAll(message, "\n", "")
	c.messageStream.Send(&pb.CreateMessage{Content: message})
}

func main() {
	c := Client{}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Pick a username: ")
	username, _ := reader.ReadString('\n')

	c.Connect(strings.ReplaceAll(username, "\n", ""))
	fmt.Println("Connected to Fenix")
	for true {
		text, _ := reader.ReadString('\n')
		fmt.Printf("%v%v", ansi.CPL, ansi.DL)
		c.SendMessage(text)
	}
}
