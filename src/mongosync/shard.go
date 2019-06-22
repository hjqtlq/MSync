package mongosync

//
//import (
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/x/bsonx"
//	"strings"
//	"sync"
//)
//
//var ss = &shardSet{set: make(map[string]*shard)}
//
//type shardSet struct {
//	set map[string]*shard
//}
//
//func (ss *shardSet) register(shardId string, client *mongo.Client) {
//	mutex := &sync.Mutex{}
//	mutex.Lock()
//	defer mutex.Unlock()
//
//	if _, ok := ss.set[shardId]; !ok {
//		client := client
//		ss.set[shardId] = &shard{
//			oplogManager: newOplogManager(shardId, client),
//		}
//		//ss.set[shardId].run()
//	}
//}
//
//func (ss *shardSet) run() {
//	for id, shard := range ss.set {
//		ts := Checkpoint.get(id)
//
//		if ts == nil && Config.Dump {
//			ms.dump(id)
//			shard.oplogManager.run(ts)
//		} else {
//			shard.oplogManager.run(ts)
//		}
//	}
//}
//
//func (ss *shardSet) process(shardDoc *bsonx.Doc) {
//	shardId := shardDoc.Lookup("_id").ObjectID().Hex()
//	mutex := &sync.Mutex{}
//	mutex.Lock()
//
//	if _, ok := ss.set[shardId]; !ok {
//		hostInfo := strings.Split(shardDoc.Lookup("host").StringValue(), "/")
//		if len(hostInfo) == 2 {
//			client := NewClient(hostInfo[1])
//			ss.set[shardId] = &shard{
//				oplogManager: newOplogManager(shardId, client),
//			}
//			mutex.Unlock()
//			ss.set[shardId].run()
//		} else {
//			mutex.Unlock()
//		}
//	} else {
//		mutex.Unlock()
//	}
//}
//
//type shard struct {
//	oplogManager *oplogManager
//}
//
//func (s *shard) run() {
//	s.oplogManager.run()
//}
