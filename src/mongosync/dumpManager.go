package mongosync

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type dumpManager struct {
	client *mongo.Client
}

func (dm *dumpManager) dump() {
	log.Println("Start dump")
	dbNames, _ := dm.client.ListDatabaseNames(SignalMonitor.Context, bson.M{})
	for _, dbName := range dbNames {
		if dbName == "config" || dbName == "local" {
			continue
		}
		db := dm.client.Database(dbName)
		collNames, err := db.ListCollectionNames(SignalMonitor.Context, bson.M{})
		if err != nil {
			log.Println(err)
		} else {
			for _, collName := range collNames {
				coll := db.Collection(collName)
				cur, err := coll.Find(SignalMonitor.Context, bson.M{}, options.Find().SetSort(bson.M{"_id": -1}))
				if err != nil {
					log.Println(err)
					continue
				}
				ns := dbName + "." + collName
				go func() {
					for cur.Next(SignalMonitor.Context) {
						var entry oplog
						err := cur.Decode(&entry)
						if err != nil {
							log.Println(err)
						}
						DocManager.Upsert(entry, ns, nil)
					}
				}()
			}
		}
	}
}

func newDumpManager(client *mongo.Client) *dumpManager {
	return &dumpManager{client: client}
}
