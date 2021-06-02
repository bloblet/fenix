package api

import (
	"context"
	pb "github.com/bloblet/fenix/protobufs/go"
	db "github.com/bloblet/fenix/server/databases"
	"github.com/bloblet/fenix/server/models"
	"github.com/bloblet/fenix/server/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"strconv"
)

var config = utils.LoadConfig()
var addr = config.API.Host + ":" + strconv.Itoa(config.API.Port)

type GRPCApi struct {
	S     *grpc.Server
	msgDB *db.MessageDB
	pb.UnimplementedUsersServer
	pb.UnimplementedMessagesServer
	authDB           *db.AuthenticationManager
	httpApi          HTTPApi
	connectedClients map[string]chan *pb.Message
}

func (api *GRPCApi) Prepare() {
	utils.Log().Trace("")
	api.S = grpc.NewServer()
	api.msgDB = db.NewMessageDB()
	api.authDB = db.NewAuthenticationManager()
	api.httpApi = HTTPApi{}
	api.connectedClients = map[string]chan *pb.Message{}

	pb.RegisterUsersServer(api.S, api)
	pb.RegisterMessagesServer(api.S, api)
}

func (api *GRPCApi) Bufconn() *bufconn.Listener {
	api.Prepare()
	b := bufconn.Listen(1024 * 1024)
	go api.Listen(b)

	return b
}

