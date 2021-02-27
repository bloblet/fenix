package models

import (
	pb "github.com/bloblet/fenix/protobufs/go"
	"github.com/go-bongo/bongo"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Message struct {
	Saved chan bool `bson:"-"`
	bongo.DocumentBase `bson:",inline"`
	UserID             string
	CreatedAt          time.Time
	ChannelID          string
	Content            string
}

func (m *Message) SetupMessage() {
	m.Saved = make(chan bool, 1)
}

func (m *Message) MarshalToPB() *pb.Message {
	message := pb.Message{}
	message.Content = m.Content
	message.UserID = m.UserID
	message.MessageID = m.Id.Hex()
	message.SentAt = timestamppb.New(m.CreatedAt)
	return &message
}

func (m *Message) AfterSave(b *bongo.Collection) error {
	m.Saved <- true
	return nil
}
