package mongosync

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"reflect"
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
	H           int64               `bson:"h"`
	Ns          string              `bson:"ns"`
	O           bson.M              `bson:"o"`
	O2          bson.M              `bson:"o2"`
	Op          string              `bson:"op"`
	Ts          primitive.Timestamp `bson:"ts"`
	T           int                 `bson:"t"`
	FromMigrate bool                `bson:"fromMigrate"`
	namespace   *Namespace
}

type oplogManager struct {
	ts     primitive.Timestamp
	client *mongo.Client
	id     string
}

func (om *oplogManager) run() {
	oplogCur := om.getOplogCursor()
	s.beforeClose(func() {
		err := oplogCur.Close(context.Background())
		fmt.Println("Close oplog cursor", err)
	})

	for _, dm := range ms.docManagers {
		s.beforeClose(func() {
			dm.Stop()
			fmt.Println("Stop doc manager.", reflect.ValueOf(dm))
		})
	}

	for oplogCur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var entry oplog
		err := oplogCur.Decode(&entry)
		if err != nil {
			log.Fatal(err)
		}
		entry.namespace = NewNamespace(entry.Ns)

		if om.shouldSkipEntry(&entry) {
			//log.Println("Skip entry:", entry)
			continue
		}

		for _, dm := range ms.docManagers {
			go func() {
				switch entry.Op {
				case OplogOpDelete:
					if id, ok := entry.O["_id"]; ok {
						dm.Remove(id.(primitive.ObjectID), entry.namespace, entry.Ts)
					}
				case OplogOpInsert:
					dm.Upsert(entry.O, entry.namespace, entry.Ts)
				case OplogOpUpdate:
					if id, ok := entry.O2["_id"].(primitive.ObjectID); ok {
						dm.Update(id, entry.O, entry.namespace, entry.Ts)
					}
				case OplogOpCommand:
					dm.HandCommand(entry.O, entry.namespace, entry.Ts)
				}
			}()
		}
		Checkpoint.set(om.id, &entry.Ts)
	}
}

func (om *oplogManager) shouldSkipEntry(entry *oplog) bool {
	// Don't replicate entries resulting from chunk moves
	if entry.FromMigrate {
		return true
	}
	// Ignore no-ops
	if entry.Op == "n" {
		return true
	}
	// Ignore none
	if entry.namespace.Db == "" || entry.namespace.Coll == "" {
		return true
	}
	//Ignore system collections
	if strings.HasPrefix(entry.namespace.Coll, "system.") {
		return true
	}
	//Ignore GridFS chunks
	if strings.HasSuffix(entry.namespace.Coll, ".chunks") {
		return true
	}
	//Ignore configured namespace
	if nsm.shouldSkipNamespace(entry.namespace.Db, entry.namespace.Coll) {
		return true
	}
	return false
}

func (om *oplogManager) getOplogCursor() *mongo.Cursor {
	coll := om.client.Database(OplogDatabaseName).Collection(OplogCollectionName)
	opt := options.Find().SetCursorType(options.TailableAwait)
	query := bson.M{"op": bson.M{"$ne": "n"}}
	//if o.timestamp > 0 {
	//	query["ts"] = bson.M{"gte": o.timestamp}
	//	opt.SetOplogReplay(true)
	//}
	cur, err := coll.Find(s.ctx, query, opt)
	if err != nil {
		fmt.Println(err)
	}
	return cur
}

func newOplogManager(id string, client *mongo.Client) *oplogManager {
	return &oplogManager{id: id, client: client}
}
