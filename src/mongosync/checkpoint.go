package mongosync

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
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

func (cp *checkpoint) run() {
	SignalMonitor.BeforeClose(func() {
		cp.write()
		err := cp.getFile().Close()
		if err != nil {
			log.Println(err)
		}
	})
	for {
		cp.write()
		time.Sleep(time.Second)
	}
}

func (cp *checkpoint) set(shardId string, ts *primitive.Timestamp) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	cp.checkpointMap[shardId] = ts
}

func (cp *checkpoint) get(shardId string) *primitive.Timestamp {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	log.Println(cp.checkpointMap)

	if len(cp.checkpointMap) == 0 {
		if cm := cp.read(); len(cm) > 0 {
			cp.checkpointMap = cm
		} else {
			log.Println(cm)
		}
	}
	if ts, ok := cp.checkpointMap[shardId]; ok {
		return ts
	}
	return nil
}

func (cp *checkpoint) getFile() *os.File {
	//log.Println(Config.CheckpointPath + "checkpoint.json")
	if cp.file == nil {
		//cp.file, _ = os.Create(Config.CheckpointPath+"checkpoint.json")
		cp.file, _ = os.OpenFile(Config.CheckpointPath+"checkpoint.json", os.O_RDWR|os.O_CREATE, 0600)
	}
	return cp.file
}

func (cp *checkpoint) write() {
	if len(cp.checkpointMap) > 0 {
		b, _ := json.Marshal(cp.checkpointMap)
		log.Println("Write checkpoint", string(b))
		err := cp.getFile().Truncate(0)
		if err != nil {
			log.Println(err)
		}
		_, err = cp.getFile().WriteAt(b, 0)
		if err != nil {
			log.Println(err)
		}
	}
}

func (cp *checkpoint) read() map[string]*primitive.Timestamp {
	bty, err := ioutil.ReadAll(cp.getFile())
	//log.Println("Read checkpoint", string(bty))

	if err != nil {
		log.Println(err)
		return nil
	} else if len(bty) > 0 {
		var ts map[string]*primitive.Timestamp
		err := json.Unmarshal(bty, &ts)
		if err != nil {
			log.Println(err)
		} else {
			return ts
		}
	}
	return nil
}

//func (cp *checkpoint) Write(shardId string, ts ) {
//
//}
