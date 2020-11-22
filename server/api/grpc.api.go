package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"net"
	"time"

	pb "github.com/bloblet/fenix-protobufs/go"
	db "github.com/bloblet/fenix/server/databases"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func generateToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	return base64.URLEncoding.EncodeToString(b), err
}

type user struct {
	ID       string
	Username string
	messages chan *pb.Message
}

type GRPCApi struct {
	S             *grpc.Server
	sessions      map[string]user
	msgDB         *db.MessageDB
	clusterConfig *gocql.ClusterConfig
	pb.UnimplementedAuthServer
	pb.UnimplementedMessagesServer
}

func (api *GRPCApi) Prepare() {
	api.S = grpc.NewServer()
	api.msgDB = db.NewMessageDB()
	api.clusterConfig = gocql.NewCluster("localhost")
	api.clusterConfig.Consistency = gocql.Consistency(1)
	api.clusterConfig.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "cassandra",
	}

	api.sessions = make(map[string]user)
	pb.RegisterAuthServer(api.S, api)
	pb.RegisterMessagesServer(api.S, api)
}

func (api *GRPCApi) Listen(lis net.Listener) {
	if err := api.S.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (api *GRPCApi) Serve() {
	api.Prepare()

	lis, err := net.Listen("tcp", "0.0.0.0:4000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	api.Listen(lis)
}

func (api GRPCApi) utilCreateMessageSession() *gocql.Session {
	session, err := db.NewMessagesSession(api.clusterConfig)
	if err != nil {
		fmt.Printf("Error starting DB session, %v", err)
		panic(err)
	}
	return session
}

// utilCheckSessionToken is a helper function that can validate and identify a request.
// If clients have more than one session-token, fenix only uses the first one.
func (api *GRPCApi) utilCheckSessionToken(ctx context.Context) user {
	md, _ := metadata.FromIncomingContext(ctx)
	token := md.Get("session-token")[0]
	return api.sessions[token]
}

// gRPC doesn't have any way of identifying clients, other than client metadata.
// To avoid cluttering all the protobuf requests with token parameters, and to avoid messy bidirectional stream workarounds,
// Fenix uses session tokens in metadata.  Clients are expected to log in and then keep that session token in metadata, and renew
// when it expires.  If anyone has a better solution, open an issue.
func (api *GRPCApi) Login(_ context.Context, in *pb.ClientAuth) (*pb.AuthAck, error) {
	sessionToken, err := generateToken(16)
	if err != nil {
		return nil, err
	}
	user := user{}
	user.messages = make(chan *pb.Message)
	user.Username = in.GetUsername()

	api.sessions[sessionToken] = user

	go func() {
		timer := time.NewTimer(5 * time.Minute)
		<-timer.C
		delete(api.sessions, sessionToken)
	}()

	return &pb.AuthAck{
		Username:     user.Username,
		SessionToken: sessionToken,
		Expiry:       timestamppb.New(time.Now().Add(5 * time.Minute)),
	}, nil
}

// TODO: Fix DB design issue
//func (api GRPCApi) GetMessageHistory(ctx context.Context, history *pb.RequestMessageHistory) (*pb.MessageHistory, error) {
//	api.utilCheckSessionToken(ctx)
//	session := api.utilCreateMessageSession()
//	msg := api.msgDB.MaybeGetMessage(func() *gocql.Session {
//		return session
//	}, history.GetLastMessageID())
//
//	return api.msgDB.FetchMessagesAfter(session, msg.GetSentAt().AsTime()), nil
//}

func (api *GRPCApi) HandleMessages(stream pb.Messages_HandleMessagesServer) error {
	user := api.utilCheckSessionToken(stream.Context())
	if user.Username == "" {
		return grpc.ErrClientConnClosing
	}

	// Pass any sent messages to the client
	go func() {
		for true {
			_ = stream.Send(<-user.messages)
		}
	}()

	session, err := db.NewMessagesSession(api.clusterConfig)
	if err != nil {
		fmt.Printf("Error starting DB session, %v", err)
		panic(err)
	}

	// Send messages the client requests
	for true {
		// Wait for the next message request
		sendRequest, err := stream.Recv()
		if err != nil {
			return grpc.ErrClientConnClosing
		}

		msg := api.msgDB.NewMessage(session, sendRequest, user.Username)
		// Notify all clients of the message
		api.notifyClientsOfMessage(msg)
	}
	return grpc.ErrClientConnClosing
}

func (api *GRPCApi) notifyClientsOfMessage(message *pb.Message) {
	for _, val := range api.sessions {
		val.messages <- message
	}
}
