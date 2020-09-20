package api

import (
	"context"
	pb "github.com/bloblet/fenix/proto/6.0.1"
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCApi struct {
	s *grpc.Server
}

func (api *GRPCApi) Serve() {
	api.s = grpc.NewServer()
	pb.RegisterUsersService(api.s, &pb.UsersService{Get: api.get})
	pb.RegisterAuthService(api.s, &pb.AuthService{Connect: api.connect, Login: api.login})

	lis, err := net.Listen("tcp", "0.0.0.0:4000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	if err := api.s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
func (api *GRPCApi) get(_ context.Context, in *pb.Authenticate) (*pb.User, error) {
	log.Printf("Received: %v", in.GetID())
	return &pb.User{ID: in.GetID() + in.GetToken()}, nil
}

func (api *GRPCApi) connect(_ context.Context, _ *pb.OpenHandshake) (*pb.FetchClientAuth, error) {
	return &pb.FetchClientAuth{}, nil
}

func (api *GRPCApi) login(_ context.Context, in *pb.ClientAuth) (*pb.AuthAck, error) {
	return &pb.AuthAck{Username: in.GetUsername()}, nil
}
