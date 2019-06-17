package mongosync

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"sync"
	"time"
)

var Checkpoint = &checkpoint{checkpointMap: make(map[string]*primitive.Timestamp), mutex: &sync.Mutex{}}

type checkpoint struct {
	checkpointMap map[string]*primitive.Timestamp
	mutex         *sync.Mutex
	file          *os.File
}

func (cp *checkpoint) set(id string, ts *primitive.Timestamp) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	cp.checkpointMap[id] = ts
}

func (cp *checkpoint) get(shardId, docManagerName string) *primitive.Timestamp {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	if ts, ok := cp.checkpointMap[shardId]; ok {
		return ts
	}
	return nil
}

func (cp *checkpoint) monitor() {
	s.beforeClose(func() {
		cp.write()
		err := cp.file.Close()
		if err != nil {
			log.Println(err)
		}
	})
	for {
		cp.write()
		time.Sleep(time.Second)
	}

}

func (cp *checkpoint) write() {
	if cp.file == nil {
		cp.file, _ = os.OpenFile(Config.CheckpointPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0600)
	}
	b, err := json.Marshal(cp.checkpointMap)
	if err != nil {
		_, err = cp.file.Write(b)
		if err != nil {
			log.Println(err)
		}
	}
}
