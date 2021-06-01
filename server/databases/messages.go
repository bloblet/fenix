package databases

import (
	pb "github.com/bloblet/fenix/protobufs/go"
	"github.com/bloblet/fenix/server/models"
	"github.com/bloblet/fenix/server/utils"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var config = utils.LoadConfig()

func NewMessageDB() *MessageDB {
	db := MessageDB{}
	db.MessageCache = make(map[string]*pb.Message)
	db.MessageCacheExpiry = make(map[string]time.Time)
	err := mgm.SetDefaultConfig(nil, config.Database.Database, options.Client().ApplyURI(config.Database.URI))

	if err != nil {
		utils.Log().WithFields(
			log.Fields{
				"host":     config.Database.URI,
				"database": config.Database.Database,
				"error":    err,
			},
		).Panic("Error connecting to mongodb")
	}

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

func (db *MessageDB) NewMessage(cMsg *pb.CreateMessage, userID string, sync ...bool) *pb.Message {
	_sync := false

	if len(sync) != 0 {
		_sync = sync[0]
	}

	msg := &models.Message{
		UserID:    userID,
		Content:   cMsg.Content,
		CreatedAt: time.Now(),
		ChannelID: "0",
	}

	msg.New()

	err := mgm.Coll(msg).Create(msg)

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

	if _sync {
		msg.WaitForSave()
	}

	db.MessageCacheExpiry[m.MessageID] = time.Now().Add(5 * time.Second)
	return m
}

func (db *MessageDB) GetMessage(id string) *pb.Message {
	return db.MessageCache[id]
}

func (db MessageDB) FetchMessage(id string) (*pb.Message, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, InvalidID{}
	}

	msg := &models.Message{}

	err := mgm.Coll(msg).FindByID(id, msg)

	if err != nil {
		return nil, err
	}

	return msg.MarshalToPB(), nil
}

func (db MessageDB) MaybeGetMessage(id string) (*pb.Message, error) {
	if msg := db.GetMessage(id); msg != nil {
		return msg, nil
	}
	return db.FetchMessage(id)
}

func (db MessageDB) FetchMessagesAfter(t time.Time) (*pb.MessageHistory, error) {
	msgs := make([]*models.Message, 0)

	err := mgm.Coll(&models.Message{}).SimpleFind(
		&msgs,
		bson.M{"channelid": bson.M{"$eq": "0"}, "createdat": bson.M{"$gt": t}},
		options.Find().SetLimit(50),
		options.Find().SetSort(bson.M{"createdat": -1}),
	)

	if err != nil {
		return nil, err
	}

	msgHistory := make([]*pb.Message, 0)

	for _, msg := range msgs {
		msgHistory = append(msgHistory, msg.MarshalToPB())
	}

	history := &pb.MessageHistory{
		Messages:         msgHistory,
		NumberOfMessages: int64(len(msgHistory)),
		Pages:            1,
	}

	return history, nil
}

type InvalidID struct {
}

func (id InvalidID) Error() string {
	return "InvalidID"
}
