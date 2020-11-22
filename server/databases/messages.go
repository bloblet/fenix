package databases

import (
	"fmt"
	pb "github.com/bloblet/fenix-protobufs/go"
	"github.com/gocql/gocql"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const (
	MessageID = "MessageID"
	Content   = "Content"
	SentAt = "CreatedAt"
	UserID = "UserID"
)

func NewMessagesSession(config *gocql.ClusterConfig) (*gocql.Session, error) {
	config.Keyspace = "messages"

	return config.CreateSession()
}

func NewMessageDB() *MessageDB {
	db := MessageDB{}
	db.MessageCache = make(map[string]*pb.Message)
	db.MessageCacheExpiry = make(map[string]time.Time)
	go db.PurgeCache()
	return &db
}

type MessageDB struct {
	// TODO: Change to a channel cache when channels are implemented
	MessageCache       map[string]*pb.Message
	MessageCacheExpiry map[string]time.Time
}

func (db *MessageDB) PurgeCache() {
	timer := time.NewTimer(5 * time.Second)
	for {
		<-timer.C
		n := time.Now()
		for id, t := range db.MessageCacheExpiry {
			// Fixes a linter warning, not sure why
			t := t
			id := id
			go func() {
				if t.After(n) {
					delete(db.MessageCache, id)
				}
			}()
		}
		timer.Reset(5 * time.Second)
	}
}

func (db *MessageDB) NewMessage(s *gocql.Session, message *pb.CreateMessage, userID string) *pb.Message {
	id := MakeULID()

	q := s.Query("INSERT INTO Messages (MessageID, UserID, Content, CreatedAt) VALUES (?, ?, ?, ?)", id.String(), userID, message.GetContent(), time.Now())
	err := q.Exec()

	if err != nil {
		panic(err)
	}
	m := pb.Message{}
	m.Content = message.GetContent()
	m.MessageID = id.String()
	m.SentAt = timestamppb.New(ulid.Time(id.Time()))
	m.UserID = userID
	db.MessageCacheExpiry[m.MessageID] = time.Now().Add(5 * time.Second)
	return &m
}

func (db *MessageDB) GetMessage(id string) *pb.Message {
	return db.MessageCache[id]
}

func (db MessageDB) FetchMessage(s *gocql.Session, id string) *pb.Message {
	rawMessage := make(map[string]interface{})

	err := s.Query("SELECT * FROM Messages WHERE MessageID = ?", id).Scan(&rawMessage)

	if err != nil {
		fmt.Printf("Error fetching message, %v", err)
		return nil
	}

	message := pb.Message{}
	message.MessageID = rawMessage[MessageID].(string)
	message.Content = rawMessage[Content].(string)

	return &message
}

func (db MessageDB) MaybeGetMessage(s func() *gocql.Session, id string) *pb.Message {
	if msg := db.GetMessage(id); msg != nil {
		return msg
	}
	return db.FetchMessage(s(), id)
}

// TODO: Fix DB design issue
func (db MessageDB) FetchMessagesBefore(s *gocql.Session, t time.Time) *pb.MessageHistory {
	messages, err := s.Query("SELECT * FROM Messages WHERE CreatedAt <= ? ORDER BY MessageID LIMIT 50", t).Iter().SliceMap()
	if err != nil {
		fmt.Printf("Error getting message history, %v", err)
		return nil
	}

	msgHistory := make([]*pb.Message, 0)

	for i := 0; i < len(messages); i++ {
		msg := messages[i]
		msgHistory = append(msgHistory, &pb.Message{
			MessageID: msg[MessageID].(string),
			UserID:    msg[UserID].(string),
			SentAt:    timestamppb.New(msg[SentAt].(time.Time)),
			Content:   msg[Content].(string),
		})
	}
	return &pb.MessageHistory{Messages: msgHistory}
}
