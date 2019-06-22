package managers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type EsDocManager struct {
	host    string
	shardId string
}

func (dm *EsDocManager) BulkUpsert(docs []interface{}, ns string, ts *primitive.Timestamp) {
	log.Println("Es bulk upsert", ns, docs)
}
func (dm *EsDocManager) Update(docId primitive.ObjectID, updateSpec interface{}, ns string, ts *primitive.Timestamp) {
	log.Println("Es update")
}
func (dm *EsDocManager) Upsert(doc interface{}, ns string, ts *primitive.Timestamp) {
	log.Println("Es upsert")
}
func (dm *EsDocManager) Remove(docId primitive.ObjectID, ns string, ts *primitive.Timestamp) {
	log.Println("Es remove")
}
func (dm *EsDocManager) HandCommand(doc interface{}, ns string, ts *primitive.Timestamp) {
	log.Println("Es hand command")
}
func (dm *EsDocManager) GetLastDoc() {
}
func (dm *EsDocManager) GetTs() (ts *primitive.Timestamp) {
	return nil
}
func (dm *EsDocManager) SetTs(ts *primitive.Timestamp) {
}
func (dm *EsDocManager) Commit() {
}
func (dm *EsDocManager) Stop() {
}

func (dm *EsDocManager) Init() bool {
	return true
}
func (dm *EsDocManager) GetName() string {
	return "es"
}

func (dm *EsDocManager) SetShardId(shardId string) {

}
func (dm *EsDocManager) GetShardId() string {
	return dm.shardId
}

//func (dm *EsDocManager) SetName(name string) {
//}
