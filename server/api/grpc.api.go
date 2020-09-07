package api

import (
	"context"
	pb "fenix/proto/6.0.1"
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCApi struct {
	s *grpc.Server
}

func (api *GRPCApi) Serve() {
	lis, err := net.Listen("tcp", "0.0.0.0:4000")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	api.s = grpc.NewServer()
	pb.RegisterUsersService(api.s, &pb.UsersService{Get: api.Get})

	if err := api.s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
func (api *GRPCApi) Get(ctx context.Context, in *pb.Authenticate) (*pb.User, error) {
	log.Printf("Received: %v", in.GetID())
	return &pb.User{ID: in.GetID() + in.GetToken()}, nil
}
