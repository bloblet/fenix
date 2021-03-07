package models

import (
	pb "github.com/bloblet/fenix/protobufs/go"
	"github.com/kamva/mgm/v3"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Message struct {
	SyncModel `bson:",inline"`
	mgm.DefaultModel `bson:",inline"`
	UserID           string
	CreatedAt        time.Time
	ChannelID        string
	Content          string
}

func (m *Message) MarshalToPB() *pb.Message {
	message := pb.Message{}
	message.Content = m.Content
	message.UserID = m.UserID
	message.MessageID = m.ID.Hex()
	message.SentAt = timestamppb.New(m.CreatedAt)
	return &message
}

func (m *Message) CollectionName() string {
	return "messages"
}
