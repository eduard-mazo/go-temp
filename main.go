package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func dbHandled(temp int, sensorID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://EASC_01:0KzQlxBvEk8yDmLV@thundersandbox.v0ydj.mongodb.net/ThunderSandBox?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	bmp280DB := client.Database("bmp280")
	collNames, err := bmp280DB.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(collNames)
	locations := bmp280DB.Collection("locations")
	/*locationAdd, err := locations.InsertOne(ctx, bson.D{
		{Key: "title", Value: "The Polyglot Developer Podcast"},
		{Key: "author", Value: "Nic Raboy"},
		{Key: "tags", Value: bson.A{"development", "programming", "coding"}},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(locationAdd)
	locationAdd, err = locations.InsertOne(ctx, bson.D{
		{Key: "sensorID", Value: 1234},
		{Key: "numSamples", Value: 1},
		{Key: "samples", Value: bson.A{
			bson.D{
				{Key: "val", Value: 70},
				{Key: "time", Value: 1535530440},
			},
			bson.D{
				{Key: "val", Value: 80},
				{Key: "time", Value: 1535530480},
			}},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(locationAdd)*/
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "sensorID", Value: sensorID}, {Key: "numSamples", Value: bson.D{{Key: "$lt", Value: 4}}}}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "samples", Value: bson.D{
				{Key: "val", Value: temp},
				{Key: "time", Value: time.Now().UnixNano()},
			}}},
		},
		{Key: "$inc", Value: bson.D{{Key: "numSamples", Value: 1}}},
	}
	locationUpdate, err := locations.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(locationUpdate)
}
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Sensor")
	}).Methods("GET")
	r.HandleFunc("/sensor/{sensorID}/temp/{temp}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sensorID, _ := strconv.Atoi(vars["sensorID"])
		temp, _ := strconv.Atoi(vars["temp"])

		fmt.Fprintf(w, "pushed sensor: %d with temp %d\n", sensorID, temp)
		dbHandled(temp, sensorID)
	})

	http.ListenAndServe(":"+port, r)
}
