package api

import (
	"context"
	pb "github.com/bloblet/fenix/proto/6.0.1"
	"log"
	"net"

	"google.golang.org/grpc"
)

func NewGRPCApi() *GRPCApi {
	api := GRPCApi{}
	api.c = make(chan interface{}) // TODO: Change to pb.Message class
	return &api
}

type GRPCApi struct {
	s *grpc.Server
	c chan interface{} // TODO: Change to pb.Message class
}
func (api *GRPCApi) Serve() {
	api.s = grpc.NewServer()
	pb.RegisterUsersService(api.s, &pb.UsersService{Get: api.get})

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

func (api *GRPCApi) notifyClientsOfMessage(message interface{}) { // TODO: Change to pb.Message class
	api.c <- message
}