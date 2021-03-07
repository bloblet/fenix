package models

import (
	pb "github.com/bloblet/fenix/protobufs/go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Token struct {
	Token string
	Expires time.Time
	TokenID string
}
func (t *Token) MarshalToPB() *pb.Token {
	p := pb.Token{}
	p.Token = t.Token
	p.ExpirationDate = timestamppb.New(t.Expires)
	p.TokenID = t.TokenID
	return &p
}