package dbo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// CONNECTIONSTRING DB connection string
const CONNECTIONSTRING = "mongodb+srv://EASC_01:0KzQlxBvEk8yDmLV@thundersandbox.v0ydj.mongodb.net/ThunderSandBox?retryWrites=true&w=majority"

// DBNAME Database name
const DBNAME = "bmp280"

// COLLNAME Collection name
const COLLNAME = "locations"

const TREEHOURS = 10800000000000

var db *mongo.Database

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(CONNECTIONSTRING))
	if err != nil {
		log.Fatal(err)
	}
	/*defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()*/
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
	db = client.Database(DBNAME)
}

func Update(sensorID int, temp float64) {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "sensorID", Value: sensorID}, {Key: "numSamples", Value: bson.D{{Key: "$lt", Value: 250}}}}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "samples", Value: bson.D{
				{Key: "val", Value: temp},
				{Key: "time", Value: time.Now().UnixNano()},
			}}},
		},
		{Key: "$inc", Value: bson.D{{Key: "numSamples", Value: 1}}},
	}
	locationUpdate, err := db.Collection(COLLNAME).UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(locationUpdate)
}

func Get(sensorID int) []bson.M {
	timeNow := time.Now().UnixNano() - TREEHOURS
	filter := bson.D{{Key: "sensorID", Value: sensorID}, {Key: "samples.time", Value: bson.D{{Key: "$gt", Value: timeNow}}}}
	cursor, err := db.Collection(COLLNAME).Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	var measures []bson.M
	var i int
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var measure bson.M
		if err = cursor.Decode(&measure); err != nil {
			log.Fatal(err)
		}
		measures = append(measures, measure)
		i++
	}
	return measures
}
