package mongosync

import (
	"go.mongodb.org/mongo-driver/x/bsonx"
	"strings"
	"sync"
)

var ss = &shardSet{set: make(map[string]*shard)}

type shardSet struct {
	set map[string]*shard
}

func (ss *shardSet) process(shardDoc *bsonx.Doc) {
	shardId := shardDoc.Lookup("_id").ObjectID().Hex()
	mutex := &sync.Mutex{}
	mutex.Lock()

	if s, ok := ss.set[shardId]; !ok {
		s.info = shardDoc
		hostInfo := strings.Split(shardDoc.Lookup("host").StringValue(), "/")
		if len(hostInfo) == 2 {
			client := NewClient(hostInfo[1])
			ss.set[shardId] = &shard{
				info:         shardDoc,
				oplogManager: newOplogManager(shardId, client),
			}
			mutex.Unlock()
			ss.set[shardId].run()
		} else {
			mutex.Unlock()
		}
	} else {
		mutex.Unlock()
	}
}

type shard struct {
	info         *bsonx.Doc
	oplogManager *oplogManager
}

func (s *shard) run() {
	s.oplogManager.run()
}