func (api *GRPCApi) Listen(lis net.Listener) {
	utils.Log().Trace("")
	go api.httpApi.ServeHTTP()
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
	utils.Log().Trace("")
	api.Prepare()
	utils.Log().Infof("Serving GRPC on %v", addr)
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

func (api *GRPCApi) authenticate(a *pb.AuthMethod) (*models.User, bool) {
	utils.Log().Trace("")

	return api.authDB.AuthenticateUser(a)
}

func (api *GRPCApi) RequestToken(_ context.Context, in *pb.AuthMethod) (*pb.AuthMethod, error) {
	utils.Log().Trace("")
	_, okid := api.authDB.GetUser(in.UserID)
	_, okem := api.authDB.GetUserByEmail(in.Password.Email)

	if !okid && !okem {
		return nil, UserDoesNotExistError{}
	}

	authenticated := false
	var u *models.User

	if okid {

		u, authenticated = api.authenticate(in)

		if authenticated {
			err := api.authDB.DeleteToken(u, in.Token.GetTokenID())
			if err != nil {
				return nil, err
			}

		}

	} else if okem {
		u, authenticated = api.authDB.PasswordAuthenticateUser(
			in.Password.GetEmail(),
			in.Password.GetPassword(),
		)
	} else {
		return nil, InvalidRequest{}
	}

	if !authenticated {
		return nil, NotAuthorized{}
	}

	token, err := api.authDB.CreateToken()
	if err != nil {
		return nil, err
	}

	err = api.authDB.AddToken(u, token)
	if err != nil {
		return nil, err
	}

	return &pb.AuthMethod{
		Token: token.MarshalToPB(),
		UserID: u.ID.String(),
	}, nil
}

func (api *GRPCApi) GetUser(_ context.Context, in *pb.RequestUser) (*pb.User, error) {
	utils.Log().Trace("")
	if _, a := api.authenticate(in.GetAuthentication()); !a {
		return nil, NotAuthorized{}
	}

	u, ok := api.authDB.GetUser(in.GetUserID())
	if !ok {
		return nil, UserDoesNotExistError{}
	}
	return u.MarshalToPB(), nil
}

func (api *GRPCApi) CreateUser(_ context.Context, in *pb.RequestUserCreation) (*pb.UserCreated, error) {
	utils.Log().Trace("")

	u, err := api.authDB.NewUser(in.GetEmail(), in.GetUsername(), in.GetPassword())
	utils.Log().Trace(err)

	if err != nil {
		return nil, err
	}

	var uc *pb.UserCreated
	for s, _ := range u.Tokens {
		uc = u.MarshalToUserCreated(s)
	}

	return uc, nil
}

func (api *GRPCApi) WaitForEmailVerification(in *pb.AuthMethod, s pb.Users_WaitForEmailVerificationServer) error {
	utils.Log().Trace("")
	u, authenticated := api.authenticate(in)
	if !authenticated {
		return NotAuthorized{}
	}

	if u.EmailVerified == true {
		return s.Send(&pb.Success{})
	}

	if _, ok := api.authDB.VerificationListeners[u.ID.Hex()]; ok {
		return InvalidRequest{}
	}

	api.authDB.VerificationListeners[u.ID.Hex()] = make(chan bool, 1)
	<-api.authDB.VerificationListeners[u.ID.Hex()]
	return s.Send(&pb.Success{})
}

func (api *GRPCApi) ResendEmailVerification(_ context.Context, in *pb.AuthMethod) (*pb.Success, error) {
	utils.Log().Trace("")
	u, authenticated := api.authenticate(in)
	if !authenticated {
		return nil, NotAuthorized{}
	}

	err := api.authDB.SendVerificationEmail(u)
	if err != nil {
		return nil, err
	}

	return &pb.Success{}, nil
}

func (api *GRPCApi) ChangeMFA(_ context.Context, in *pb.MFAStatus) (*pb.Success, error) {
	utils.Log().Trace("")
	u, authenticated := api.authenticate(in.GetAuthentication())
	if !authenticated {
		return nil, NotAuthorized{}
	}

	err := api.authDB.ChangeMFA(u, in.GetStatus())
	if err != nil {
		return nil, err
	}
	return &pb.Success{}, nil
}

func (api *GRPCApi) GetMFALink(_ context.Context, in *pb.RequestMFALink) (*pb.MFALink, error) {
	utils.Log().Trace("")
	u, authenticated := api.authenticate(in.GetAuthentication())
	if !authenticated {
		return nil, NotAuthorized{}
	}

	l, err := api.authDB.Generate2FALink(u)
	if err != nil {
		return nil, err
	}

	return &pb.MFALink{Link: l}, nil
}

func (api *GRPCApi) ChangeUsername(_ context.Context, in *pb.ChangeUsernameRequest) (*pb.User, error) {
	utils.Log().Trace("")
	u, authenticated := api.authenticate(in.GetAuthentication())
	if !authenticated {
		return nil, NotAuthorized{}
	}
	err := api.authDB.ChangeUsername(u, in.GetUsername())
	if err != nil {
		return nil, err
	}
	return u.MarshalToPB(), nil
}

func (api *GRPCApi) ChangePassword(_ context.Context, in *pb.ChangePasswordRequest) (*pb.UserCreated, error) {
	utils.Log().Trace("")
	u, authenticated := api.authDB.PasswordAuthenticateUser(
		in.GetAuthentication().GetEmail(),
		in.GetAuthentication().GetPassword(),
	)
	if !authenticated {
		return nil, NotAuthorized{}
	}

	t, err := api.authDB.ChangePassword(u, in.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.UserCreated{
		User:  u.MarshalToPB(),
		Token: t.MarshalToPB(),
	}, nil
}

func (api GRPCApi) GetMessageHistory(ctx context.Context, history *pb.RequestMessageHistory) (*pb.MessageHistory, error) {
	utils.Log().Trace("")
	user, authenticated := api.authenticate(history.GetAuthentication())

	if !authenticated {
		return nil, NotAuthorized{}
	}

	messageHistory, err := api.msgDB.FetchMessagesAfter(history.GetLastMessageTime().AsTime())

	if err != nil {
		utils.Log().Error(err)
		return nil, err
	}

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
	utils.Log().Trace("")
	m, err := stream.Recv()
	if err != nil {
		return err
	}
	user, authenticated := api.authenticate(m.GetAuthentication())

	if !authenticated {
		return NotAuthorized{}
	}

	api.connectedClients[user.ID.Hex()] = make(chan *pb.Message, 1)

	// Pass any sent messages to the client
	go func() {
		for true {
			_ = stream.Send(<-api.connectedClients[user.ID.Hex()])
		}
	}()

	p, _ := peer.FromContext(stream.Context())

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
	for _, c := range api.connectedClients {
		c <- message
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

type UserDoesNotExistError struct {
}

func (e UserDoesNotExistError) Error() string {
	return "UserDoesNotExist"
}

type InvalidRequest struct {
}

func (e InvalidRequest) Error() string {
	return "InvalidRequest"
}

type NotAuthorized struct {
}

func (e NotAuthorized) Error() string {
	return "NotAuthorized"
}
