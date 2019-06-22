package mongosync

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
)

const (
	OplogDatabaseName   = "local"
	OplogCollectionName = "oplog.rs"
	OplogOpUpdate       = "u"
	OplogOpDelete       = "d"
	OplogOpInsert       = "i"
	OplogOpCommand      = "c"
)

type oplog struct {
	H           int64                `bson:"h"`
	Ns          string               `bson:"ns"`
	O           bson.M               `bson:"o"`
	O2          bson.M               `bson:"o2"`
	Op          string               `bson:"op"`
	Ts          *primitive.Timestamp `bson:"ts"`
	T           int                  `bson:"t"`
	FromMigrate bool                 `bson:"fromMigrate"`
	Db          string
	Coll        string
}

func (o *oplog) Clear() {
	o.H = 0
	o.Ns = ""
	o.O = nil
	o.O2 = nil
	o.Op = ""
	o.Ts = nil
	o.T = 0
	o.FromMigrate = false
	o.Db = ""
	o.Coll = ""
}

type oplogManager struct {
	dm     *docManager
	ts     primitive.Timestamp
	client *mongo.Client
	id     string
}

func (om *oplogManager) run(ts *primitive.Timestamp) {

	oplogCur := om.getOplogCursor(ts)
	//SignalMonitor.BeforeClose(func() {
	//	err := oplogCur.Close(context.Background())
	//	log.Println("Close oplog cursor.", err)
	//})

	for oplogCur.Next(context.TODO()) {
		var entry oplog
		err := oplogCur.Decode(&entry)
		if err != nil {
			log.Println(err)
		}
		entry.Db, entry.Coll = ParseNamespace(entry.Ns)
		if om.shouldSkipEntry(&entry) {
			continue
		}
		om.dm.Process(entry)
	}
}

func (om *oplogManager) getLastTs() *primitive.Timestamp {
	var ts oplog
	err := om.getOplogColl().FindOne(
		SignalMonitor.Context,
		bson.M{"op": bson.M{"$ne": "n"}},
		options.FindOne().SetSort(bson.M{"$natural": -1}),
	).Decode(ts)
	if err != nil {
		log.Print(err)
	}
	return ts.Ts
}

func (om *oplogManager) shouldSkipEntry(entry *oplog) bool {
	// Ignore wrong entry
	if entry.O == nil {
		return true
	}
	// Don't replicate entries resulting from chunk moves
	if entry.FromMigrate {
		return true
	}
	// Ignore no-ops
	if entry.Op == "n" {
		return true
	}
	// Ignore none
	if entry.Db == "" || entry.Coll == "" {
		return true
	}
	//Ignore system collections
	if strings.HasPrefix(entry.Coll, "system.") {
		return true
	}
	//Ignore GridFS chunks
	if strings.HasSuffix(entry.Coll, ".chunks") {
		return true
	}
	//Ignore configured namespace
	if nsm.shouldSkipNamespace(entry.Db, entry.Coll) {
		return true
	}
	return false
}

func (om *oplogManager) getOplogColl() *mongo.Collection {
	return om.client.Database(OplogDatabaseName).Collection(OplogCollectionName)
}

func (om *oplogManager) getOplogCursor(ts *primitive.Timestamp) *mongo.Cursor {
	opt := options.Find().SetCursorType(options.TailableAwait)
	query := bson.M{"op": bson.M{"$ne": "n"}}
	if ts.T > 0 {
		query["ts"] = bson.M{"$gt": ts}
		opt.SetOplogReplay(true)
	}
	cur, err := om.getOplogColl().Find(SignalMonitor.Context, query, opt)
	if err != nil {
		fmt.Println(err)
	}
	return cur
}

func newOplogManager(id string, client *mongo.Client) *oplogManager {
	return &oplogManager{id: id, client: client, dm: NewDocManager(id)}
}
