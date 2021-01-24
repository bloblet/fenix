package databases

import (
	"fmt"
	pb "github.com/bloblet/fenix/protobufs/go"
	"github.com/bloblet/fenix/server/models/database"
	"github.com/hailocab/gocassa"
	"github.com/oklog/ulid/v2"
	"time"
)

func NewMessagesSession(hosts []string, username string, password string) (gocassa.KeySpace, error) {
	return gocassa.ConnectToKeySpace("messages", hosts, username, password)
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

func (db MessageDB) getMessageTable(k gocassa.KeySpace) gocassa.Table {
	mTable := k.Table("Messages", database.Message{}, gocassa.Keys{PartitionKeys: []string{"MessageID"}})

	mTable.CreateIfNotExist()

	return mTable
}

func (db *MessageDB) NewMessage(k *gocassa.KeySpace, cMsg *pb.CreateMessage, userID string) *pb.Message {
	id := MakeULID()

	mTable := db.getMessageTable(*k)

	msg := database.Message{
		MessageID: id.String(),
		UserID:    userID,
		Content:   cMsg.Content,
		CreatedAt: ulid.Time(id.Time()),
	}

	err := mTable.Set(msg).Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil
	}

	m := msg.MarshalToPB()

	db.MessageCacheExpiry[m.MessageID] = time.Now().Add(5 * time.Second)
	return m
}

func (db *MessageDB) GetMessage(id string) *pb.Message {
	return db.MessageCache[id]
}

func (db MessageDB) FetchMessage(k *gocassa.KeySpace, id string) *pb.Message {
	mTable := db.getMessageTable(*k)
	msg := database.Message{}
	err := mTable.Where(gocassa.Eq("id", id)).ReadOne(&msg).Run()

	if err != nil {
		panic(err)
	}

	return msg.MarshalToPB()
}

func (db MessageDB) MaybeGetMessage(k func() *gocassa.KeySpace, id string) *pb.Message {
	if msg := db.GetMessage(id); msg != nil {
		return msg
	}
	return db.FetchMessage(k(), id)
}

func (db MessageDB) FetchMessagesBefore(k *gocassa.KeySpace, t time.Time) *pb.MessageHistory {
	mTable := db.getMessageTable(*k)

	messages := make([]database.Message, 0)

	err := mTable.Where(gocassa.LTE("createdat", t)).Read(&messages).WithOptions(gocassa.Options{
		ClusteringOrder: []gocassa.ClusteringOrderColumn{
			{gocassa.DESC, "createdat"},
		},
		AllowFiltering: true,
	}).Run()

	if err != nil {
		panic(err)
	}

	msgHistory := make([]*pb.Message, 0)

	for _, msg := range messages {
		msgHistory = append(msgHistory, msg.MarshalToPB())
	}

	return &pb.MessageHistory{Messages: msgHistory}
}
