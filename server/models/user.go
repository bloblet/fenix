package models

import (
	pb "github.com/bloblet/fenix/protobufs/go"
	"github.com/kamva/mgm/v3"
)

// User is the current datatype for fenix users.
type User struct {
	SyncModel        `bson:",inline"`
	mgm.DefaultModel `bson:",inline"`
	AuthSecret       string
	Tokens           map[string]Token
	Email            string
	Salt             []byte `json:"-"`
	Password         []byte `json:"-"`
	OTT              string
	EmailVerified    bool
	Username         string
	Discriminator    string
	MFAEnabled       bool
}

func (u *User) MarshalToPB() *pb.User {
	p := pb.User{}
	p.Discriminator = u.Discriminator
	p.UserID = u.ID.Hex()
	p.Username = u.Username
	return &p
}

func (u *User) MarshalToUserCreated(tokenID string) *pb.UserCreated {
	p := pb.UserCreated{}
	p.User = u.MarshalToPB()
	t := u.Tokens[tokenID]
	p.Token = t.MarshalToPB()
	return &p
}

func (u *User) CollectionName() string {
	return "users"
}
