package mongosync

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.mongodb.org/mongo-driver/x/network/result"
	"gopkg.in/mgo.v2/bson"
	"log"
	"mongosync/doc-managers"
	"sync"
)

var ms = &MongoSync{}

type MongoSync struct {
	config      *config
	docManagers []doc_managers.IDocManager
	client      *mongo.Client
}

func (mc *MongoSync) run() {
	mc.client = NewClient(Config.Mongo.Url)
	var res interface{}
	err := mc.client.Database("admin").RunCommand(s.ctx, bsonx.Doc{{"isdbgrid", bsonx.Boolean(true)}}).Decode(&res)
	if err != nil {
		log.Println(err)
		return
	}
	go Checkpoint.monitor()
	if res == nil {
		// Replica set
		isMasterRes := mc.client.Database("admin").RunCommand(s.ctx, bsonx.Doc{{"ismaster", bsonx.Boolean(true)}})
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
		go newOplogManager("replicaSet", mc.client).run()
	} else {
		shards, err := mc.client.Database("config").Collection("shards").Find(s.ctx, bson.D{})
		if err != nil {
			log.Println(err)
			return
		}
		for shards.Next(s.ctx) {
			shard := bsonx.Doc{}
			res := shards.Decode(&shard)
			if res != nil {
				log.Println(res)
			} else {
				go ss.process(&shard)
			}
		}
	}
}

func Run(config string) {
	if c, err := InitConfig(config); err == nil {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		ms.config = c
		ms.run()
		wg.Wait()
		defer wg.Done()
	} else {
		log.Println(err)
	}
}
