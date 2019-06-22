package mongosync

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.mongodb.org/mongo-driver/x/network/result"
	"gopkg.in/mgo.v2/bson"
	"log"
	"mongosync/managers"
	"strings"
	"sync"
	"time"
)

type MongoSync struct {
	client     *mongo.Client
	shardSet   map[string]*oplogManager
	runningMap map[string]bool
}

func (ms *MongoSync) run() {
	DocManager.Register(&managers.EsDocManager{})

	log.Println(Config.Mongo.Url)
	ms.client = NewClient(Config.Mongo.Url)
	if ms.client == nil {
		return
	}

	go Checkpoint.run()

	// if is mongos
	res := ms.client.Database("admin").RunCommand(SignalMonitor.Context, &bsonx.Doc{{"isDbGrid", bsonx.Boolean(true)}})
	log.Println(res.Err())
	if res.Err() != nil {
		log.Println("Start sync replica set")
		// Replica set
		isMasterRes := ms.client.Database("admin").RunCommand(SignalMonitor.Context, bsonx.Doc{{"isMaster", bsonx.Boolean(true)}})
		if isMasterRes.Err() != nil {
			log.Println(isMasterRes.Err())
			return
		}
		var isMaster result.IsMaster
		err := isMasterRes.Decode(&isMaster)
		if err != nil {
			log.Println(err)
			return
		}
		if isMaster.SetName == "" {
			log.Println("A replica set is required.")
			return
		}
		shardId := "replicaSet"
		ms.register(shardId, ms.client)
	} else {
		shards, err := ms.client.Database("config").Collection("shards").Find(SignalMonitor.Context, bson.D{})
		if err != nil {
			log.Println(err)
			return
		}
		for shards.Next(SignalMonitor.Context) {
			shard := bsonx.Doc{}
			res := shards.Decode(&shard)
			if res != nil {
				log.Println(res)
			} else {
				host := strings.SplitN(shard.Lookup("host").StringValue(), "/", 2)
				if len(host) == 2 {
					shardId := shard.Lookup("_id").ObjectID().Hex()
					if shardId != "" {
						ms.register(shardId, NewClient(host[1]))
					}
				}
			}
		}
	}
	go ms.monitor()
}

func (ms *MongoSync) monitor() {
	for {
		for id, oplogManager := range ms.shardSet {
			if !ms.runningMap[id] {
				ms.runningMap[id] = true

				ts := Checkpoint.get(id)

				go func() {
					if ts == nil {
						ts := newOplogManager(id, ms.client).getLastTs()
						Checkpoint.set(id, ts)
						Checkpoint.write()
						log.Println("Start dump", id)
						newDumpManager(oplogManager.client).dump()
						log.Println("Dump done", id)
					}
					log.Println("Start to run oplog manager", id, ts)
					oplogManager.run(ts)
				}()
			}
		}
		time.Sleep(time.Second * 3)
	}
}

func (ms *MongoSync) register(shardId string, client *mongo.Client) {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := ms.shardSet[shardId]; !ok {
		client := client
		ms.shardSet[shardId] = newOplogManager(shardId, client)
		ms.runningMap[shardId] = false
	}
}

func NewMongoSync() *MongoSync {
	return &MongoSync{
		shardSet:   map[string]*oplogManager{},
		runningMap: map[string]bool{},
	}
}

//func (ms *MongoSync) register() {
//
//}

func Run(config string) {
	if err := InitConfig(config); err == nil {
		var ms = NewMongoSync()
		wg := &sync.WaitGroup{}
		wg.Add(1)
		ms.run()
		wg.Wait()
		defer wg.Done()
	} else {
		log.Println(err)
	}
}
