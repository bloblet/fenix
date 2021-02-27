package client

import (
	"context"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"time"

	pb "github.com/bloblet/fenix/protobufs/go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var timeout = 10 * time.Second

type Client struct {
	token         string
	Username      string
	Conn          *grpc.ClientConn
	MessageStream pb.Messages_HandleMessagesClient
	Debug         bool
	SessionTokens chan *pb.AuthAck
	Messages      chan *pb.Message
	MsgClient     pb.MessagesClient
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
			panic(err)
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
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	c.Conn = conn

	if err != nil {
		panic(err)
	}
}

func (c *Client) bufdial(lis *bufconn.Listener)  {
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithInsecure(), grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return lis.Dial()
	}))
	c.Conn = conn

	if err != nil {
		panic(err)
	}
}

func (c *Client) initAuthClient(username string) {
	a := pb.NewAuthClient(c.Conn)

	// The channel is to make sure we don't try to do anything with a null SessionToken.
	c.SessionTokens = make(chan *pb.AuthAck)

	go c.keepalive(a, username, c.SessionTokens)
	<-c.SessionTokens
}

func (c *Client) initMessageClient() {
	c.MsgClient = pb.NewMessagesClient(c.Conn)

	messageStream, err := c.MsgClient.HandleMessages(c.auth(context.Background()))

	if err != nil {
		panic(err)
	}
	c.MessageStream = messageStream
	c.Messages = make(chan *pb.Message)

	go func() {
		for true {
			msg, err := c.MessageStream.Recv()
			if err != nil {
				panic(err)
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

func (c *Client) BuffConnect(username string, lis *bufconn.Listener, setup bool) {
	c.bufdial(lis)

	if setup {
		c.initAuthClient(username)

		c.initMessageClient()
	}
}

func (c *Client) RequestMessageHistory(lastMessageTime time.Time) []*pb.Message {
	history, err := c.MsgClient.GetMessageHistory(c.auth(context.Background()), &pb.RequestMessageHistory{LastMessageTime: timestamppb.New(time.Now().Add(-time.Hour))})
	if err != nil {
		panic(err)
	}

	return history.GetMessages()
}

func (c *Client) SendMessage(message string) error {
	return c.MessageStream.Send(&pb.CreateMessage{Content: message})
}
