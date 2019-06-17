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
	client, err := mongo.Connect(s.ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Println(err)
		return nil
	}
	defer s.beforeClose(func() {
		log.Println("Close mongo connection.", url)
		err := client.Disconnect(s.ctx)
		if err != nil {
			log.Println(err)
		}
	})
	return client
}
