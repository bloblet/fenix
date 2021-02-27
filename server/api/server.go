package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	pb "github.com/bloblet/fenix/protobufs/go"
	db "github.com/bloblet/fenix/server/databases"
	"github.com/bloblet/fenix/server/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"strconv"
	"time"
)

var config = utils.LoadConfig()
var addr = config.API.Host + ":" + strconv.Itoa(config.API.Port)

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
	S        *grpc.Server
	sessions map[string]user
	msgDB    *db.MessageDB
	pb.UnimplementedAuthServer
	pb.UnimplementedMessagesServer
}

func (api *GRPCApi) Prepare() {
	api.S = grpc.NewServer()
	api.msgDB = db.NewMessageDB()
	api.sessions = make(map[string]user)
	pb.RegisterAuthServer(api.S, api)
	pb.RegisterMessagesServer(api.S, api)
}

func (api *GRPCApi) Bufconn() *bufconn.Listener {
	api.Prepare()
	b := bufconn.Listen(1024 * 1024)
	go api.Listen(b)

	return b
}

func (api *GRPCApi) Listen(lis net.Listener) {
	if err := api.S.Serve(lis); err != nil {
		utils.Log().WithFields(
			log.Fields{
				"addr": addr,
				"err":  err,
			},
		).Panic("Failed to serve API")
	}
}

func (api *GRPCApi) Serve() {
	api.Prepare()
	utils.Log().Infof("Serving on %v", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		utils.Log().WithFields(
			log.Fields{
				"addr": addr,
				"err":  err,
			},
		).Panic("Failed to listen on address")
	}
	api.Listen(lis)
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
func (api *GRPCApi) Login(ctx context.Context, in *pb.ClientAuth) (*pb.AuthAck, error) {
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
	p, _ := peer.FromContext(ctx)

	utils.Log().WithFields(
		log.Fields{
			"userID":   user.ID,
			"username": user.Username,
			"ip":       p.Addr,
		},
	).Trace("login")

	return &pb.AuthAck{
		Username:     user.Username,
		SessionToken: sessionToken,
		Expiry:       timestamppb.New(time.Now().Add(5 * time.Minute)),
	}, nil
}

func (api GRPCApi) GetMessageHistory(ctx context.Context, history *pb.RequestMessageHistory) (*pb.MessageHistory, error) {
	user := api.utilCheckSessionToken(ctx)
	messageHistory := api.msgDB.FetchMessagesAfter(history.GetLastMessageTime().AsTime())

	p, _ := peer.FromContext(ctx)
	utils.Log().WithFields(
		log.Fields{
			"userID":        user.ID,
			"username":      user.Username,
			"ip":            p.Addr,
			"numOfMessages": messageHistory.NumberOfMessages,
		},
	).Trace("getMessageHistory")

	return messageHistory, nil
}

func (api *GRPCApi) HandleMessages(stream pb.Messages_HandleMessagesServer) error {
	ctx := stream.Context()

	user := api.utilCheckSessionToken(ctx)

	if user.Username == "" {
		return InvalidUsername{}
	}

	// Pass any sent messages to the client
	go func() {
		for true {
			_ = stream.Send(<-user.messages)
		}
	}()

	p, _ := peer.FromContext(ctx)

	utils.Log().WithFields(
		log.Fields{
			"userID":   user.ID,
			"username": user.Username,
			"ip":       p.Addr,
		},
	).Trace("messageStream")

	// Send messages the client requests
	for true {
		// Wait for the next message request
		sendRequest, err := stream.Recv()
		if err != nil {
			return err
		}

		msg := api.msgDB.NewMessage(sendRequest, user.Username)

		utils.Log().WithFields(
			log.Fields{
				"userID":        user.ID,
				"contentLength": len(sendRequest.GetContent()),
				"messageID":     msg.GetMessageID(),
			},
		).Trace("createMessage")

		// Notify all clients of the message
		api.notifyClientsOfMessage(msg)
	}
	return ConnectionClosed{}
}

func (api *GRPCApi) notifyClientsOfMessage(message *pb.Message) {
	for _, val := range api.sessions {
		val.messages <- message
	}
}

type InvalidUsername struct{}

func (e InvalidUsername) Error() string {
	return "InvalidUsername"
}

type ConnectionClosed struct{}

func (e ConnectionClosed) Error() string {
	return "ConnectionClosed"
}
