package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	mrand "math/rand"
	"net"
	"strconv"
	"time"

	pb "github.com/bloblet/fenix/proto/6.0.1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func generateToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	return base64.URLEncoding.EncodeToString(b), err
}

type GRPCApi struct {
	s        *grpc.Server
	sessions map[string]string
}

func (api *GRPCApi) Serve() {
	api.s = grpc.NewServer()
	api.sessions = make(map[string]string)
	
	pb.RegisterAuthService(api.s, &pb.AuthService{Login: api.login})

	lis, err := net.Listen("tcp", "0.0.0.0:4000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := api.s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// gRPC doesn't have any way of identifying clients, other than client metadata.
// To avoid cluttering all the protobuf requests with token parameters, and to avoid messy bidirectional stream workarounds,
// Fenix uses session tokens in metadata.  Clients are expected to log in and then keep that session token in metadata, and renew 
// when it expires.  If anyone has a better solution, open an issue.
func (api *GRPCApi) login(_ context.Context, in *pb.ClientAuth) (*pb.AuthAck, error) {
	sessionToken, err := generateToken(16)
	if err != nil {
		return nil, err
	}
	psudeoUniqueUsername := in.GetUsername() + strconv.Itoa(mrand.Intn(1000))
	api.sessions[sessionToken] = psudeoUniqueUsername

	defer func() {
		time.NewTimer(5 * time.Minute)
		delete(api.sessions, sessionToken)
	}()

	return &pb.AuthAck{
		Username: psudeoUniqueUsername, 
		SessionToken: sessionToken, 
		Expiry: timestamppb.New(time.Now().Add(5 * time.Minute)),
		}, nil
}
