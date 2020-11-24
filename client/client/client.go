package client

import (
	"context"
	"log"

	"time"

	pb "github.com/bloblet/fenix-protobufs/go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var timeout = 10 * time.Second

type Client struct {
	token         string
	Username      string
	conn          *grpc.ClientConn
	messageStream pb.Messages_HandleMessagesClient
	Debug         bool
	SessionTokens chan *pb.AuthAck
	Messages      chan *pb.Message
	msgClient     pb.MessagesClient
	LastMessageID string
}

func (c *Client) keepalive(a pb.AuthClient, username string, sessionTokens chan *pb.AuthAck) {
	// Timeout of 10 seconds

	// Keepalive loop
	for true {
		// Send login request
		// Timeout of 10 seconds
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		loginAck, err := a.Login(ctx, &pb.ClientAuth{Username: username})
		cancel()

		if err != nil {
			log.Fatal(err)
		}
		// Update values
		c.token = loginAck.GetSessionToken()
		c.Username = loginAck.GetUsername()
		// Send updated signal on channel
		sessionTokens <- loginAck

		// Wait until 30 seconds before token expires, for no "jank"
		time.Sleep(loginAck.GetExpiry().AsTime().Sub(time.Now()))
	}
}

func (c *Client) auth(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{"session-token": c.token})
	return metadata.NewOutgoingContext(ctx, md)
}

func (c *Client) dial(addr string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(timeout), grpc.WithBlock())
	c.conn = conn
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
}

func (c *Client) initAuthClient(username string) {
	a := pb.NewAuthClient(c.conn)

	// The channel is to make sure we don't try to do anything with a null SessionToken.
	c.SessionTokens = make(chan *pb.AuthAck)

	go c.keepalive(a, username, c.SessionTokens)
	<-c.SessionTokens
}

func (c *Client) initMessageClient() {
	c.msgClient = pb.NewMessagesClient(c.conn)

	messageStream, err := c.msgClient.HandleMessages(c.auth(context.Background()))

	if err != nil {
		log.Fatal(err)
	}
	c.messageStream = messageStream
	c.Messages = make(chan *pb.Message)

	go func() {
		for true {
			msg, err := c.messageStream.Recv()
			if err != nil {
				log.Fatal(err)
			}
			c.LastMessageID = msg.MessageID
			c.Messages <- msg
		}
	}()
}

func (c *Client) Connect(username string, addr string) {
	c.dial(addr)

	c.initAuthClient(username)

	c.initMessageClient()
}

// TODO: Uncomment once DB design issue is fixed
//func (c *Client) RequestMessageHistory(lastMessageTime time.Time) []*pb.Message {
//	history, err := c.msgClient.GetMessageHistory(c.auth(context.Background()), &pb.RequestMessageHistory{LastMessageTime: lastMessageTime})
//	if err != nil {
//		panic(err)
//	}
//
//	return history.GetMessages()
//}

func (c *Client) SendMessage(message string) {
	c.messageStream.Send(&pb.CreateMessage{Content: message})
}
