package database

import (
	pb "github.com/bloblet/fenix/protobufs/go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Message struct {
	MessageID string
	UserID    string
	CreatedAt time.Time
	ChannelID string
	//ServerID  string
	Content   string
}

func (m Message) MarshalToPB() *pb.Message {
	message := pb.Message{}

	message.Content = m.Content
	message.UserID = m.UserID
	message.MessageID = m.MessageID
	message.SentAt = timestamppb.New(m.CreatedAt)
	return &message
}