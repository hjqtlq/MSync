package managers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDocManager interface {
	BulkUpsert(docs []interface{}, ns string, ts *primitive.Timestamp)
	Update(docId primitive.ObjectID, updateSpec interface{}, ns string, ts *primitive.Timestamp)
	Upsert(doc interface{}, ns string, ts *primitive.Timestamp)
	Remove(docId primitive.ObjectID, ns string, ts *primitive.Timestamp)
	HandCommand(doc interface{}, ns string, ts *primitive.Timestamp)
	GetLastDoc()
	GetTs() (ts *primitive.Timestamp)
	SetTs(ts *primitive.Timestamp)
	Commit()
	Stop()

	Init() bool
	GetName() string
	//SetName(name string)
	GetShardId() string
	SetShardId(shardId string)
}
