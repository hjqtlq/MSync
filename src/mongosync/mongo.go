package mongosync

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

//type mongoManager struct {
//	client *mongo.Client
//}

func NewClient(url string) *mongo.Client {
	client, err := mongo.Connect(SignalMonitor.Context, options.Client().ApplyURI(url))
	if err != nil {
		log.Println(err, url)
		return nil
	}
	defer SignalMonitor.BeforeClose(func() {
		log.Println("Close mongo connection.", url)
		if err := client.Disconnect(SignalMonitor.Context); err != nil {
			log.Println(err)
		}
	})
	return client
}
