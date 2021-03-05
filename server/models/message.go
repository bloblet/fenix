package models

import (
	pb "github.com/bloblet/fenix/protobufs/go"
	"github.com/kamva/mgm/v3"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Message struct {
	onSave           chan bool `bson:"-"`
	mgm.DefaultModel `bson:",inline"`
	UserID           string
	CreatedAt        time.Time
	ChannelID        string
	Content          string
}

func (m *Message) SetupMessage() {
	m.onSave = make(chan bool, 1)
}

func (m *Message) MarshalToPB() *pb.Message {
	message := pb.Message{}
	message.Content = m.Content
	message.UserID = m.UserID
	message.MessageID = m.ID.Hex()
	message.SentAt = timestamppb.New(m.CreatedAt)
	return &message
}

func (m *Message) Saved() error {
	m.onSave <- true
	return nil
}

func (m *Message) WaitForSave() {
	<-m.onSave
}

func (m *Message) CollectionName() string {
	return "messages"
}
