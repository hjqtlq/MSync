package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://94.191.9.80:27017"))

	coll := client.Database("test").Collection("a")

	fmt.Println("1")

	lcur, _ := coll.Database().ListCollections(context.Background(), bson.D{})
	defer func() {
		err := lcur.Close(ctx)
		if err != nil {
			log.Println(err)
		}
	}()
	for lcur.Next(ctx) {
		var result bson.M
		err := lcur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	//mongosync.Run("D:\\mongosync\\src\\config.yaml")
	//
	//fmt.Println("进程启动...")
	//sum := 0
	//for {
	//	sum++
	//	fmt.Println("sum:", sum)
	//	time.Sleep(time.Second)
	//}

}
