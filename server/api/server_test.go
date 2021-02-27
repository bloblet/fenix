package api

import (
	"context"
	"github.com/bloblet/fenix/client/client"
	pb "github.com/bloblet/fenix/protobufs/go"
	"github.com/bloblet/fenix/server/models"
	"github.com/bloblet/fenix/server/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"testing"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func makeString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func startBufconn() (*GRPCApi, *bufconn.Listener) {
	grpcAPI := GRPCApi{}
	return &grpcAPI, grpcAPI.Bufconn()
}

func login(username string, t *testing.T) (*GRPCApi, *client.Client, *pb.AuthClient, *pb.AuthAck) {
	api, conn := startBufconn()
	cli := client.Client{}
	cli.BuffConnect("testApiLogin", conn, false)

	a := pb.NewAuthClient(cli.Conn)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*4)

	clientAuth := pb.ClientAuth{
		Username: "testApiLogin",
	}

	ack, err := a.Login(ctx, &clientAuth)

	if err != nil {
		t.Error(err)
	}

	return api, &cli, &a, ack
}

func connectToMessages(username string, t *testing.T) (*GRPCApi, *client.Client, pb.MessagesClient, pb.Messages_HandleMessagesClient, *pb.AuthAck) {
	api, cli, _, ack := login(username, t)

	token := ack.SessionToken

	m := pb.NewMessagesClient(cli.Conn)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*4)

	md := metadata.New(map[string]string{"session-token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	msgs, err := m.HandleMessages(ctx)

	if err != nil {
		t.Error(err)
	}
	return api, cli, m, msgs, ack
}
func TestGRPCApi_Login(t *testing.T) {
	_, _, _, ack := login("testApiLogin", t)

	if testing.Verbose() {
		utils.Log().WithFields(
			log.Fields{
				"username":     ack.Username,
				"expiry":       ack.Expiry.AsTime().Format(time.RFC3339),
				"sessionToken": ack.SessionToken,
			},
		).Info("TestAPILogin")
	}
}

func TestGRPCApi_HandleMessages(t *testing.T) {
	api, _, _, msgs, _ := connectToMessages("testApiHandleMessages", t)

	msgChan := make(chan *pb.Message)

	go func() {
		msg, err := msgs.Recv()

		if err != nil {
			t.Error(err)
		}
		msgChan <- msg
	}()

	err := msgs.Send(&pb.CreateMessage{
		Content: makeString(50),
	})

	if err != nil {
		t.Fatal(err)
	}

	go func() {
		time.Sleep(time.Second * 10)
		t.Fatal("Timeout on message receiving")
	}()

	msg := <-msgChan

	if testing.Verbose() {
		utils.Log().WithFields(
			log.Fields{
				"content":   msg.Content,
				"sentAt":    msg.SentAt.AsTime().Format(time.RFC3339),
				"messageID": msg.MessageID,
				"userID":    msg.UserID,
			},
		).Info("TestApiHandleMessages")
	}

	_, err = api.msgDB.FetchMessage(msg.GetMessageID())

	if err != nil {
		t.Error(err)
	}
}

func TestGRPCApi_GetMessageHistory(t *testing.T) {
	api, _, m, _, ack := connectToMessages("testApiGetMessageHistory", t)

	// Populate Database
	before := time.Now()

	for i := 0; i <= 50; i++ {
		msg := &models.Message{
			UserID:    ack.Username,
			Content:   makeString(50),
			CreatedAt: time.Now(),
			ChannelID: "0",
		}
		err := api.msgDB.Conn.Collection("messages").Save(msg)
		if err != nil {
			t.Fatal(err)
		}
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*4)

	md := metadata.New(map[string]string{"session-token": ack.SessionToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &pb.RequestMessageHistory{
		LastMessageTime: timestamppb.New(before),
	}

	time.Sleep(1000)
	history, err := m.GetMessageHistory(ctx, req)

	if err != nil {
		t.Fatal(err)
	}

	if history.NumberOfMessages != 50 {
		t.Errorf("Invalid number of messages after %v: %v", before.Format(time.RFC3339), history.NumberOfMessages)
	}
}
