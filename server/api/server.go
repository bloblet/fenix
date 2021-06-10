package api

import (
	"context"
	"fmt"
	pb "github.com/bloblet/fenix/protobufs/go"
	db "github.com/bloblet/fenix/server/databases"
	"github.com/bloblet/fenix/server/models"
	"github.com/bloblet/fenix/server/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func (api *GRPCApi) authenticate(a *pb.AuthMethod) (*models.User, bool, error) {
	u, auth, err := api.authDB.AuthenticateUser(a)
	if auth == false && err == nil {
		err = NotAuthorized{}
	}
	return u, auth, err
}

func (api *GRPCApi) RequestToken(_ context.Context, in *pb.AuthMethod) (*pb.AuthMethod, error) {
	_, okid := api.authDB.GetUser(in.UserID)
	_, okem := api.authDB.GetUserByEmail(in.Password.Email)

	if !okid && !okem {
		return nil, UserDoesNotExistError{}
	}

	authenticated := false
	var u *models.User
	if okid {

		u, authenticated, err := api.authenticate(in)
		if err != nil {
			return nil, err
		}
		if authenticated {
			err := api.authDB.DeleteToken(u, in.Token.GetTokenID())
			if err != nil {
				return nil, err
			}

		}

	} else if okem {
		u, authenticated, _ = api.authDB.PasswordAuthenticateUser(
			in.GetPassword().GetEmail(),
			in.GetPassword().GetPassword(),
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
		Token:  token.MarshalToPB(),
		UserID: u.ID.Hex(),
	}, nil
}

func (api *GRPCApi) GetUser(_ context.Context, in *pb.RequestUser) (*pb.User, error) {
	utils.Log().Trace("")
	if _, a, err := api.authenticate(in.GetAuthentication()); !a {
		return nil, NotAuthorized{Message: err.Error()}
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
	u, a, err := api.authenticate(in)
	if !a {
		return NotAuthorized{Message: err.Error()}
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
	u, a, err := api.authenticate(in)
	if !a {
		return nil, NotAuthorized{Message: err.Error()}
	}

	err = api.authDB.SendVerificationEmail(u)
	if err != nil {
		return nil, err
	}

	return &pb.Success{}, nil
}

func (api *GRPCApi) ChangeMFA(_ context.Context, in *pb.MFAStatus) (*pb.Success, error) {
	utils.Log().Trace("")
	u, a, err := api.authenticate(in.GetAuthentication())
	if !a {
		return nil, NotAuthorized{Message: err.Error()}
	}

	err = api.authDB.ChangeMFA(u, in.GetStatus())
	if err != nil {
		return nil, err
	}
	return &pb.Success{}, nil
}

func (api *GRPCApi) GetMFALink(_ context.Context, in *pb.RequestMFALink) (*pb.MFALink, error) {
	utils.Log().Trace("")
	u, a, err := api.authenticate(in.GetAuthentication())
	if !a {
		return nil, NotAuthorized{Message: err.Error()}
	}

	l, err := api.authDB.Generate2FALink(u)
	if err != nil {
		return nil, err
	}

	return &pb.MFALink{Link: l}, nil
}

func (api *GRPCApi) ChangeUsername(_ context.Context, in *pb.ChangeUsernameRequest) (*pb.User, error) {
	utils.Log().Trace("")
	u, a, err := api.authenticate(in.GetAuthentication())
	if !a {
		return nil, NotAuthorized{Message: err.Error()}
	}
	err = api.authDB.ChangeUsername(u, in.GetUsername())
	if err != nil {
		return nil, err
	}
	return u.MarshalToPB(), nil
}

func (api *GRPCApi) ChangePassword(_ context.Context, in *pb.ChangePasswordRequest) (*pb.UserCreated, error) {
	utils.Log().Trace("")
	u, authenticated, err := api.authDB.PasswordAuthenticateUser(
		in.GetAuthentication().GetEmail(),
		in.GetAuthentication().GetPassword(),
	)
	if !authenticated {
		return nil, NotAuthorized{Message: err.Error()}
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

func (api GRPCApi) GetMessageHistory(ctx context.Context, in *pb.RequestMessageHistory) (*pb.MessageHistory, error) {
	utils.Log().Trace("")
	u, a, err := api.authenticate(in.GetAuthentication())
	if !a {
		return nil, NotAuthorized{Message: err.Error()}
	}

	messageHistory, err := api.msgDB.FetchMessagesAfter(in.GetLastMessageTime().AsTime())

	if err != nil {
		utils.Log().Error(err)
		return nil, err
	}

	p, _ := peer.FromContext(ctx)
	utils.Log().WithFields(
		log.Fields{
			"userID":        u.ID,
			"username":      u.Username,
			"ip":            p.Addr,
			"numOfMessages": messageHistory.NumberOfMessages,
		},
	).Trace("getMessageHistory")

	return messageHistory, nil
}

func (api *GRPCApi) HandleMessages(stream pb.Messages_HandleMessagesServer) error {
	p, _ := peer.FromContext(stream.Context())
	meta, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return NotAuthorized{Message: "Invalid Metadata"}
	}

	_tokenID := meta["tokenid"]
	_userID := meta["userid"]
	_token := meta["token"]

	if len(_tokenID) != 1 || len(_userID) != 1 || len(_token) != 1 {
		return NotAuthorized{Message: "Invalid Metadata (Not exactly one of each metadata variable)"}
	}

	tokenID := _tokenID[0]
	userID := _userID[0]
	token := _token[0]

	if tokenID == "" || userID == "nil" || token == "" {
		return NotAuthorized{Message: "Either tokenID, userID, or token are nil."}
	}

	utils.Log().WithField("ip", p.Addr).Trace()

	utils.Log().WithFields(log.Fields{
		"userID":  userID,
		"token":   token,
		"tokenID": tokenID,
	}).Trace("HandleMessages Auth")

	user, authenticated, err := api.authDB.TokenAuthenticateUser(userID, token, tokenID)

	if err != nil {
		return NotAuthorized{err.Error()}
	}

	if !authenticated {
		utils.Log().Trace()
		return NotAuthorized{Message: "Invalid token/tokenid/userid"}
	}

	api.connectedClients[user.ID.Hex()] = make(chan *pb.Message, 1)

	// Pass any sent messages to the client
	go func() {
		for true {
			_ = stream.Send(<-api.connectedClients[user.ID.Hex()])
		}
	}()

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
		go api.authDB.RenewToken(token, tokenID, user)
		if err != nil {
			return err
		}

		msg := api.msgDB.NewMessage(sendRequest, user)

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
	Message string
}

func (e NotAuthorized) Error() string {
	return fmt.Sprintf("Not Authorized: %v", e.Message)
}
