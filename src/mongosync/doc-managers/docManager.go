package doc_managers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mongosync"
)

//type docManger struct {
//}
//
//func (dm *docManger) applyUpdate(doc interface{}, updateSpec interface{}) {
//
//}

type IDocManager interface {
	BulkUpsert(docs []interface{}, ns *mongosync.Namespace, ts primitive.Timestamp)
	Update(docId primitive.ObjectID, updateSpec interface{}, ns *mongosync.Namespace, ts primitive.Timestamp)
	Upsert(doc interface{}, ns *mongosync.Namespace, ts primitive.Timestamp)
	Remove(docId primitive.ObjectID, ns *mongosync.Namespace, ts primitive.Timestamp)
	HandCommand(doc interface{}, ns *mongosync.Namespace, ts primitive.Timestamp)
	GetLastDoc()
	Commit()
	Stop()
}
