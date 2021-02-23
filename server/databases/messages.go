package databases

import (
	pb "github.com/bloblet/fenix/protobufs/go"
	"github.com/bloblet/fenix/server/models/database"
	"github.com/bloblet/fenix/server/utils"
	"github.com/go-bongo/bongo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var config = utils.LoadConfig()

func NewMessageDB() *MessageDB {
	db := MessageDB{}
	db.MessageCache = make(map[string]*pb.Message)
	db.MessageCacheExpiry = make(map[string]time.Time)
	conn, err := bongo.Connect(&bongo.Config{ConnectionString: config.Database.Host, Database: config.Database.Database})

	if err != nil {
		utils.Log().WithFields(
			log.Fields{
				"host":     config.Database.Host,
				"database": config.Database.Database,
				"error":    err,
			},
		).Panic("Error connecting to mongodb")
	}
	db.conn = conn

	go db.PurgeCache()
	return &db
}

type MessageDB struct {
	// TODO: Change to a channel cache when channels are implemented
	MessageCache       map[string]*pb.Message
	MessageCacheExpiry map[string]time.Time
	conn               *bongo.Connection
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

func (db *MessageDB) NewMessage(cMsg *pb.CreateMessage, userID string) *pb.Message {

	msg := &database.Message{
		UserID:    userID,
		Content:   cMsg.Content,
		CreatedAt: time.Now(),
		ChannelID: "0",
	}
	err := db.conn.Collection("messages").Save(msg)

	if err != nil {
		utils.Log().WithFields(
			log.Fields{
				"userID":        userID,
				"contentLength": len(msg.Content),
				"createdAt":     msg.CreatedAt,
				"channelID":     msg.ChannelID,
				"err":           err,
			},
		).Error("Error creating message")
		return nil
	}

	m := msg.MarshalToPB()

	db.MessageCacheExpiry[m.MessageID] = time.Now().Add(5 * time.Second)
	return m
}

func (db *MessageDB) GetMessage(id string) *pb.Message {
	return db.MessageCache[id]
}

func (db MessageDB) FetchMessage(id string) *pb.Message {
	msg := &database.Message{}
	err := db.conn.Collection("messages").FindById(bson.ObjectId(id), msg)

	if err != nil {
		utils.Log().WithFields(
			log.Fields{
				"messageID": id,
				"err":       err,
			},
		).Error("Error fetching message")
	}

	return msg.MarshalToPB()
}

func (db MessageDB) MaybeGetMessage(id string) *pb.Message {
	if msg := db.GetMessage(id); msg != nil {
		return msg
	}
	return db.FetchMessage(id)
}

func (db MessageDB) FetchMessagesAfter(t time.Time) *pb.MessageHistory {
	resultSet := db.conn.Collection("messages").Find(bson.D{{"channelid", "0"}, {"createdat", bson.D{{"$gt", t}}}})

	results := resultSet.Query.Sort("createdAt").Limit(50).Iter()

	msgHistory := make([]*pb.Message, 0)

	var result database.Message
	for results.Next(&result) {
		msgHistory = append(msgHistory, result.MarshalToPB())
	}

	return &pb.MessageHistory{Messages: msgHistory}
}
